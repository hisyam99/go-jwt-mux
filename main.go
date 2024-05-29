package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hisyam99/go-jwt-mux/controllers/authcontroller"
	"github.com/hisyam99/go-jwt-mux/controllers/productcontroller"
	"github.com/hisyam99/go-jwt-mux/middlewares"
	"github.com/hisyam99/go-jwt-mux/models"
)

func main() {
	models.ConnectDatabase()
	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productcontroller.Index).Methods("GET")
	api.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
