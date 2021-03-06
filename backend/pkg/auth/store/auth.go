package store

import (
	"context"
	"errors"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/exceptions"
	graph2 "github.com/commonpool/backend/pkg/graph"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/gommon/log"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"strings"
)

type UserStore struct {
	graphDriver graph2.Driver
}

var _ Store = &UserStore{}

func NewUserStore(graphDriver graph2.Driver) *UserStore {
	return &UserStore{
		graphDriver: graphDriver,
	}
}

func (us *UserStore) GetByKeys(ctx context.Context, keys *keys.UserKeys) (*models.Users, error) {

	session := us.graphDriver.GetSession()
	defer session.Close()

	return us.getByKeys(session, keys)
}

func (us *UserStore) Upsert(key keys.UserKey, email string, username string) error {

	session := us.graphDriver.GetSession()
	defer session.Close()

	usr, err := us.getByKey(session, key)

	if err == nil {

		if usr.Username == username && usr.Email == email {
			return nil
		}

		updateResult, err := session.Run(`
			MATCH (n:User {id:$id})
			SET n += {username: $username, email: $email}
			RETURN n`,
			map[string]interface{}{
				"id":       key.String(),
				"username": username,
				"email":    email,
			})

		if err != nil {
			log.Errorf("could not update user: %v", err)
			return err
		}

		if updateResult.Err() != nil {
			log.Errorf("could not update user: %v", updateResult.Err())
			return updateResult.Err()
		}

	} else if errors.Is(err, exceptions.ErrUserNotFound) {

		createResult, err := session.Run(`
			CREATE (u:User {id:$id, username: $username, email:$email}) 
			RETURN u`,
			map[string]interface{}{
				"id":       key.String(),
				"username": username,
				"email":    email,
			})

		if err != nil {
			log.Errorf("could not create user: %v", err)
			return err
		}

		if createResult.Err() != nil {
			log.Errorf("could not create user: %v", createResult.Err())
			return createResult.Err()
		}

	} else {
		return err
	}

	return nil

}

func (us *UserStore) getByKey(session neo4j.Session, key keys.UserKey) (*models.User, error) {

	getResult, err := session.Run(`
		MATCH (n:User {id:$id}) 
		RETURN n`,
		map[string]interface{}{
			"id": key.String(),
		})

	if err != nil {
		log.Errorf("could not get user: %v", err)
		return nil, err
	}

	if getResult.Err() != nil {
		log.Errorf("could not get user: %v", getResult.Err())
		return nil, getResult.Err()
	}

	if getResult.Next() {

	} else {
		return nil, exceptions.ErrUserNotFound
	}

	record := getResult.Record()
	userRecord, _ := record.Get("n")
	return MapUserNode(userRecord.(neo4j.Node)), nil

}

func (us *UserStore) getByKeys(session neo4j.Session, key *keys.UserKeys) (*models.Users, error) {

	getResult, err := session.Run(`
		MATCH (n:User) 
		WHERE n.id IN $ids
		RETURN n`,
		map[string]interface{}{
			"ids": key.Strings(),
		})

	if err != nil {
		log.Errorf("could not get users: %v", err)
		return nil, err
	}

	if getResult.Err() != nil {
		log.Errorf("could not get users: %v", getResult.Err())
		return nil, getResult.Err()
	}

	var users []*models.User

	for getResult.Next() {
		record := getResult.Record()
		userRecord, _ := record.Get("n")
		node := userRecord.(neo4j.Node)
		users = append(users, MapUserNode(node))
	}

	_, err = getResult.Consume()

	return models.NewUsers(users), err

}

func (us *UserStore) GetByKey(key keys.UserKey) (*models.User, error) {
	session := us.graphDriver.GetSession()
	defer session.Close()

	u, err := us.getByKey(session, key)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *UserStore) GetUsername(key keys.UserKey) (string, error) {
	u, err := us.GetByKey(key)
	if err != nil {
		return "", err
	}
	return u.Username, err
}

func (us *UserStore) Find(query Query) (*models.Users, error) {

	session := us.graphDriver.GetSession()
	defer session.Close()

	var whereClauses []string
	var params = map[string]interface{}{}

	cypher := []string{
		"MATCH (u:User)",
	}

	if strings.TrimSpace(query.Query) != "" {
		whereClauses = append(whereClauses, "u.username ~= '.*$query.*'")
		params["query"] = strings.TrimSpace(query.Query)
	}

	if query.NotInGroup != nil {
		whereClauses = append(whereClauses, "NOT (u)-[:IsMemberOf]->(:Group {id:$groupId})")
		params["groupId"] = query.NotInGroup.String()
	}

	if len(whereClauses) > 0 {
		cypher = append(cypher, "WHERE "+strings.Join(whereClauses, " AND "))
	}

	cypher = append(cypher, "RETURN u")
	cypher = append(cypher, "ORDER BY u.username")

	result, err := session.Run(strings.Join(cypher, "\n"), params)

	if err != nil {
		return nil, err
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	var users []*models.User

	for result.Next() {
		record := result.Record()
		field, _ := record.Get("u")
		users = append(users, MapUserNode(field.(neo4j.Node)))
	}

	return models.NewUsers(users), err
}

func IsUserNode(node neo4j.Node) bool {
	return graph2.NodeHasLabel(node, "User")
}

func MapUserNode(node neo4j.Node) *models.User {
	return &models.User{
		ID:       node.Props["id"].(string),
		Username: node.Props["username"].(string),
		Email:    node.Props["email"].(string),
	}
}
