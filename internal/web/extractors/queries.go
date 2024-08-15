package extractors

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

var queryDecoder = schema.NewDecoder()

// init is required for queryDecoder and there are no side effects
func init() { //nolint:gochecknoinits
	queryDecoder.RegisterConverter(uuid.UUID{}, func(s string) reflect.Value {
		v, err := uuid.Parse(s)
		if err != nil {
			// TODO: what to do here?
			v = uuid.Nil
		}
		return reflect.ValueOf(v)
	})
}

func Query[T any](r *http.Request) (T, error) {
	var v T
	err := queryDecoder.Decode(&v, r.URL.Query())
	if err != nil {
		return v, err
	}

	// return validate.Check(v)
	return v, nil
}
