package openapi

import (
	"example/admin/bff/cmd/http/container"
	"net/http"
)

func NewStrictResponseErrorHandler(ctr *container.Container) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		switch err.(type) {
		case nil:
		case *InvalidParamFormatError:
			// Сгенерированный код использует эту ошибку при некорректных хедерах, вроде X-Trace-Id.
			// Поскольку эти хедеры генерируются в `rproxy`, для клиента это означает ошибку сервера --
			// поэтому нужно отвечать с кодом 500, а не 400.
			ctr.Logger.CtxErrorFf(r.Context(), "NewStrictResponseErrorHandler: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			ctr.Logger.CtxErrorFf(r.Context(), "NewStrictResponseErrorHandler: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
