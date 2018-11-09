package comments

import (
	"net/http"
	"reviewer/api/auth"
	"reviewer/api/config"
	"reviewer/api/notifications"
	"reviewer/api/store"
	"reviewer/api/utils"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	frontendReviewPrefix = config.BaseURL + "/review/"
)

func notifyForComment(reviewID string, user store.User, comment string) {
	rev, err := store.Reviews.FindReviewByID(reviewID)
	if err != nil {
		logrus.Errorf("Cannot find review for notification: %+v", err)
		return
	}

	text := "<b>" + rev.Name + "</b>\n\n"
	text += user.FirstName + " " + user.LastName + " опубликовал комментарий:\n"
	text += "<code>" + comment + "</code>"

	if user.Login != rev.Owner {
		u, err := store.Auth.FindUserByLogin(rev.Owner)
		if err == nil {
			notifications.TelegramSend(u, text)
		}
	}

	for _, login := range rev.Reviewers {
		if user.Login == login {
			continue
		}
		u, err := store.Auth.FindUserByLogin(login)
		if err != nil {
			continue
		}
		notifications.TelegramSend(u, text)
	}
}

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
		Author:   user.Login,
		Created:  time.Now().Unix(),
		Text:     form.Text,
		ParentID: form.Parent,
		LineID:   form.LineID,
	}

	if len(form.Parent) > 0 {
		if exists, err := store.Comments.CheckExists(form.ReviewID, comment.ParentID); err != nil || !exists {
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

	notifyForComment(form.ReviewID, user, comment.Text)

	utils.Ok(w, nil)
})
