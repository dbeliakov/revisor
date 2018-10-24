package comments

import (
	"net/http"
	"reviewer/api/auth"
	"reviewer/api/store"
	"reviewer/api/utils"
	"time"

	"github.com/sirupsen/logrus"
)

// AddComment to line of review
var AddComment = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	var form struct {
		Text     string `json:"text" validate:"required"`
		Parent   string `json:"parent"`
		ReviewID string `json:"review_id" validate:"required"`
		LineID   string `json:"line_id" validate:"required"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	comment := store.Comment{
		Author:  user.Login,
		Created: time.Now().Unix(),
		Text:    form.Text,
		ParetID: form.Parent,
		LineID:  form.LineID,
	}

	if len(form.Parent) > 0 {
		if exists, err := store.Comments.CheckExists(form.ReviewID, comment.ParetID); err != nil || !exists {
			logrus.Warnf("Cannot find parent comment: %+v", err)
			utils.Error(w, utils.JSONErrorResponse{
				Status:        http.StatusBadRequest,
				Message:       "Cannot find parent comment",
				ClientMessage: "Сервер получил некорректный запрос",
			})
			return
		}
	}
	err = store.Comments.CreateComment(form.ReviewID, &comment)
	if err != nil {
		logrus.Errorf("Cannot save new comment: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot save comment to database"))
		return
	}

	utils.Ok(w, nil)
})
