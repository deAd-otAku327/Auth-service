package app

import (
	"auth-service/internal/controller"
	"auth-service/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func initRoutes(controller controller.Controller, mw middleware.Middleware) *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/auth").Subrouter()

	router.HandleFunc("/login", controller.HandleLogin()).Methods(http.MethodPost)
	router.HandleFunc("/refresh", controller.HandleRefresh()).Methods(http.MethodPost)

	protected := router.NewRoute().Subrouter()
	protected.Use(mw.Auth())

	protected.HandleFunc("/current", controller.HandleGetCurrentUser()).Methods(http.MethodGet)
	protected.HandleFunc("/logout", controller.HandleLogout()).Methods(http.MethodPost)

	return router
}
