// agent-cluster-backend: decision 003 first vertical slice.
//
// Starts the HTTP server in internal/slice on $PORT (default 8080).
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kimjooyoon/agent-cluster-backend/internal/slice"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("agent-cluster-backend: decision 003 first vertical slice\n")
	fmt.Printf("listening on http://localhost%s\n", addr)
	fmt.Printf("  POST /graphql   { workItems { id title state } }\n")
	fmt.Printf("  GET  /events    SSE WorkItemCreated\n")
	log.Fatal(http.ListenAndServe(addr, slice.Handler()))
}
