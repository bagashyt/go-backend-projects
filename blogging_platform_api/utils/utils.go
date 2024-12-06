package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func WriteError(w *gin.Context, status int, v any) error {
	w.Writer.Header().Add("Content-Type", "application/json")
	w.Writer.WriteHeader(status)

	return json.NewEncoder(w.Writer).Encode(v)

}

func WriteJSON(w *gin.Context, status int, v any) error {
	w.Writer.Header().Add("Content-Type", "application/json")
	w.Writer.WriteHeader(status)

	return json.NewEncoder(w.Writer).Encode(v)
}

func ParseJSON(r *gin.Context, payload any) error {
	if r.Request.Response.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Request.Body).Decode(payload)
}
