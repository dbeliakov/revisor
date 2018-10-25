package main

import (
	"net/http"
	"os"
	"reviewer/api/auth"
	"reviewer/api/comments"
	"reviewer/api/config"
	"reviewer/api/review"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func addAPIHandlers(base string, r *mux.Router) {
	// Auth handlers
	r.HandleFunc(base+"/auth/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc(base+"/auth/signup", auth.SignUpHandler).Methods("POST")
	r.HandleFunc(base+"/auth/user", auth.UserInfoHandler).Methods("GET")
	r.HandleFunc(base+"/auth/change/password", auth.ChangePasswordHandler).Methods("POST")

	// Review handlers
	r.HandleFunc(base+"/reviews/outgoing", review.OutgoingReviews).Methods("GET")
	r.HandleFunc(base+"/reviews/incoming", review.IncomingReviews).Methods("GET")
	r.HandleFunc(base+"/reviews/new", review.NewReview).Methods("POST")
	r.HandleFunc(base+"/reviews/{id}", review.Review).Methods("GET")
	r.HandleFunc(base+"/reviews/{id}/update", review.UpdateReview).Methods("POST")
	r.HandleFunc(base+"/reviews/{id}/accept", review.Accept).Methods("GET")
	r.HandleFunc(base+"/reviews/{id}/decline", review.Decline).Methods("GET")
	r.HandleFunc(base+"/users/search", review.SearchReviewer).Methods("GET")

	// Comments handlers
	r.HandleFunc(base+"/comments/add", comments.AddComment).Methods("POST")
}

func addClientFilesHandlers(r *mux.Router) {
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.ClientFilesDir)))
}

func main() {
	r := mux.NewRouter()
	addAPIHandlers("/api", r)
	addClientFilesHandlers(r)

	handler := gh.CombinedLoggingHandler(os.Stdout, r)
	listenAddres := "[::]:80"

	if config.Debug {
		handler = gh.CORS(
			gh.AllowedOrigins([]string{"*"}),
			gh.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			gh.AllowedHeaders([]string{"Content-Type", "Authorization"}),
			gh.ExposedHeaders([]string{"Authorization"}),
		)(handler)
		listenAddres = "[::]:8090"
	}
	err := http.ListenAndServe(listenAddres, handler)
	if err != nil {
		logrus.Panicf("Error from http server: %+v", err)
	}
}
