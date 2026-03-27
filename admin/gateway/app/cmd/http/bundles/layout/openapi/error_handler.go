package openapi

import (
	"github.com/selyukovn/go-std/logger"
	"net/http"
)

func NewStrictResponseErrorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		switch err.(type) {
		case nil:
		default:
			logger.ErrorFf(r.Context(), "NewStrictResponseErrorHandler: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
