// Package extractors contains extractor functions for getting data out
// of the request.
package extractors

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Body decodes the request body into a struct and validates the struct using
// the `validate` tags
func Body[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		if errors.Is(err, io.EOF) {
			return v, errors.New("empty body")
		}

		return v, err
	}

	/* err := validate.Check(v) */
	/* if err != nil { */
	/* 	return v, err */
	/* } */
	/**/
	/* return v, err */
	return v, nil
}

// ID extracts the id field from the path, validates it, and returns it as
// a ID. If the id is not a valid ID, an error is returned.
func ID(r *http.Request, key string) (uuid.UUID, error) {
	v := chi.URLParam(r, key)

	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, err
		// return uuid.Nil, validate.NewRouteKeyErrorWithMessage(key, "unable to parse UUID")
	}

	return id, nil
}

// BodyWithID combines the calls of Body and ID into one call and extracts both
// the ID and the body of the request.
func BodyWithID[T any](r *http.Request, key string) (uuid.UUID, T, error) {
	id, err := ID(r, key)
	if err != nil {
		var v T
		return uuid.Nil, v, err
	}

	body, err := Body[T](r)
	if err != nil {
		var v T
		return uuid.Nil, v, err
	}

	return id, body, nil
}


