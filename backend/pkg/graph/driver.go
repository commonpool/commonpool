package graph

import (
	"fmt"
	"github.com/commonpool/backend/pkg/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Driver interface {
	GetSession() neo4j.Session
}

type Neo4jGraphDriver struct {
	driver       neo4j.Driver
	databaseName string
}

func (n Neo4jGraphDriver) GetSession() neo4j.Session {
	return n.driver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		Bookmarks:    nil,
		DatabaseName: n.databaseName,
	})
}

var _ Driver = &Neo4jGraphDriver{}

func NewNeo4jDriver(appConfig *config.AppConfig, databaseName string) (*Neo4jGraphDriver, error) {

	driver, err := neo4j.NewDriver(
		appConfig.BoltUrl,
		neo4j.BasicAuth(appConfig.BoltUsername, appConfig.BoltPassword, ""),
		func(c *neo4j.Config) {
		})

	if err != nil {
		return nil, fmt.Errorf("could not create neo4j driver: %v", err)
	}

	return &Neo4jGraphDriver{
		driver:       driver,
		databaseName: databaseName,
	}, nil

}
