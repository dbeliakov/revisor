package main

import (
	"net/http"
	"os"
	"reviewer/api/auth"
	"reviewer/api/comments"
	"reviewer/api/config"
	"reviewer/api/revisions"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func addAPIHandlers(base string, r *mux.Router) {
	// Auth handlers
	r.HandleFunc(base+"/auth/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc(base+"/auth/signup", auth.SignUpHandler).Methods("POST")
	r.HandleFunc(base+"/auth/user", auth.UserInfoHandler).Methods("GET")
	r.HandleFunc(base+"/auth/change/password", auth.ChangePasswordHandler).Methods("POST")

	// Review handlers
	r.HandleFunc(base+"/reviews/outgoing", revisions.OutgoingReviews).Methods("GET")
	r.HandleFunc(base+"/reviews/incoming", revisions.IncomingReviews).Methods("GET")
	r.HandleFunc(base+"/reviews/new", revisions.NewReview).Methods("POST")
	r.HandleFunc(base+"/reviews/{id}", revisions.Review).Methods("GET")
	r.HandleFunc(base+"/reviews/{id}/update", revisions.UpdateReview).Methods("POST")
	r.HandleFunc(base+"/reviews/{id}/accept", revisions.Accept).Methods("GET")
	r.HandleFunc(base+"/reviews/{id}/decline", revisions.Decline).Methods("GET")

	// Comments handlers
	r.HandleFunc(base+"/comments/add", comments.AddComment).Methods("POST")
}

func addClientFilesHandlers(r *mux.Router) {
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./client")))
}

func main() {
	r := mux.NewRouter()
	addAPIHandlers("/api", r)
	addClientFilesHandlers(r)

	if config.Debug {
		corsRouter := gh.CORS(
			gh.AllowedOrigins([]string{"*"}),
			gh.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			gh.AllowedHeaders([]string{"Content-Type", "Authorization"}),
			gh.ExposedHeaders([]string{"Authorization"}),
		)(r)
		http.ListenAndServe(":8090", gh.LoggingHandler(os.Stdout, corsRouter))
	} else {
		http.ListenAndServe(":80", gh.LoggingHandler(os.Stdout, r))
	}
}
