package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anuragch/go_auth/configs"
	userHandler "github.com/Anuragch/go_auth/handlers/user"
	"github.com/Anuragch/go_auth/repository"
	"github.com/Anuragch/go_auth/utils/dbutils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	dbInstance, err := dbutils.InitDb(&configs.MongoConfig)
	repo := repository.NewUserRepo(dbInstance, "Users")
	if err != nil {
		fmt.Println("Could not connect to db instance")
		return
	}
	sm := mux.NewRouter()
	uh := userHandler.NewUserHandler(log, repo)
	sr := sm.Methods(http.MethodPost).Subrouter()
	sr.HandleFunc("/signup", uh.CreateUser)
	sr.HandleFunc("/login", uh.AuthenticateUser)
	sr.HandleFunc("/refresh", uh.RefreshToken)

	sr = sm.Methods(http.MethodGet).Subrouter()
	sr.HandleFunc("/users", uh.GetUsers)
	sr.HandleFunc("/user/{id}", uh.GetUser)
	s := http.Server{
		Addr:         "localhost:8080",
		Handler:      sm,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(&s)
}

func defaultHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Calling deafult handler")
}
func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
