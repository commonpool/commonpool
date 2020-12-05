package graph

import (
	"fmt"
	"github.com/commonpool/backend/config"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"strings"
)

func InitGraphDatabase(appConfig *config.AppConfig) error {

	systemLeaderSessionPtr, err := GetDatabaseLeaderSession(appConfig, "system")
	if err != nil {
		return err
	}
	defer (*systemLeaderSessionPtr).Close()
	systemLeaderSession := *systemLeaderSessionPtr

	result, err := systemLeaderSession.Run("CREATE DATABASE "+appConfig.Neo4jDatabase+" IF NOT EXISTS", map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("could not create database: %v", err)
	}
	if result.Err() != nil {
		return fmt.Errorf("could not create database: %v", err)
	}

	return nil

}

type GraphDriver interface {
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

var _ GraphDriver = &Neo4jGraphDriver{}

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
	defer session.Close()

	if err != nil {
		return nil, fmt.Errorf("could not open connection: %v", err)
	}

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

func findDatabaseLeaderBoltUrl(session neo4j.Session, databaseName string) (string, error) {

	result, err := session.Run("CALL dbms.cluster.overview()", map[string]interface{}{})

	if err != nil {
		return "", fmt.Errorf("could not create database: %v", err)
	}

	if result.Err() != nil {
		return "", fmt.Errorf("could not create database: %v", result.Err())
	}

	for result.Next() {

		addressesIntf, ok := result.Record().Get("addresses")
		if !ok {
			return "", fmt.Errorf("could not get addresses field: %v", err)
		}
		databasesIntf, ok := result.Record().Get("databases")
		if !ok {
			return "", fmt.Errorf("could not get databases field: %v", err)
		}

		var addresses []string
		err = mapstructure.Decode(addressesIntf, &addresses)
		if err != nil {
			return "", fmt.Errorf("could not decode addresses: %v", err)
		}

		var boltAddress string
		for _, address := range addresses {
			if strings.Index(address, "bolt://") == 0 {
				boltAddress = address
			}
		}

		var databases map[string]string
		err = mapstructure.Decode(databasesIntf, &databases)
		if err != nil {
			return "", fmt.Errorf("could not decode detabases: %v", err)
		}

		if databases[databaseName] == "LEADER" {
			return boltAddress, nil
		}
	}

	return "", fmt.Errorf("could not find leader bolt address")

}

func GetDatabaseLeaderSession(appConfig *config.AppConfig, databaseName string) (*neo4j.Session, error) {

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
	defer session.Close()

	if err != nil {
		return nil, fmt.Errorf("could not open connection: %v", err)
	}

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

	leaderSession, err := leaderDriver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		Bookmarks:    nil,
		DatabaseName: databaseName,
	})

	return &leaderSession, nil

}
