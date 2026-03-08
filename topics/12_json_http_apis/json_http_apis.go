// json_http_apis.go
// Demonstrates building a simple REST API with net/http and JSON.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := Message{Text: "Hello, Go API!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("Server running at http://localhost:8080/hello")
	http.ListenAndServe(":8080", nil)
}
// Documentation:
// - Use net/http for REST APIs.
// - Use encoding/json for JSON encoding/decoding.
// - Edge cases: error handling, malformed JSON, status codes.
