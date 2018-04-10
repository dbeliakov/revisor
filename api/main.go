package main

import (
	"net/http"
	"os"
	"reviewer/api/handlers"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")

	corsRouter := gh.CORS(
		gh.AllowedOrigins([]string{"*"}),
		gh.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		gh.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gh.ExposedHeaders([]string{"Authorization"}),
	)(r)

	http.ListenAndServe(":8090", gh.LoggingHandler(os.Stdout, corsRouter))
}
