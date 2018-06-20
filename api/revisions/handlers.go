package revisions

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	auth "reviewer/api/auth/database"
	"reviewer/api/auth/middlewares"
	comments "reviewer/api/comments/database"
	"reviewer/api/revisions/database"
	"reviewer/api/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/globalsign/mgo/bson"
)

func idsByLogins(logins []string) ([]string, error) {
	var reviewers []string
	for _, reviewer := range logins {
		reviewer := strings.TrimSpace(reviewer)
		user, err := auth.UserByLogin(reviewer)
		if err != nil {
			return nil, errors.New("No such reviewer: " + reviewer)
		}
		reviewers = append(reviewers, user.ID.Hex())
	}
	return reviewers, nil
}

func stringsToBsonIDs(strings []string) []bson.ObjectId {
	var result []bson.ObjectId
	for _, str := range strings {
		result = append(result, bson.ObjectIdHex(str))
	}
	return result
}

func hasAccess(id string, review database.Review) bool {
	hasAccess := review.OwnerID.Hex() == id
	if !hasAccess {
		for _, reviewer := range review.ReviewersID {
			if reviewer.Hex() == id {
				hasAccess = true
				break
			}
		}
	}
	return hasAccess
}

// OutgoingReviews returns reviews
var OutgoingReviews = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	reviews, err := database.ReviewsByConditions(bson.M{"ownerid": bson.ObjectIdHex(id.(string))})
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot load outgoing reviews")
		return
	}
	utils.Ok(w, reviews)
})

// IncomingReviews return reviews
var IncomingReviews = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	reviews, err := database.ReviewsByConditions(bson.M{"reviewersid": bson.ObjectIdHex(id.(string))})
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot load outgoing reviews")
		return
	}
	utils.Ok(w, reviews)
})

// NewReview creates new review
var NewReview = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	var form struct {
		Name        string `json:"name" validate:"required"`
		Reviewers   string `json:"reviewers" validate:"required"`
		FileName    string `json:"file_name" validate:"required"`
		FileContent string `json:"file_content" validate:"required"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	fileContent, err := base64.StdEncoding.DecodeString(form.FileContent)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Incorrect base64 file content")
		return
	}

	var reviewers []string
	if reviewers, err = idsByLogins(strings.Split(form.Reviewers, ",")); err != nil {
		utils.Error(w, http.StatusNotAcceptable, err.Error())
		return
	}

	review := database.NewReview(
		form.FileName,
		difflib.SplitLines(string(fileContent)),
		form.Name, id.(string),
		reviewers)

	err = review.Save()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot save review")
	}
})

// Review returns information about review
var Review = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "No review with id: "+reviewID)
		return
	}

	getParamOr := func(param string, def int) (int, error) {
		value := r.URL.Query()[param]
		if len(value) > 0 {
			return strconv.Atoi(value[0])
		}
		return def, nil
	}

	startRev, err := getParamOr("start_rev", 0)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Incorrect start revision")
		return
	}
	endRev, err := getParamOr("end_rev", review.File.RevisionsCount()-1)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Incorrect end revision")
		return
	}

	if !hasAccess(id.(string), review) {
		utils.Error(w, http.StatusForbidden, "No access to this review")
		return
	}

	content, err := review.File.Diff(startRev, endRev)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot calculate diff")
		return
	}

	comments, err := comments.RootCommentsForReview(review.ID.Hex())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot calculate diff")
		return
	}

	utils.Ok(w, &map[string]interface{}{
		"info":     review,
		"diff":     content,
		"comments": comments,
	})
})

// UpdateReview information or add revision
var UpdateReview = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "No review with id: "+reviewID)
		return
	}
	if review.OwnerID.Hex() != id {
		utils.Error(w, http.StatusForbidden, "Only owner allowed")
		return
	}

	var form struct {
		Name        string `json:"name" validate:"required"`
		Reviewers   string `json:"reviewers" validate:"required"`
		NewRevision string `json:"new_revision"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	fileContent, err := base64.StdEncoding.DecodeString(form.NewRevision)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Incorrect base64 file content")
		return
	}
	review.Name = form.Name

	var reviewers []string
	if reviewers, err = idsByLogins(strings.Split(form.Reviewers, ",")); err != nil {
		utils.Error(w, http.StatusNotAcceptable, err.Error())
		return
	}
	review.ReviewersID = stringsToBsonIDs(reviewers)

	if len(form.NewRevision) > 0 {
		err = review.File.AddRevision(difflib.SplitLines(string(fileContent)))
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Cannot add revision")
			return
		}
		review.Updated = time.Now()
	}
	err = review.Save()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot update reiew")
		return
	}
	utils.Ok(w, nil)
})

// Decline review
var Decline = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "No review with id: "+reviewID)
		return
	}
	if !hasAccess(id.(string), review) {
		utils.Error(w, http.StatusForbidden, "No access to this review")
		return
	}

	review.Closed = true
	review.Accepted = false
	review.Updated = time.Now()
	err = review.Save()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot update reiew")
		return
	}
})

// Accept review
var Accept = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "No review with id: "+reviewID)
		return
	}
	if review.OwnerID.Hex() == id || !hasAccess(id.(string), review) {
		utils.Error(w, http.StatusForbidden, "No access to this review")
		return
	}

	review.Closed = true
	review.Accepted = true
	review.Updated = time.Now()
	err = review.Save()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot update reiew")
		return
	}
})
