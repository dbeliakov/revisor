package main

import (
	"net/http"
	"os"

	"github.com/dbeliakov/revisor/api/auth"
	"github.com/dbeliakov/revisor/api/comments"
	"github.com/dbeliakov/revisor/api/config"
	"github.com/dbeliakov/revisor/api/review"
	"github.com/dbeliakov/revisor/api/store"
	"github.com/go-chi/chi"
	gh "github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func addAuthHandlers(base string, logger zap.Logger, service auth.Service, r *chi.Mux) {
	authRouter := r.Route("/auth", nil)
	authRouter.Use(auth.NewRequiredMiddleware(logger))
	authRouter.Post("/login", service.HandleLogin)
	authRouter.Post("/signup", service.HandleSignUp)
	authRouter.Post("/user", service.HandleUserInfo)
	// TODO replace with other URL
	authRouter.Post("change/password", service.HandleChangePassword)
}

func addAPIHandlers(base string, r *chi.Mux) {
	// Review handlers
	r.Get(base+"/reviews/outgoing", review.OutgoingReviews)
	r.Get(base+"/reviews/incoming", review.IncomingReviews)
	r.Post(base+"/reviews/new", review.NewReview)
	r.Get(base+"/reviews/{id}", review.Review)
	r.Post(base+"/reviews/{id}/update", review.UpdateReview)
	r.Post(base+"/reviews/{id}/accept", review.Accept)
	r.Post(base+"/reviews/{id}/decline", review.Decline)
	r.Get(base+"/users/search", review.SearchReviewer)

	// Comments handlers
	r.Post(base+"/comments/add", comments.AddComment)
}

func addClientFilesHandlers(r *chi.Mux) {
	// TODO check if it is correct
	r.Handle("/*", http.FileServer(http.Dir(config.ClientFilesDir)))
}

func main() {
	store.InitStore()

	r := chi.NewRouter()
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
