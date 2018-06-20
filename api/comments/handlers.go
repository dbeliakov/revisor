package comments

import (
	"errors"
	"log"
	"net/http"
	"reviewer/api/auth/middlewares"
	"reviewer/api/comments/database"
	"reviewer/api/utils"
)

// AddComment to line of review
var AddComment = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
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

	var parent *database.Comment = nil
	if form.Parent != "" {
		p, err := database.CommentByID(form.Parent)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Cannot find parent comment")
			return
		}
		parent = &p
	}

	comment := database.NewComment(id.(string), form.Text, form.ReviewID, form.LineID, parent == nil)
	err := comment.Save()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot save comment to database")
		return
	}

	if parent != nil {
		parent.ChildIDs = append(parent.ChildIDs, comment.ID)
		err = parent.Save()
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot update parent comment")
		return
	}

	utils.Ok(w, nil)
})
