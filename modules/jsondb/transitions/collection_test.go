package transitions

import "testing"

func TestCollectionGetAllLenClear(t *testing.T) {
	withTempWorkingDir(t)
	dbPath := "collection-db"
	doc := &documentTransitions{}
	col := &collectionTransitions{}

	if r := doc.Set(dbPath, "items", "a", map[string]interface{}{"v": 1}, nil); !r.Success {
		t.Fatalf("Set a failed: %v", r.Error)
	}
	if r := doc.Set(dbPath, "items", "b", map[string]interface{}{"v": 2}, nil); !r.Success {
		t.Fatalf("Set b failed: %v", r.Error)
	}

	// filejsondb.Set always writes an id ".idx" blob alongside the document
	// itself (see AddIndex), so 2 documents show up as 4 stored records.
	lenResult := col.Len(dbPath, "items")
	if !lenResult.Success {
		t.Fatalf("Len failed: %v", lenResult.Error)
	}
	count, _ := lenResult.Response.(map[string]interface{})["count"].(uint64)
	if count != 4 {
		t.Fatalf("expected count=4 (2 documents + 2 id-index blobs), got %v", lenResult.Response)
	}

	getAllResult := col.GetAll(dbPath, "items")
	if !getAllResult.Success {
		t.Fatalf("GetAll failed: %v", getAllResult.Error)
	}
	items, ok := getAllResult.Response.(map[string]interface{})
	if !ok || len(items) != 4 {
		t.Fatalf("unexpected GetAll response: %#v", getAllResult.Response)
	}

	clearResult := col.Clear(dbPath, "items")
	if !clearResult.Success {
		t.Fatalf("Clear failed: %v", clearResult.Error)
	}
}
