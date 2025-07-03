package app

import (
	"auth-service/internal/controller"

	"github.com/gorilla/mux"
)

func initRoutes(controller controller.Controller) *mux.Router {
	router := mux.NewRouter()

	return router
}
