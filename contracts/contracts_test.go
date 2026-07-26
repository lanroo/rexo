package contracts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSchemasAreValidJSON(t *testing.T) {
	files, err := filepath.Glob("*.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no schema files found")
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("%s: %v", file, err)
			continue
		}
		var document map[string]any
		if err := json.Unmarshal(data, &document); err != nil {
			t.Errorf("%s is not valid JSON: %v", file, err)
		}
		if document["$schema"] == nil || document["$id"] == nil {
			t.Errorf("%s must declare $schema and $id", file)
		}
	}
}
