// Package main is the entry point for the sample Go application.
package main

import (
	"fmt"
	"net/http"

	"example.com/sample/handler"
)

// Port is the default HTTP server port.
const Port = "8080"

// main starts the HTTP server and registers routes.
func main() {
	mux := setupRouter()
	fmt.Printf("Server starting on port %s\n", Port)
	http.ListenAndServe(":"+Port, mux)
}

// setupRouter creates and configures the HTTP router.
// It registers all application routes and returns the configured ServeMux.
func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.GetUser)
	return mux
}
