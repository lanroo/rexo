package project

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesPortableProject(t *testing.T) {
	target := filepath.Join(t.TempDir(), "Course Builder")
	result, err := Init(target)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result.Name != "Course Builder" {
		t.Fatalf("result.Name = %q", result.Name)
	}

	required := []string{
		"rexo.project.json",
		"AGENTS.md",
		"REXO_BOOTSTRAP.md",
		"REXO_STATE.md",
		".rexo/artifacts/README.md",
		".rexo/memory/README.md",
	}
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(name))); err != nil {
			t.Errorf("required file %s: %v", name, err)
		}
	}

	data, err := os.ReadFile(filepath.Join(target, "rexo.project.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}
}

func TestInitDoesNotOverwrite(t *testing.T) {
	target := filepath.Join(t.TempDir(), "existing")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := Init(target)
	if !errors.Is(err, ErrTargetExists) {
		t.Fatalf("Init() error = %v, want ErrTargetExists", err)
	}
}

func TestSlug(t *testing.T) {
	if got := slug(" Curso de Java! "); got != "curso-de-java" {
		t.Fatalf("slug() = %q", got)
	}
}
