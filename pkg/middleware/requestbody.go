package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"io"
	"net/http"
)

var ValidatorInstance = validator.New()

type RequestBodyMiddleware[T any] struct {
	Logger utils.ILogger
}

func (dep *RequestBodyMiddleware[T]) RequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			dep.Logger.Error("request body is nil")
			utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid request body"})
			return
		}

		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			dep.Logger.Error(err.Error())
			utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid request body"})
			return
		}

		if err := ValidatorInstance.Struct(payload); err != nil {
			dep.Logger.Error(err)
			utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid request body"})
			return
		}

		by, err := json.Marshal(payload)
		if err != nil {
			dep.Logger.Error(err.Error())
			utils.ErrorResponse(w, errors.New("internal server error"))
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(by))

		next.ServeHTTP(w, r)
	})
}
