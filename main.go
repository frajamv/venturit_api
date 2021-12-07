package main

import (
	"log"
	"moviesapi/db"
	"moviesapi/routes"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	con, err := db.GetDBconnection() // Connect to DB to test it's connectivity.

	if err != nil {
		log.Println("error with database " + err.Error())
	} else {
		err = con.Ping()
		if err != nil {
			log.Println("error making conection to DB, error: " + err.Error())
			return
		}
	}

	router := mux.NewRouter()
	routes.SetupRoutes(router) // Setup all routes to be fetched by the REST client.

	port := ":5500"

	server := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
