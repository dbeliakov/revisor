package comments

import (
	"net/http"
	"reviewer/api/auth/middlewares"
	"reviewer/api/comments/database"
	"reviewer/api/utils"

	"github.com/sirupsen/logrus"
)

// AddComment to line of review
var AddComment = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromRequest(r)
	if err != nil {
		logrus.Error("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
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

	var parent *database.Comment = nil
	if form.Parent != "" {
		logrus.Info("Parent id: " + form.Parent)
		p, err := database.CommentByID(form.Parent)
		if err != nil {
			logrus.Warnf("Cannot find parent comment: %+v", err)
			utils.Error(w, utils.JSONErrorResponse{
				Status:        http.StatusBadRequest,
				Message:       "Cannot find parent comment",
				ClientMessage: "Сервер получил некорректный запрос",
			})
			return
		}
		parent = &p
	}

	comment := database.NewComment(user.ID.Hex(), form.Text, form.ReviewID, form.LineID, parent == nil)
	err = comment.Save()
	if err != nil {
		logrus.Errorf("Cannot save new comment: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot save comment to database"))
		return
	}

	if parent != nil {
		parent.ChildIDs = append(parent.ChildIDs, comment.ID)
		err = parent.Save()
	}
	if err != nil {
		logrus.Errorf("Cannot save parent comment: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot update parent comment"))
		return
	}

	utils.Ok(w, nil)
})
