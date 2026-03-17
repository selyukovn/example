package openapi

import (
	"example/admin/gateway/cmd/http/container"
	"net/http"
)

func NewStrictResponseErrorHandler(ctr *container.Container) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		switch err.(type) {
		case nil:
		default:
			ctr.Logger.CtxErrorFf(r.Context(), "NewStrictResponseErrorHandler: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
