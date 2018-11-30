package review

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reviewer/api/auth"
	"reviewer/api/store"
	"reviewer/api/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
)

func hasAccess(login string, review store.Review) bool {
	hasAccess := review.Owner == login
	if !hasAccess {
		for _, reviewer := range review.Reviewers {
			if reviewer == login {
				hasAccess = true
				break
			}
		}
	}
	return hasAccess
}

// APIReview represents api result struct
type APIReview struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Updated        int64          `json:"updated"`
	Closed         bool           `json:"closed"`
	Accepted       bool           `json:"accepted"`
	Owner          auth.APIUser   `json:"owner"`
	Reviewers      []auth.APIUser `json:"reviewers"`
	RevisionsCount int            `json:"revisions_count"`
	CommentsCount  int            `json:"comments_count"`
}

// NewAPIReview creates new api review from store review
func NewAPIReview(review store.Review) (APIReview, error) {
	result := APIReview{
		ID:       review.ID,
		Name:     review.Name,
		Updated:  review.Updated,
		Closed:   review.Closed,
		Accepted: review.Accepted,
	}
	owner, err := store.Auth.FindUserByLogin(review.Owner)
	if err != nil {
		return APIReview{}, err
	}
	result.Owner = auth.NewAPIUser(owner)
	result.Reviewers = make([]auth.APIUser, 0)
	for _, login := range review.Reviewers {
		reviewer, err := store.Auth.FindUserByLogin(login)
		if err != nil {
			return APIReview{}, err
		}
		result.Reviewers = append(result.Reviewers, auth.NewAPIUser(reviewer))
	}
	var p fastjson.Parser
	v, err := p.ParseBytes(review.File)
	if err != nil {
		return result, err
	}
	result.RevisionsCount = len(v.GetArray("Revisions"))
	comments, err := store.Comments.CommentsForReview(review.ID)
	if err != nil {
		return result, err
	}
	result.CommentsCount = len(comments)
	return result, nil
}

func newAPIReviews(reviews []store.Review) ([]APIReview, error) {
	result := make([]APIReview, 0)
	for _, r := range reviews {
		ar, err := NewAPIReview(r)
		if err != nil {
			return nil, err
		}
		result = append(result, ar)
	}
	return result, nil
}

// APIComment represents api result struct
type APIComment struct {
	ID      string        `json:"id"`
	Author  auth.APIUser  `json:"author"`
	Created int64         `json:"created"`
	Text    string        `json:"text"`
	LineID  string        `json:"line_id"`
	Childs  []*APIComment `json:"childs"`
}

// NewAPIComment creates new api comment from store comment
func NewAPIComment(comment store.Comment) (APIComment, error) {
	result := APIComment{
		ID:      comment.ID,
		Created: comment.Created,
		Text:    comment.Text,
		LineID:  comment.LineID,
		Childs:  make([]*APIComment, 0),
	}
	author, err := store.Auth.FindUserByLogin(comment.Author)
	if err != nil {
		return result, err
	}
	result.Author = auth.NewAPIUser(author)
	return result, nil
}

// OutgoingReviews returns reviews
var OutgoingReviews = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	reviews, err := store.Reviews.FindReviewsByOwner(user.Login)
	if err != nil {
		logrus.Errorf("Cannot load outgoing reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load outgoing reviews"))
		return
	}
	res, err := newAPIReviews(reviews)
	if err != nil {
		logrus.Errorf("Cannot load users from outgoing reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load outgoing reviews"))
		return
	}
	utils.Ok(w, res)
})

// IncomingReviews return reviews
var IncomingReviews = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	reviews, err := store.Reviews.FindReviewsByReviewer(user.Login)
	if err != nil {
		logrus.Errorf("Cannot load outgoing reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load incoming reviews"))
		return
	}
	res, err := newAPIReviews(reviews)
	if err != nil {
		logrus.Errorf("Cannot load users from incoming reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load incoming reviews"))
		return
	}
	utils.Ok(w, res)
})

