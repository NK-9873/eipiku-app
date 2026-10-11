package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "user-service ok")
	})
	log.Println("user-service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
