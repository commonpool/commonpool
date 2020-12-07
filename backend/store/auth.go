package store

import (
	"context"
	"errors"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/model"
	"github.com/labstack/gommon/log"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"gorm.io/gorm"
	"strings"
)

type AuthStore struct {
	db          *gorm.DB
	graphDriver graph.GraphDriver
}

var _ auth.Store = &AuthStore{}

func NewAuthStore(db *gorm.DB, graphDriver graph.GraphDriver) *AuthStore {
	return &AuthStore{
		db:          db,
		graphDriver: graphDriver,
	}
}

type UserStore struct {
	db *gorm.DB
}

func (as *AuthStore) GetByKeys(ctx context.Context, keys []model.UserKey) (*auth.Users, error) {

	session, err := as.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return as.getByKeys(session, model.NewUserKeys(keys))
}

func (as *AuthStore) Upsert(key model.UserKey, email string, username string) error {

	session, err := as.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	usr, err := as.getByKey(session, key)

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

	} else if errors.Is(err, errs.ErrUserNotFound) {

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

func (as *AuthStore) getByKey(session neo4j.Session, key model.UserKey) (*auth.User, error) {

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
		return nil, errs.ErrUserNotFound
	}

	record := getResult.Record()
	userRecord, _ := record.Get("n")
	return MapUserNode(userRecord.(neo4j.Node)), nil

}

func (as *AuthStore) getByKeys(session neo4j.Session, key *model.UserKeys) (*auth.Users, error) {

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

	var users []*auth.User

	for getResult.Next() {
		record := getResult.Record()
		userRecord, _ := record.Get("n")
		node := userRecord.(neo4j.Node)
		users = append(users, MapUserNode(node))
	}

	_, err = getResult.Consume()

	return auth.NewUsers(users), err

}

func (as *AuthStore) GetByKey(key model.UserKey) (*auth.User, error) {
	session, err := as.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	user, err := as.getByKey(session, key)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (as *AuthStore) GetUsername(key model.UserKey) (string, error) {
	user, err := as.GetByKey(key)
	if err != nil {
		return "", err
	}
	return user.Username, err
}

func (as *AuthStore) Find(query auth.UserQuery) ([]*auth.User, error) {

	session, err := as.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var whereClauses = []string{}
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

	var users []*auth.User

	for result.Next() {
		record := result.Record()
		field, _ := record.Get("u")
		users = append(users, MapUserNode(field.(neo4j.Node)))
	}

	return users, err
}
func IsUserNode(node neo4j.Node) bool {
	return NodeHasLabel(node, "User")
}

func MapUserNode(node neo4j.Node) *auth.User {
	return &auth.User{
		ID:       node.Props()["id"].(string),
		Username: node.Props()["username"].(string),
		Email:    node.Props()["email"].(string),
	}
}
