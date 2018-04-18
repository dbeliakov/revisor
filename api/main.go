package main

import (
	"net/http"
	"os"
	"reviewer/api/handlers"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func newRouter(base string) *mux.Router {
	r := mux.NewRouter()

	// Auth handlers
	r.HandleFunc(base+"/auth/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc(base+"/auth/signup", handlers.SignUpHandler).Methods("POST")
	r.HandleFunc(base+"/auth/user", handlers.UserInfoHandler).Methods("GET")
	r.HandleFunc(base+"/auth/change/password", handlers.ChangePasswordHandler).Methods("POST")

	// Review handlers
	r.HandleFunc(base+"/reviews/outgoing", handlers.OutgoingReviews).Methods("GET")
	r.HandleFunc(base+"/reviews/incoming", handlers.IncomingReviews).Methods("GET")
	r.HandleFunc(base+"/reviews/new", handlers.NewReview).Methods("POST")

	return r
}

func main() {
	/*a, _ := ioutil.ReadFile("test/a.cpp")
	//b, _ := ioutil.ReadFile("test/b.cpp")
	rw := database.NewReview("test.cpp", difflib.SplitLines(string(a)), "Тестовое задание", "5acdd916ed9b7ebab904f03b", []string{"5acdd916ed9b7ebab904f03b"})
	rw.Save()*/

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
