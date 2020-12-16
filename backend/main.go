package main

import (
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/pkg/server"
	"log"
	"os"
	"os/signal"
)

// @title commonpool api
// @version 1.0
// @description resources api
// @termsOfService http://swagger.io/terms
// @host 127.0.0.1:8585
// @basePath /api/v1
func main() {

	var (
		srv *server.Server
		err error
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if srv, err = server.NewServer(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit

	srv.Shutdown()

}
