package transitions

import "testing"

func TestIndexSetGetExistsDelete(t *testing.T) {
	withTempWorkingDir(t)
	dbPath := "index-db"
	doc := &documentTransitions{}
	idx := &indexTransitions{}

	value := map[string]interface{}{"email": "anar@kuetix.com"}
	if r := doc.Set(dbPath, "users", "u1", value, []string{"email"}); !r.Success {
		t.Fatalf("Set failed: %v", r.Error)
	}

	existsResult := idx.IsExistsIndexes(dbPath, "users", "u1")
	if !existsResult.Success {
		t.Fatalf("IsExistsIndexes operation failed: %v", existsResult.Error)
	}
	if exists, _ := existsResult.Response.(map[string]interface{})["exists"].(bool); !exists {
		t.Fatalf("expected indexes to exist after Set")
	}

	getResult := idx.GetIndexes(dbPath, "users", "u1")
	if !getResult.Success {
		t.Fatalf("GetIndexes failed: %v", getResult.Error)
	}

	setIndexResult := idx.SetIndex(dbPath, "users", "u1", value, []string{"email"})
	if !setIndexResult.Success {
		t.Fatalf("SetIndex failed: %v", setIndexResult.Error)
	}

	deleteResult := idx.DeleteIndexes(dbPath, "users", "u1")
	if !deleteResult.Success {
		t.Fatalf("DeleteIndexes failed: %v", deleteResult.Error)
	}
}
