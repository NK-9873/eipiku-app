package handler

import (
	"fmt"
	"net/http"
)

// Health は GET /health
func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "learning-service ok")
}
