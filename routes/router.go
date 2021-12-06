package routes

import (
	"moviesapi/controllers"
	"moviesapi/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		utils.SendSuccessResponse("REST API connection success", res)
	}).Methods(http.MethodGet)

	subRouter := router.PathPrefix("/movies").Subrouter()
	subRouter.HandleFunc("/all", controllers.GetAllMovies).Methods(http.MethodGet)
	subRouter.HandleFunc("/update", controllers.UpdateMovie).Methods(http.MethodPatch)
}
