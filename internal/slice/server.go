// Package slice holds the HTTP surface for decision 003's first vertical
// slice. It exposes:
//
//   POST /graphql    — accepts exactly the query
//                       { workItems { id title state } }
//                      and responds with a hardcoded list constructed from the
//                      generated contracts.WorkItem type.
//   GET  /events     — SSE; emits one WorkItemCreated event on connect, then
//                      closes the stream.
//
// No GraphQL library is used. The minimal surface is intentional — see
// decision 003. A real GraphQL library lands behind its own decision.
package slice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kimjooyoon/agent-cluster-backend/internal/contracts"
)

// SeedWorkItems is the single fixture the slice serves. Backend never
// invents WorkItem fields — they come from the generated contracts package.
var SeedWorkItems = []contracts.WorkItem{
	{Id: "WI-001", Title: "Bootstrap contracts repo", State: "done"},
	{Id: "WI-002", Title: "Wire first vertical slice", State: "in_progress"},
}

// SeedCreatedEvent is the single event emitted by /events for this slice.
var SeedCreatedEvent = contracts.WorkItemCreated{
	WorkItemId: "WI-002",
	Title:      "Wire first vertical slice",
}

// Handler returns the HTTP mux for the slice.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", handleGraphQL)
	mux.HandleFunc("/events", handleSSE)
	mux.HandleFunc("/", handleRoot)
	return mux
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "agent-cluster-backend: decision 003 first vertical slice")
	fmt.Fprintln(w, "POST /graphql  -> { workItems { id title state } }")
	fmt.Fprintln(w, "GET  /events   -> SSE WorkItemCreated")
}

// graphQLRequest is the minimum subset of the GraphQL HTTP shape we accept.
type graphQLRequest struct {
	Query string `json:"query"`
}

// graphQLResponse is similarly minimal.
type graphQLResponse struct {
	Data   any              `json:"data,omitempty"`
	Errors []graphQLMessage `json:"errors,omitempty"`
}

type graphQLMessage struct {
	Message string `json:"message"`
}

type workItemsData struct {
	WorkItems []contracts.WorkItem `json:"workItems"`
}

// supportedQuery is the only query string this slice answers. The hand-written
// parser only matches by normalized whitespace; a real GraphQL impl arrives
// behind its own decision.
const supportedQuery = "{ workItems { id title state } }"

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeGraphQLError(w, "POST required")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<16))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeGraphQLError(w, "read body: "+err.Error())
		return
	}
	var req graphQLRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeGraphQLError(w, "parse body: "+err.Error())
		return
	}
	if normalizeQuery(req.Query) != supportedQuery {
		w.WriteHeader(http.StatusBadRequest)
		writeGraphQLError(w, "this slice only supports query: "+supportedQuery)
		return
	}
	_ = json.NewEncoder(w).Encode(graphQLResponse{
		Data: workItemsData{WorkItems: SeedWorkItems},
	})
}

func writeGraphQLError(w http.ResponseWriter, msg string) {
	_ = json.NewEncoder(w).Encode(graphQLResponse{
		Errors: []graphQLMessage{{Message: msg}},
	})
}

// normalizeQuery collapses runs of whitespace to single spaces and trims, so
// "{\n  workItems {\n    id title state\n  }\n}" matches supportedQuery.
func normalizeQuery(q string) string {
	var b strings.Builder
	prevSpace := true
	for _, r := range q {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		b.WriteRune(r)
		prevSpace = false
	}
	return strings.TrimSpace(b.String())
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported by this server", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	payload, err := json.Marshal(SeedCreatedEvent)
	if err != nil {
		http.Error(w, "marshal event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Use the generated event-name constant so backend never invents the
	// wire identifier. If the IR changes the event name, this string changes.
	fmt.Fprintf(w, "event: %s\n", contracts.WorkItemCreatedEventName)
	fmt.Fprintf(w, "data: %s\n\n", payload)
	flusher.Flush()
}
