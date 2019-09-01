package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	contentTypeHeader   = "Content-Type"
	contentLengthHeader = "Content-Length"
)

func writeJSONContentType(rw http.ResponseWriter) {
	rw.Header().Set(contentTypeHeader, "application/json")
}

func writeContentLength(rw http.ResponseWriter, length int) {
	rw.Header().Set(contentLengthHeader, strconv.Itoa(length))
}

// JSONErrorResponse presents general structure of error response
type jsonErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"client_message"`
}

// JSONErrorResponse is alias for jsonErrorResponse
// TODO use jsonErrorResponse instead
type JSONErrorResponse jsonErrorResponse

func writeErrorResponse(rw http.ResponseWriter, response jsonErrorResponse) error {
	writeJSONContentType(rw)
	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		return xerrors.Errorf("cannot marshal json response: %w", err)
	}

	writeContentLength(rw, len(marshaledResponse))
	rw.WriteHeader(http.StatusBadRequest)
	_, err = rw.Write(marshaledResponse)
	if err != nil {
		return err
	}
	return nil
}

// HTTPError writes error response to response writer
func HTTPError(rw http.ResponseWriter, status int, message string) error {
	err := writeErrorResponse(rw, jsonErrorResponse{
		Status:  status,
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}

// HTTPBadRequest writes bad request response to response writer
func HTTPBadRequest(rw http.ResponseWriter) error {
	err := HTTPError(rw, http.StatusBadRequest, "Сервер получил некорректный запрос")
	if err != nil {
		return err
	}
	return nil
}

// HTTPInternalServerError writes internal server error response to response writer
func HTTPInternalServerError(rw http.ResponseWriter) error {
	err := HTTPError(rw, http.StatusInternalServerError,
		"Внутренняя ошибка сервера. Пожалуйста, попробуйте повторить позднее")
	if err != nil {
		return err
	}
	return nil
}

// HTTPUnauthorized writes unauthorized response to response writer
func HTTPUnauthorized(rw http.ResponseWriter) error {
	err := HTTPError(rw, http.StatusUnauthorized,
		"Для продолжения работы необходимо авторизоваться")
	if err != nil {
		return err
	}
	return nil
}

// HTTPOK writes ok response with optional data
func HTTPOK(rw http.ResponseWriter, data interface{}) error {
	result := map[string]interface{}{
		"status": http.StatusOK,
	}
	if data != nil {
		result["data"] = data
	}
	bytes, err := json.Marshal(&result)
	if err != nil {
		return xerrors.Errorf("cannot marshal json: %w", err)
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// InternalErrorResponse template with custom message
func InternalErrorResponse(message string) JSONErrorResponse {
	return JSONErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "Внутренняя ошибка сервера. Пожалуйста, попробуйте повторить позднее",
	}
}

// Unauthorized writes unauthorized response
func Unauthorized(w http.ResponseWriter) {
	bytes, err := json.Marshal(JSONErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: "Для продолжения работы необходимо авторизоваться",
	})
	if err != nil {
		logrus.Errorf("Cannot create json for unauthorized message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write(bytes)
}

// Error writes error response
func Error(w http.ResponseWriter, response JSONErrorResponse) {
	bytes, err := json.Marshal(response)
	if err != nil {
		logrus.Errorf("Cannot create json for unauthorized message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(response.Status)
	_, _ = w.Write(bytes)
}

// Ok writes ok response
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
		_, _ = w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

var (
	// ErrIncorrectBody error
	ErrIncorrectBody = xerrors.New("Incorrect body")
	// ErrIncorectFormFields error
	ErrIncorectFormFields = xerrors.New("Incorrect form fields")
)

// UnmarshalForm from body and validate it. If error occurs, writes error message to response writer
// TODO use JSONUnmarshalAndValidate instead
func UnmarshalForm(w http.ResponseWriter, r *http.Request, to interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Cannot read request body: %+v", err)
		Error(w, JSONErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Внутренняя ошибка сервера. Пожалуйста, повторите позднее",
		})
		return xerrors.Errorf("cannot read request body: %w", ErrIncorrectBody)
	}
	err = json.Unmarshal(body, to)
	if err != nil {
		logrus.Errorf("Cannot unmarshal request body: %+v", err)
		Error(w, JSONErrorResponse{
			Status:  http.StatusNotAcceptable,
			Message: "Сервер получил некорректный запрос",
		})
		return xerrors.Errorf("cannot unmarshal request body: %w", ErrIncorrectBody)
	}

	validate := validator.New()
	err = validate.Struct(to)
	if err != nil {
		logrus.Errorf("Invalid request body: %+v", err)
		Error(w, JSONErrorResponse{
			Status:  http.StatusNotAcceptable,
			Message: "Сервер получил некорректный запрос",
		})
		return xerrors.Errorf("invalid request body: %w", ErrIncorectFormFields)
	}
	return nil
}
