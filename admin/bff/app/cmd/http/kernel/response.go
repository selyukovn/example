package kernel

import (
	"encoding/json"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------
// REDIRECT
// ---------------------------------------------------------------------------------------------------------------------

func Redirect302(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

func Redirect307(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// ---------------------------------------------------------------------------------------------------------------------
// ERRORS
// ---------------------------------------------------------------------------------------------------------------------

func Error400(w http.ResponseWriter, customMessage ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	message := http.StatusText(http.StatusBadRequest)
	if len(customMessage) > 0 && len(customMessage[0]) > 0 {
		message = customMessage[0]
	}
	http.Error(w, message, http.StatusBadRequest)
}

func Error401(w http.ResponseWriter, customMessage ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	message := http.StatusText(http.StatusUnauthorized)
	if len(customMessage) > 0 && len(customMessage[0]) > 0 {
		message = customMessage[0]
	}
	http.Error(w, message, http.StatusUnauthorized)
}

func Error403(w http.ResponseWriter, customMessage ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	message := http.StatusText(http.StatusForbidden)
	if len(customMessage) > 0 && len(customMessage[0]) > 0 {
		message = customMessage[0]
	}
	http.Error(w, message, http.StatusForbidden)
}

func Error404(w http.ResponseWriter, customMessage ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	message := http.StatusText(http.StatusNotFound)
	if len(customMessage) > 0 && len(customMessage[0]) > 0 {
		message = customMessage[0]
	}
	http.Error(w, message, http.StatusNotFound)
}

func Error422(w http.ResponseWriter, customMessage ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	message := http.StatusText(http.StatusUnprocessableEntity)
	if len(customMessage) > 0 && len(customMessage[0]) > 0 {
		message = customMessage[0]
	}
	http.Error(w, message, http.StatusUnprocessableEntity)
}

func Error424(w http.ResponseWriter, customMessage ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	message := http.StatusText(http.StatusFailedDependency)
	if len(customMessage) > 0 && len(customMessage[0]) > 0 {
		message = customMessage[0]
	}
	http.Error(w, message, http.StatusFailedDependency)
}

func Error500(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ---------------------------------------------------------------------------------------------------------------------
// RENDER
// ---------------------------------------------------------------------------------------------------------------------

func RenderJson(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