// NewReview creates new review
var NewReview = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
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
		logrus.Warnf("Incorrect base64 file content: %s, error: %+v", form.FileContent, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Incorrect base64 file content",
			ClientMessage: "Некорректная кодировка файла",
		})
		return
	}

	reviewers := strings.Split(form.Reviewers, ",")
	for _, r := range reviewers {
		exists, err := store.Auth.CheckExists(r)
		if !exists || err != nil {
			logrus.Infof("Incorrect list of reviewers: %+v", err)
			utils.Error(w, utils.JSONErrorResponse{
				Status:        http.StatusNotAcceptable,
				Message:       "Incorrect list of reviewers",
				ClientMessage: "Некорректный список ревьюеров",
			})
			return
		}
	}

	file := NewVersionedFile(form.FileName, difflib.SplitLines(string(fileContent)))
	bytesFile, err := json.Marshal(&file)
	if err != nil {
		logrus.Errorf("Cannot serialize versioned file: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot create versioned file"))
		return
	}
	review := store.Review{
		File:      bytesFile,
		Name:      form.Name,
		Owner:     user.Login,
		Reviewers: reviewers,
		Updated:   time.Now().Unix(),
		Closed:    false,
		Accepted:  false,
	}

	err = store.Reviews.CreateReview(&review)
	if err != nil {
		logrus.Errorf("Cannot save new review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot save review"))
	}
	utils.Ok(w, nil)
})

// Review returns information about review
var Review = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := store.Reviews.FindReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review with id: %+s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	var file VersionedFile
	err = json.Unmarshal(review.File, &file)
	if err != nil {
		logrus.Errorf("Error while deserializing versioned file: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot deserialize versioned file"))
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
	endRev, err := getParamOr("end_rev", file.RevisionsCount()-1)
	if err != nil {
		logrus.Warnf("Incorrect end revision: %+v", err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Incorrect end revision",
			ClientMessage: "Некорректный номер конечной ревизии",
		})
		return
	}

	if !hasAccess(user.Login, review) {
		logrus.Warnf("User %s han no access to review %s", user.Login, review.ID)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "No access to this review",
			ClientMessage: "У вас недостаточно прав для просмотра ревью",
		})
		return
	}

	content, err := file.Diff(startRev, endRev)
	if err != nil {
		logrus.Errorf("Cannot calculate diff: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot calculate diff"))
		return
	}

	comments, err := store.Comments.CommentsForReview(review.ID)
	if err != nil {
		logrus.Errorf("Cannot load comments for review: %s, error: %+v", review.ID, err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load comments"))
		return
	}
	apiComments := make(map[string]*APIComment)
	for _, comment := range comments {
		ac, err := NewAPIComment(comment)
		if err != nil {
			logrus.Errorf("Cannot create API comment for review: %s, error: %+v", review.ID, err)
			utils.Error(w, utils.InternalErrorResponse("Cannot load comments"))
			return
		}
		if len(comment.ParentID) == 0 {
			apiComments[comment.ID] = &ac
		} else {
			parent, exists := apiComments[comment.ParentID]
			if !exists {
				logrus.Errorf("Cannot find parent comment for review: %s, error: %+v", review.ID, err)
				utils.Error(w, utils.InternalErrorResponse("Cannot load comments"))
				return
			}
			parent.Childs = append(parent.Childs, &ac)
			apiComments[comment.ParentID] = parent
			apiComments[comment.ID] = &ac
		}
	}
	resComments := make([]APIComment, 0)
	for _, comment := range comments {
		if len(comment.ParentID) == 0 {
			resComments = append(resComments, *apiComments[comment.ID])
		}
	}

	res, err := NewAPIReview(review)
	if err != nil {
		logrus.Errorf("Cannot load users from incoming reviews: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot load incoming reviews"))
		return
	}
	utils.Ok(w, &map[string]interface{}{
		"info":     res,
		"diff":     content,
		"comments": resComments,
	})
})

