package revisions

import (
	"encoding/base64"
	"errors"
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
	"github.com/sirupsen/logrus"

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
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	reviews, err := database.ReviewsByConditions(bson.M{"ownerid": bson.ObjectIdHex(user.ID.Hex())})
	if err != nil {
		logrus.Errorf("Cannot load outgoing reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load outgoing reviews"))
		return
	}
	utils.Ok(w, reviews)
})

// IncomingReviews return reviews
var IncomingReviews = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	reviews, err := database.ReviewsByConditions(bson.M{"reviewersid": bson.ObjectIdHex(user.ID.Hex())})
	if err != nil {
		logrus.Errorf("Cannot load outgoing reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load incoming reviews"))
		return
	}
	utils.Ok(w, reviews)
})

// NewReview creates new review
var NewReview = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
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
		logrus.Warnf("Incorrect base64 file content: %s, error: %+v", form.FileContent, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Incorrect base64 file content",
			ClientMessage: "Некорректная кодировка файла",
		})
		return
	}

	var reviewers []string
	if reviewers, err = idsByLogins(strings.Split(form.Reviewers, ",")); err != nil {
		logrus.Infof("Incorrect list of reviewers: %+v", err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Incorrect list of reviewers",
			ClientMessage: "Некорректный список ревьюеров",
		})
		return
	}

	review := database.NewReview(
		form.FileName,
		difflib.SplitLines(string(fileContent)),
		form.Name, user.ID.Hex(),
		reviewers)

	err = review.Save()
	if err != nil {
		logrus.Errorf("Cannot save new review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot save review"))
	}
	utils.Ok(w, nil)
})

// Review returns information about review
var Review = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review with id: %+s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
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
		logrus.Warnf("Incorrect start revision: %+v", err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Incorrect start revision",
			ClientMessage: "Некорректный номер начальной ревизии",
		})
		return
	}
	endRev, err := getParamOr("end_rev", review.File.RevisionsCount()-1)
	if err != nil {
		logrus.Warnf("Incorrect end revision: %+v", err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Incorrect end revision",
			ClientMessage: "Некорректный номер конечной ревизии",
		})
		return
	}

	if !hasAccess(user.ID.Hex(), review) {
		logrus.Warnf("User %s han no access to review %s", user.Login, review.ID.Hex())
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "No access to this review",
			ClientMessage: "У вас недостаточно прав для просмотра ревью",
		})
		return
	}

	content, err := review.File.Diff(startRev, endRev)
	if err != nil {
		logrus.Errorf("Cannot calculate diff: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot calculate diff"))
		return
	}

	comments, err := comments.RootCommentsForReview(review.ID.Hex())
	if err != nil {
		logrus.Errorf("Cannot load comments for review: %s, error: %+v", review.ID.Hex(), err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load comments"))
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
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review: %s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	if review.OwnerID.Hex() != user.ID.Hex() {
		logrus.Warnf("User %s tries to update review %s without being the creator", user.Login, review.ID.Hex())
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "Only owner allowed to update review",
			ClientMessage: "Только автор ревью может его обновлять",
		})
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
		logrus.Warnf("Incorrect base64 file content: %s, error: %+v", form.NewRevision, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Incorrect base64 file content",
			ClientMessage: "Некорректная кодировка файла",
		})
		return
	}
	review.Name = form.Name

	var reviewers []string
	if reviewers, err = idsByLogins(strings.Split(form.Reviewers, ",")); err != nil {
		logrus.Warnf("Incorrect list of reviewers: %+v", err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Incorrect list of reviewers",
			ClientMessage: "Некорректный список ревьюеров",
		})
		return
	}
	review.ReviewersID = stringsToBsonIDs(reviewers)

	if len(form.NewRevision) > 0 {
		err = review.File.AddRevision(difflib.SplitLines(string(fileContent)))
		if err != nil {
			logrus.Errorf("Cannot add revision: %+v", err)
			utils.Error(w, utils.InternalErrorResponse("Cannot add revision"))
			return
		}
		review.Updated = time.Now()
	}
	err = review.Save()
	if err != nil {
		logrus.Errorf("Cannot save updated review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update review"))
		return
	}
	utils.Ok(w, nil)
})

// Decline review
var Decline = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review: %s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	if !hasAccess(user.ID.Hex(), review) {
		logrus.Warnf("User %s has no access to review %s", user.Login, review.ID.Hex())
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "No access to this review",
			ClientMessage: "У вас недостаточно прав для закрытия ревью",
		})
		return
	}

	review.Closed = true
	review.Accepted = false
	review.Updated = time.Now()
	err = review.Save()
	if err != nil {
		logrus.Errorf("Cannot save review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update review"))
		return
	}
})

// Accept review
var Accept = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := database.ReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review: %s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	if review.OwnerID.Hex() == user.ID.Hex() || !hasAccess(user.ID.Hex(), review) {
		logrus.Warnf("User %s has no access to review %s", user.Login, review.ID.Hex())
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "No access to this review",
			ClientMessage: "У вас недостаточно прав для закрытия ревью",
		})
		return
	}

	review.Closed = true
	review.Accepted = true
	review.Updated = time.Now()
	err = review.Save()
	if err != nil {
		logrus.Errorf("Cannot save review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update review"))
		return
	}
})
