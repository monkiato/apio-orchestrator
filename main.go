package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"monkiato/apio-orchestrator/pkg/server"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)

	if debugMode, found := os.LookupEnv("DEBUG_MODE"); found {
		if val, _ := strconv.Atoi(debugMode); val == 1 {
			log.SetLevel(log.DebugLevel)
		}
	}

	port := "80"
	if customPort, found := os.LookupEnv("SERVER_PORT"); found {
		port = customPort
	}


	mainRoute := mux.NewRouter().PathPrefix("/api/").Subrouter()
	addAPIRoutes(mainRoute)

	srv := &http.Server{
		Handler: mainRoute,
		Addr:    fmt.Sprintf(":%s", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Debugf("server ready. Running at %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func addAPIRoutes(router *mux.Router) {
	apiRoute := router.PathPrefix("node").Subrouter()
	apiRoute.HandleFunc("/", server.ParseBody(server.CreateNodeHandler)).Methods(http.MethodPut)
	apiRoute.HandleFunc("/{id}", server.ParseBody(server.EditNodeHandler)).Methods(http.MethodPost)
	apiRoute.HandleFunc("/{id}", server.RemoveNodeHandler).Methods(http.MethodDelete)
}