package graph

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/config"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.uber.org/zap"
	"strings"
)

func InitGraphDatabase(ctx context.Context, appConfig *config.AppConfig) error {

	l := logging.WithContext(ctx)

	systemLeaderSession, err := getDatabaseLeaderSession(appConfig, "system")
	if err != nil {
		l.Error("could not get database leader session", zap.Error(err))
		return err
	}
	defer systemLeaderSession.Close()

	result, err := systemLeaderSession.Run("CREATE DATABASE "+appConfig.Neo4jDatabase+" IF NOT EXISTS", map[string]interface{}{})
	if err != nil {
		l.Error("could not create database", zap.Error(err))
		return err
	}
	if result.Err() != nil {
		l.Error("could not create database", zap.Error(result.Err()))
		return result.Err()
	}

	dbSession, err := getDatabaseLeaderSession(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		l.Error("could not get database leader session", zap.Error(err))
		return err
	}

	return initGraphConstraints(ctx, dbSession)

}

func initGraphConstraints(ctx context.Context, session neo4j.Session) error {

	l := logging.WithContext(ctx)

	nodeNames := []string{
		"User",
		"Resource",
		"Offer",
		"OfferItem",
		"Group",
	}
	for _, nodeName := range nodeNames {
		if err := createIdConstraint(ctx, session, nodeName); err != nil {
			l.Error("could not create constraint", zap.Error(err))
			return err
		}
	}
	return nil
}

func createIdConstraint(ctx context.Context, session neo4j.Session, nodeName string) error {

	l := logging.WithContext(ctx)

	result, err := session.Run(`CREATE CONSTRAINT IF NOT EXISTS idx`+nodeName+` ON (n:`+nodeName+`) ASSERT n.id IS UNIQUE`, map[string]interface{}{})
	if err != nil {
		l.Error("could not create constraint", zap.Error(err))
		return err
	}

	if result.Err() != nil {

	}

	return nil

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

func getDatabaseLeaderSession(appConfig *config.AppConfig, databaseName string) (neo4j.Session, error) {

	tempDriver, err := neo4j.NewDriver(
		appConfig.BoltUrl,
		neo4j.BasicAuth(appConfig.BoltUsername, appConfig.BoltPassword, ""),
		func(c *neo4j.Config) {

		})

	if err != nil {
		return nil, fmt.Errorf("could not create neo4j driver: %v", err)
	}

	session := tempDriver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		Bookmarks:    nil,
		DatabaseName: databaseName,
	})

	return session, nil

	defer session.Close()

	leaderBoltUrl, err := findDatabaseLeaderBoltUrl(session, databaseName)
	if err != nil {
		return nil, err
	}

	leaderDriver, err := neo4j.NewDriver(
		leaderBoltUrl,
		neo4j.BasicAuth(appConfig.BoltUsername, appConfig.BoltPassword, ""),
		func(c *neo4j.Config) {

		})
	if err != nil {
		return nil, fmt.Errorf("could not create system leader neo4j driver: %v", err)
	}

	leaderSession := leaderDriver.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		Bookmarks:    nil,
		DatabaseName: databaseName,
	})

	return leaderSession, nil

}