// UpdateReview information or add revision
var UpdateReview = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := store.Reviews.FindReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review: %s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	if review.Owner != user.Login {
		logrus.Warnf("User %s tries to update review %s without being the creator", user.Login, review.ID)
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

	reviewers := strings.Split(form.Reviewers, ",")
	for _, r := range reviewers {
		exists, err := store.Auth.CheckExists(r)
		if !exists || err != nil {
			logrus.Infof("Incorrect list of reviewers: %+v", err)
			utils.Error(w, utils.JSONErrorResponse{
				Status:        http.StatusNotAcceptable,
				Message:       "Incorrect list of reviewers",
				ClientMessage: "Некорректный список ревьюеров",
			})
			return
		}
	}
	review.Reviewers = reviewers

	var file VersionedFile
	err = json.Unmarshal(review.File, &file)
	if err != nil {
		logrus.Errorf("Error while deserializing versioned file: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot deserialize versioned file"))
		return
	}
	if len(form.NewRevision) > 0 {
		err = file.AddRevision(difflib.SplitLines(string(fileContent)))
		if err != nil {
			logrus.Errorf("Cannot add revision: %+v", err)
			utils.Error(w, utils.InternalErrorResponse("Cannot add revision"))
			return
		}
		review.Updated = time.Now().Unix()
	}
	bytesFile, err := json.Marshal(&file)
	if err != nil {
		logrus.Errorf("Cannot serialize versioned file: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot create versioned file"))
		return
	}
	review.File = bytesFile
	err = store.Reviews.UpdateReview(review)
	if err != nil {
		logrus.Errorf("Cannot save updated review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update review"))
		return
	}
	utils.Ok(w, nil)
})

// Decline review
var Decline = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := store.Reviews.FindReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review: %s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	if !hasAccess(user.Login, review) {
		logrus.Warnf("User %s has no access to review %s", user.Login, review.ID)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "No access to this review",
			ClientMessage: "У вас недостаточно прав для закрытия ревью",
		})
		return
	}

	review.Closed = true
	review.Accepted = false
	err = store.Reviews.UpdateReview(review)
	if err != nil {
		logrus.Errorf("Cannot save review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update review"))
		return
	}
})

// Accept review
var Accept = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	vars := mux.Vars(r)
	reviewID := vars["id"]
	review, err := store.Reviews.FindReviewByID(reviewID)
	if err != nil {
		logrus.Warnf("Cannot find review: %s, error: %+v", reviewID, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotFound,
			Message:       "No review with id: " + reviewID,
			ClientMessage: "Не удалось найти ревью",
		})
		return
	}
	if review.Owner == user.Login || !hasAccess(user.Login, review) {
		logrus.Warnf("User %s has no access to review %s", user.Login, review.ID)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusForbidden,
			Message:       "No access to this review",
			ClientMessage: "У вас недостаточно прав для закрытия ревью",
		})
		return
	}

	review.Closed = true
	review.Accepted = true
	err = store.Reviews.UpdateReview(review)
	if err != nil {
		logrus.Errorf("Cannot save review: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update review"))
		return
	}
})

// SearchReviewer by login or by name
var SearchReviewer = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
	}

	query, ok := r.URL.Query()["query"]
	if !ok || len(query) == 0 || len(query[0]) == 0 {
		logrus.Warnf("Empty query for search")
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusBadRequest,
			Message:       "Empty Query",
			ClientMessage: "Пустой запрос для поиска",
		})
		return
	}

	results, err := store.Auth.FindUsers(query[0], user.Login)
	if err != nil {
		logrus.Errorf("Cannot find reviewers: %+v", err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusInternalServerError,
			Message:       "Cannot find users",
			ClientMessage: "Не удалось произвести поиск пользователей",
		})
		return
	}
	res := make([]auth.APIUser, 0)
	for _, r := range results {
		res = append(res, auth.NewAPIUser(r))
	}
	utils.Ok(w, res)
})
