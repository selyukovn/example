package kernel

import (
	"encoding/json"
	"net/http"
)

func ParseRequestJson[T any](r *http.Request, dataExample T) *T {
	p := &dataExample
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return nil
	}
	return p
}
