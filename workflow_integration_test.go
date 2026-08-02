package jsondb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kuetix/std-jsondb/modules"

	"github.com/kuetix/engine"
	"github.com/kuetix/engine/engine/domain"
	engineModule "github.com/kuetix/engine/modules"
)

// runWorkflow runs a jsondb workflow through the real engine (parsing,
// dependency injection, and transition dispatch), exercising the exact path a
// consuming app would use, with args passed as native Go values via Context
// (as an HTTP handler or another workflow would), rather than as CLI strings.
func runWorkflow(t *testing.T, name string, args map[string]interface{}) *workflowResult {
	t.Helper()

	engineModule.Enable()
	modules.Enable()

	responses := engine.RunWorkflow("production", &domain.Options{
		EngineName: "jsondb-test",
		ConfigName: "engine",
		Workflow:   "@jsondb/" + name,
		Amount:     1,
		Retry:      1,
		LogPath:    "stdout",
		Context: map[string]interface{}{
			"args": args,
		},
	})

	res, ok := responses["main"]
	if !ok {
		t.Fatalf("workflow %q: no response returned", name)
	}
	return &workflowResult{StatusCode: res.StatusCode, Err: res.GetError(), Response: res.Response}
}

type workflowResult struct {
	StatusCode int
	Err        error
	Response   interface{}
}

func TestWorkflowsEndToEnd(t *testing.T) {
	dir := t.TempDir()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir into temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	if err := os.Symlink(filepath.Join(original, "workflows"), filepath.Join(dir, "workflows")); err != nil {
		t.Fatalf("failed to symlink workflows dir: %v", err)
	}

	dbPath := "wf-test-db"

	r := runWorkflow(t, "ping", map[string]interface{}{"path": dbPath})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("ping failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "set", map[string]interface{}{
		"path":       dbPath,
		"collection": "users",
		"id":         "u1",
		"value":      map[string]interface{}{"name": "Anar", "role": "maintainer"},
		"indexes":    []string{"role"},
	})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("set failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "get", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("get failed: err=%v status=%d", r.Err, r.StatusCode)
	}
	value, ok := r.Response.(map[string]interface{})
	if !ok || value["name"] != "Anar" {
		t.Fatalf("unexpected get response: %#v", r.Response)
	}

	r = runWorkflow(t, "exists", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("exists failed: err=%v status=%d", r.Err, r.StatusCode)
	}
	if exists, _ := r.Response.(map[string]interface{})["exists"].(bool); !exists {
		t.Fatalf("expected exists=true, got %#v", r.Response)
	}

	r = runWorkflow(t, "update", map[string]interface{}{
		"path": dbPath, "collection": "users", "id": "u1",
		"value": map[string]interface{}{"name": "Anar", "role": "lead"},
	})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("update failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "list", map[string]interface{}{"path": dbPath, "collection": "users"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("list failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "count", map[string]interface{}{"path": dbPath, "collection": "users"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("count failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "set_index", map[string]interface{}{
		"path": dbPath, "collection": "users", "id": "u1",
		"value": map[string]interface{}{"name": "Anar", "role": "lead"}, "fields": []string{"role"},
	})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("set_index failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "get_indexes", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("get_indexes failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "index_exists", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("index_exists failed: err=%v status=%d", r.Err, r.StatusCode)
	}
	if exists, _ := r.Response.(map[string]interface{})["exists"].(bool); !exists {
		t.Fatalf("expected index exists=true, got %#v", r.Response)
	}

	r = runWorkflow(t, "delete_index", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("delete_index failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "clear", map[string]interface{}{"path": dbPath, "collection": "users"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("clear failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "delete", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("delete failed: err=%v status=%d", r.Err, r.StatusCode)
	}

	r = runWorkflow(t, "exists", map[string]interface{}{"path": dbPath, "collection": "users", "id": "u1"})
	if r.Err != nil || r.StatusCode != 200 {
		t.Fatalf("final exists failed: err=%v status=%d", r.Err, r.StatusCode)
	}
	if exists, _ := r.Response.(map[string]interface{})["exists"].(bool); exists {
		t.Fatalf("expected exists=false after delete, got %#v", r.Response)
	}
}
