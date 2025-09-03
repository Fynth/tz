package main

import (
	"net/http"
	handler "tz/internal/handlers"
	"tz/internal/repository"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscription API
// @version 1.0
// @description API для управления подписками
// @host localhost:8080
// @BasePath /
func main() {
	repo := repository.NewSubscriptionRepository()

	handler := handler.NewSubscriptionHandler(repo)

	r := mux.NewRouter()
	r.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/subs/", handler.Create).Methods("POST")
	r.HandleFunc("/subs/", handler.List).Methods("GET")
	r.HandleFunc("/subs/total/", handler.Total).Methods("GET")
	r.HandleFunc("/subs/{id}/", handler.Retrieve).Methods("GET")
	r.HandleFunc("/subs/{id}/", handler.Update).Methods("PUT")
	r.HandleFunc("/subs/{id}/", handler.Delete).Methods("DELETE")

	log.Info()
	http.ListenAndServe(":8080", r)
}
