package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"golang.org/x/xerrors"
	validator "gopkg.in/go-playground/validator.v9"
)

// JSONUnmarshalAndValidate unmarshals json from source to destination
func JSONUnmarshalAndValidate(source io.Reader, target interface{}) error {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return xerrors.Errorf("cannot read data: %w", ErrIncorrectBody)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return xerrors.Errorf("cannot unmarshal data: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(target)
	if err != nil {
		return xerrors.Errorf("invalid request body: %w", err)
	}

	return nil
}
