package transitions

import (
	"os"
	"testing"
)

// filejsondb resolves its path relative to the process's working directory
// (its own test suite does the same), so tests chdir into a temp dir first.
func withTempWorkingDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir into temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(original)
	})
}

func TestDocumentSetGetUpdateDeleteExists(t *testing.T) {
	withTempWorkingDir(t)
	dbPath := "test-db"
	doc := &documentTransitions{}

	setResult := doc.Set(dbPath, "users", "u1", map[string]interface{}{"name": "Anar", "role": "maintainer"}, []string{"role"})
	if !setResult.Success {
		t.Fatalf("Set failed: %v (status=%d)", setResult.Error, setResult.StatusCode)
	}

	getResult := doc.Get(dbPath, "users", "u1")
	if !getResult.Success {
		t.Fatalf("Get failed: %v (status=%d)", getResult.Error, getResult.StatusCode)
	}
	value, ok := getResult.Response.(map[string]interface{})
	if !ok || value["name"] != "Anar" {
		t.Fatalf("unexpected Get response: %#v", getResult.Response)
	}

	existsResult := doc.Exists(dbPath, "users", "u1")
	if !existsResult.Success {
		t.Fatalf("Exists operation failed: %v", existsResult.Error)
	}
	if exists, _ := existsResult.Response.(map[string]interface{})["exists"].(bool); !exists {
		t.Fatalf("expected document to exist")
	}

	updateResult := doc.Update(dbPath, "users", "u1", map[string]interface{}{"name": "Anar", "role": "lead"})
	if !updateResult.Success {
		t.Fatalf("Update failed: %v (status=%d)", updateResult.Error, updateResult.StatusCode)
	}

	deleteResult := doc.Delete(dbPath, "users", "u1")
	if !deleteResult.Success {
		t.Fatalf("Delete failed: %v (status=%d)", deleteResult.Error, deleteResult.StatusCode)
	}

	afterDelete := doc.Exists(dbPath, "users", "u1")
	if exists, _ := afterDelete.Response.(map[string]interface{})["exists"].(bool); exists {
		t.Fatalf("expected document to no longer exist after delete")
	}
}
