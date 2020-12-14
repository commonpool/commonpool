package graph

import (
	"fmt"
	"github.com/commonpool/backend/pkg/config"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Driver interface {
	GetSession() (neo4j.Session, error)
}

type Neo4jGraphDriver struct {
	driver       neo4j.Driver
	databaseName string
}

func (n Neo4jGraphDriver) GetSession() (neo4j.Session, error) {
	sess, err := n.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		Bookmarks:    nil,
		DatabaseName: n.databaseName,
	})
	return sess, err
}

var _ Driver = &Neo4jGraphDriver{}

func NewNeo4jDriver(appConfig *config.AppConfig, databaseName string) (*Neo4jGraphDriver, error) {

	tempDriver, err := neo4j.NewDriver(
		appConfig.BoltUrl,
		neo4j.BasicAuth(appConfig.BoltUsername, appConfig.BoltPassword, ""),
		func(c *neo4j.Config) {
			c.Encrypted = false
		})

	if err != nil {
		return nil, fmt.Errorf("could not create neo4j driver: %v", err)
	}

	session, err := tempDriver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		Bookmarks:    nil,
		DatabaseName: "system",
	})
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %v", err)
	}

	defer session.Close()

	leaderBoltUrl, err := findDatabaseLeaderBoltUrl(session, databaseName)
	if err != nil {
		return nil, err
	}

	leaderDriver, err := neo4j.NewDriver(
		leaderBoltUrl,
		neo4j.BasicAuth(appConfig.BoltUsername, appConfig.BoltPassword, ""),
		func(c *neo4j.Config) {
			c.Encrypted = false
		})
	if err != nil {
		return nil, fmt.Errorf("could not create system leader neo4j driver: %v", err)
	}

	return &Neo4jGraphDriver{
		driver:       leaderDriver,
		databaseName: databaseName,
	}, nil

}
