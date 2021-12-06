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
	con, err := db.GetDBconnection()

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
	routes.SetupRoutes(router)

	port := ":5500"

	server := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server started at %s", port)
	log.Fatal(server.ListenAndServe())
}
