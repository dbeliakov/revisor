package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// APIResponse struct to prevent incorrect data in response
type APIResponse interface {
	Check()
}

// JSONErrorResponse presents general structure of error response
type JSONErrorResponse struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	ClientMessage string `json:"client_message"`
}

// InternalErrorResponse template with custom message
func InternalErrorResponse(message string) JSONErrorResponse {
	return JSONErrorResponse{
		Status:        http.StatusInternalServerError,
		Message:       message,
		ClientMessage: "Внутренняя ошибка сервера. Пожалуйста, попробуйте повторить позднее",
	}
}

// Unauthorized writes unauthorized response
func Unauthorized(w http.ResponseWriter) {
	bytes, err := json.Marshal(JSONErrorResponse{
		Status:        http.StatusUnauthorized,
		Message:       "Unauthorized",
		ClientMessage: "Для продолжения работы необходимо авторизоваться",
	})
	if err != nil {
		logrus.Errorf("Cannot create json for unauthorized message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write(bytes)
}

// Error writes error response
func Error(w http.ResponseWriter, response JSONErrorResponse) {
	bytes, err := json.Marshal(response)
	if err != nil {
		logrus.Errorf("Cannot create json for unauthorized message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(response.Status)
	w.Write(bytes)
}

// Ok writes error response
func Ok(w http.ResponseWriter, data interface{}) {
	result := map[string]interface{}{
		"status": http.StatusOK,
	}
	if data != nil {
		result["data"] = data
	}
	bytes, err := json.Marshal(&result)
	if err != nil {
		logrus.Errorf("Cannot create json for unauthorized message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

var (
	// ErrIncorrectBody error
	ErrIncorrectBody = errors.New("Incorrect body")
	// ErrIncorectFormFields error
	ErrIncorectFormFields = errors.New("Incorrect form fields")
)

// UnmarshalForm from body and validate it. If error occurs, writes error message to response writer
func UnmarshalForm(w http.ResponseWriter, r *http.Request, to interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Cannot read request body: %+v", err)
		Error(w, JSONErrorResponse{
			Status:        http.StatusInternalServerError,
			Message:       "Cannot read request body",
			ClientMessage: "Внутренняя ошибка сервера. Пожалуйста, повторите позднее",
		})
		return ErrIncorrectBody
	}
	err = json.Unmarshal(body, to)
	if err != nil {
		logrus.Errorf("Cannot unmarshal request body: %+v", err)
		Error(w, JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Cannot unmarshal request body",
			ClientMessage: "Сервер получил некорректный запрос",
		})
		return ErrIncorrectBody
	}

	validate := validator.New()
	err = validate.Struct(to)
	if err != nil {
		logrus.Errorf("Invalid request body: %+v", err)
		Error(w, JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Invalid request",
			ClientMessage: "Сервер получил некорректный запрос",
		})
		return ErrIncorrectBody
	}
	return nil
}
