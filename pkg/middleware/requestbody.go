package middleware

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type RequestBodyMiddleware[T any] struct {
	Logger utils.ILogger
}

var validate = validator.New()

func (dep *RequestBodyMiddleware[T]) RequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			dep.Logger.Error(err)
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "invalid request body",
				},
			)
			return
		}

		if err := validate.Struct(payload); err != nil {
			dep.Logger.Error(err)
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: err.Error(),
				},
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
