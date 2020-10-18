package main

import (
	"fmt"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/store"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	d  *gorm.DB
	rs resource.Store
	e  *echo.Echo
)

type DbConfig struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

// @title resources api
// @version 1.0
// @description resources api
// @termsOfService http://swagger.io/terms
// @host 127.0.0.1:8585
// @basePath /api/v1
func main() {
	r := router.NewRouter()

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api/v1")

	db := getDb()

	store.AutoMigrate(db)
	rs = store.NewResourceStore(db)
	h := handler.NewHandler(rs)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}

func getDb() *gorm.DB {
	dbConfig, err := getDbConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}

	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port)
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

type EnvReader func(string) (string, bool)
type FileReader func(string) ([]byte, error)

func getDbConfig(readEnv EnvReader, readFile FileReader) (*DbConfig, error) {
	dbUser, hasDbUser := readEnv("DB_USER")
	dbUserFile, hasDbUserFile := readEnv("DB_USER_FILE")

	if !hasDbUser && !hasDbUserFile {
		panic(fmt.Errorf("DB_USER or DB_USER_FILE environment variable is required"))
	}

	if hasDbUserFile {
		dbUserIo, err := readFile(dbUserFile)
		if err != nil {
			return nil, err
		}
		dbUser = string(dbUserIo)
	}

	dbPassword, hasDbPassword := readEnv("DB_PASSWORD")
	dbPasswordFile, hasDbPasswordFile := readEnv("DB_PASSWORD_FILE")

	if !hasDbPassword && !hasDbPasswordFile {
		return nil, fmt.Errorf("DB_PASSWORD or DB_PASSWORD_FILE environment variable is required")
	}

	if hasDbPasswordFile {
		dbUserIo, err := readFile(dbPasswordFile)
		if err != nil {
			return nil, err
		}
		dbPassword = string(dbUserIo)
	}

	dbName, ok := readEnv("DB_NAME")
	if !ok {
		return nil, fmt.Errorf("DB_NAME environment variable is required")
	}
	dbPortStr, ok := readEnv("DB_PORT")
	if !ok {
		return nil, fmt.Errorf("DB_PORT environment variable is required")
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, err
	}

	dbHost, ok := readEnv("DB_HOST")
	if !ok {
		return nil, fmt.Errorf("DB_HOST environment variable is required")
	}

	dbConfig := &DbConfig{
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
		Username: dbUser,
		Password: dbPassword,
	}
	return dbConfig, nil
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
