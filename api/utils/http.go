package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v9"
)

// Unauthorized writes unauthorized response
func Unauthorized(w http.ResponseWriter) {
	bytes, err := json.Marshal(&map[string]string{
		"status":  strconv.Itoa(http.StatusUnauthorized),
		"message": "Invalid authorization token",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write(bytes)
}

// Error writes error response
func Error(w http.ResponseWriter, code int, message string) {
	bytes, err := json.Marshal(&map[string]string{
		"status":  strconv.Itoa(code),
		"message": message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(bytes)
}

// Ok writes error response
func Ok(w http.ResponseWriter, data interface{}) {
	result := map[string]interface{}{
		"status": strconv.Itoa(http.StatusOK),
	}
	if data != nil {
		result["data"] = data
	}
	bytes, err := json.Marshal(&result)
	if err != nil {
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

// UnmarshalForm from body and validate it
func UnmarshalForm(r *http.Request, to interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrIncorrectBody
	}
	err = json.Unmarshal(body, to)
	if err != nil {
		return ErrIncorrectBody
	}

	validate := validator.New()
	err = validate.Struct(to)
	if err != nil {
		return ErrIncorectFormFields
	}
	return nil
}
