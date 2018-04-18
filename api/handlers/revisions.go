package handlers

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"reviewer/api/database"
	"reviewer/api/middlewares"
	"reviewer/api/utils"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/globalsign/mgo/bson"
)

type userInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func writeReviews(w http.ResponseWriter, reviews []database.Review) {
	type reviewInfo struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		Owner     userInfo   `json:"owner"`
		Reviewers []userInfo `json:"reviewers"`
		Updated   int64      `json:"updated"`
		Closed    bool       `json:"closed"`
		Accepted  bool       `json:"accepted"`
	}

	var info = make([]reviewInfo, 0)
	for _, review := range reviews {
		i := reviewInfo{
			ID:      review.ID.Hex(),
			Name:    review.Name,
			Updated: review.Updated.Unix(),
		}
		o, err := database.UserByID(review.Owner.Hex())
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Cannot load info")
			return
		}
		i.Owner = userInfo{
			FirstName: o.FirstName,
			LastName:  o.LastName,
			Username:  o.Login,
		}
		for _, reviewer := range review.Reviewers {
			r, err := database.UserByID(reviewer.Hex())
			if err != nil {
				utils.Error(w, http.StatusInternalServerError, "Cannot load info")
				return
			}
			i.Reviewers = append(i.Reviewers, userInfo{
				FirstName: r.FirstName,
				LastName:  r.LastName,
				Username:  r.Login,
			})
		}
		info = append(info, i)
	}
	utils.Ok(w, info)
}

// OutgoingReviews returns reviews
var OutgoingReviews = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	reviews, err := database.FindReviews(bson.M{"owner": bson.ObjectIdHex(id.(string))})
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot load outgoing reviews")
		return
	}
	writeReviews(w, reviews)
})

// IncomingReviews return reviews
var IncomingReviews = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	reviews, err := database.FindReviews(bson.M{"reviewers": bson.ObjectIdHex(id.(string))})
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot load outgoing reviews")
		return
	}
	writeReviews(w, reviews)
})

// NewReview creates new review
var NewReview = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	type newReviewForm struct {
		Name        string `json:"name"`
		Reviewers   string `json:"reviewers"`
		FileName    string `json:"file_name"`
		FileContent string `json:"file_content"`
	}

	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	var form newReviewForm
	err := utils.UnmarshalForm(r, &form)
	if err == utils.ErrIncorrectBody {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	} else if err == utils.ErrIncorectFormFields {
		utils.Error(w, http.StatusNotAcceptable, err.Error())
		return
	} else if err != nil {
		log.Panic("Assert: unexpected error: " + err.Error())
	}

	fileContent, err := base64.StdEncoding.DecodeString(form.FileContent)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Incorrect base64 file content")
		return
	}

	var reviewers []string
	for _, reviewer := range strings.Split(form.Reviewers, ",") {
		reviewer := strings.TrimSpace(reviewer)
		user, err := database.UserByLogin(reviewer)
		if err != nil {
			utils.Error(w, http.StatusNotAcceptable, "No such reviewer: "+reviewer)
			return
		}
		reviewers = append(reviewers, user.ID.Hex())
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
	review, err := database.LoadReview(reviewID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "No review with id: "+reviewID)
		return
	}

	startRevParam := r.URL.Query()["start_rev"]
	var startRev int
	if len(startRevParam) == 0 {
		startRev = 0
	} else {
		startRev, err = strconv.Atoi(startRevParam[0])
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Incorrect start revision")
		}
	}

	endRevParam := r.URL.Query()["end_rev"]
	var endRev int
	if len(endRevParam) == 0 {
		endRev = review.File.RevisionsCount() - 1
	} else {
		endRev, err = strconv.Atoi(endRevParam[0])
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Incorrect end revision")
			return
		}
	}

	hasAccess := review.Owner.Hex() == id
	if !hasAccess {
		for _, reviewer := range review.Reviewers {
			if reviewer.Hex() == id {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		utils.Error(w, http.StatusForbidden, "No access to this review")
		return
	}

	type diffRange struct {
		From int
		To   int
	}

	type lineInfo struct {
		ID       string `json:"id"`
		Revision int    `json:"revision"`
		Content  string `json:"content"`
	}

	type diffGroup struct {
		OldRange diffRange
		NewRange diffRange
		Lines    lineInfo
	}

	type reviewContent struct {
		ID        string     `json:"id"`
		Content   []lineInfo `json:"content"`
		Filename  string
		Name      string     `json:"name"`
		Owner     userInfo   `json:"user_info"`
		Reviewers []userInfo `json:"reviewers"`
		Updated   int64      `json:"updated"`
		Closed    bool       `json:"closed"`
		Accepted  bool       `json:"accepted"`
	}

	content, err := review.File.Diff(startRev, endRev)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot calculate diff")
		return
	}

	result := reviewContent{
		ID:       review.ID.Hex(),
		Name:     review.Name,
		Updated:  review.Updated.Unix(),
		Closed:   review.Closed,
		Accepted: review.Accepted,
	}

	owner, err := database.UserByID(id.(string))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot read owner info")
		return
	}
	result.Owner = userInfo{
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Username:  owner.Login,
	}
	for _, reviewerID := range review.Reviewers {
		reviewer, err := database.UserByID(reviewerID.Hex())
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Cannot read owner info")
			return
		}
		result.Reviewers = append(result.Reviewers, userInfo{
			FirstName: reviewer.FirstName,
			LastName:  reviewer.LastName,
			Username:  reviewer.Login,
		})
	}

	for _, _ = range content.Groups {
		// TODO
		/*group := diffGroup{
			OldRange: diffRange{
				From: gr.OldRange.From,
				To:   gr.OldRange.To,
			},
			NewRange: diffRange{
				From: gr.NewRange.From,
				To:   gr.NewRange.To,
			},
		}
		for _, l := range gr.Lines {
			group.Lines = append(group.Lines, lineInfo{
				ID: l.
			})
		}*/
	}
})
