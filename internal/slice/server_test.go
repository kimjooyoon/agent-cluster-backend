package slice

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kimjooyoon/agent-cluster-backend/internal/contracts"
)

func TestGraphQLWorkItemsReturnsGeneratedType(t *testing.T) {
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	body := `{"query":"{ workItems { id title state } }"}`
	resp, err := http.Post(srv.URL+"/graphql", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got struct {
		Data struct {
			WorkItems []contracts.WorkItem `json:"workItems"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if len(got.Data.WorkItems) != len(SeedWorkItems) {
		t.Fatalf("got %d items, want %d", len(got.Data.WorkItems), len(SeedWorkItems))
	}
	if got.Data.WorkItems[0].Id != "WI-001" {
		t.Errorf("first item id = %q, want WI-001", got.Data.WorkItems[0].Id)
	}
}

func TestGraphQLRejectsOtherQueries(t *testing.T) {
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	body := `{"query":"{ otherField }"}`
	resp, err := http.Post(srv.URL+"/graphql", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestSSEEmitsWorkItemCreatedWithGeneratedEventName(t *testing.T) {
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("content-type = %q, want text/event-stream", ct)
	}
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	got := string(buf[:n])
	wantEvent := "event: " + contracts.WorkItemCreatedEventName
	if !strings.Contains(got, wantEvent) {
		t.Errorf("SSE body missing %q\n--- body ---\n%s", wantEvent, got)
	}
	if !strings.Contains(got, `"work_item_id":"WI-002"`) {
		t.Errorf("SSE body missing work_item_id wire field\n--- body ---\n%s", got)
	}
}

func TestNormalizeQuery(t *testing.T) {
	cases := map[string]string{
		"{ workItems { id title state } }":            "{ workItems { id title state } }",
		"{\n  workItems {\n    id title state\n  }\n}": "{ workItems { id title state } }",
		"  {  workItems  {  id  title  state  }  }  ": "{ workItems { id title state } }",
	}
	for in, want := range cases {
		if got := normalizeQuery(in); got != want {
			t.Errorf("normalizeQuery(%q) = %q, want %q", in, got, want)
		}
	}
}
