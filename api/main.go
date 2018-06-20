package main

import (
	"net/http"
	"os"
	"reviewer/api/auth"
	"reviewer/api/comments"
	"reviewer/api/revisions"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func newRouter(base string) *mux.Router {
	r := mux.NewRouter()

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

	return r
}

func main() {
	r := newRouter("/api")

	// DEBUG
	corsRouter := gh.CORS(
		gh.AllowedOrigins([]string{"*"}),
		gh.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		gh.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gh.ExposedHeaders([]string{"Authorization"}),
	)(r)

	// TODO logging to file
	http.ListenAndServe(":8090", gh.LoggingHandler(os.Stdout, corsRouter))
}
