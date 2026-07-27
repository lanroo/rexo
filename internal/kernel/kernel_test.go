package kernel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lanroo/rexo/internal/workflow"
)

func helloWorkflow() *workflow.Workflow {
	return &workflow.Workflow{
		SchemaVersion: "0.1.0",
		ID:            "hello",
		Version:       "0.1.0",
		Steps: []workflow.Step{
			{
				ID:         "greeting",
				Capability: "text.constant",
				Needs:      []string{},
				With:       map[string]any{"value": "hello world"},
			},
			{
				ID:         "shout",
				Capability: "text.uppercase",
				Needs:      []string{"greeting"},
				With:       map[string]any{"text": map[string]any{"from_task": "greeting"}},
			},
		},
	}
}

func fixedNow() func() time.Time {
	t := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func TestRunAndReplayHello(t *testing.T) {
	dir := t.TempDir()
	opt := Options{ProjectDir: dir, RunID: "run-testfixed01", Now: fixedNow()}

	trace, err := Run(helloWorkflow(), opt)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if trace.Status != "succeeded" {
		t.Fatalf("status = %q, want succeeded", trace.Status)
	}
	if len(trace.Tasks) != 2 {
		t.Fatalf("tasks = %d, want 2", len(trace.Tasks))
	}

	shout := trace.Tasks[1]
	if shout.OutputArtifact == nil || *shout.OutputArtifact != "shout" {
		t.Fatalf("shout output artifact = %v", shout.OutputArtifact)
	}

	// Resolve the shout output through its manifest and verify the bytes.
	manifestData, err := os.ReadFile(filepath.Join(dir, ".rexo", "runs", opt.RunID, "artifacts", "shout.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m ArtifactManifest
	if err := json.Unmarshal(manifestData, &m); err != nil {
		t.Fatal(err)
	}
	hexHash := strings.TrimPrefix(m.ContentHash, "sha256:")
	content, err := os.ReadFile(filepath.Join(dir, ".rexo", "artifacts", hexHash))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "HELLO WORLD" {
		t.Fatalf("shout content = %q, want %q", content, "HELLO WORLD")
	}

	res, err := Replay(helloWorkflow(), opt, opt.RunID)
	if err != nil {
		t.Fatalf("Replay() error = %v", err)
	}
	if res.Checked != 2 {
		t.Fatalf("replay checked = %d, want 2", res.Checked)
	}
}

func TestRunUnknownCapabilityIsDiagnosable(t *testing.T) {
	dir := t.TempDir()
	wf := &workflow.Workflow{
		SchemaVersion: "0.1.0",
		ID:            "bad",
		Version:       "0.1.0",
		Steps: []workflow.Step{
			{ID: "x", Capability: "text.nope", Needs: []string{}, With: map[string]any{}},
		},
	}
	opt := Options{ProjectDir: dir, RunID: "run-badfixed001", Now: fixedNow()}

	trace, err := Run(wf, opt)
	if err == nil {
		t.Fatal("expected an error for an unknown capability")
	}
	if trace == nil || trace.Status != "failed" {
		t.Fatalf("trace = %+v, want failed status", trace)
	}
	if trace.Tasks[0].Status != "failed" {
		t.Fatalf("task status = %q, want failed", trace.Tasks[0].Status)
	}
	// The partial trace must remain on disk for diagnosis.
	if _, err := os.Stat(filepath.Join(dir, ".rexo", "runs", opt.RunID, "trace.json")); err != nil {
		t.Fatalf("trace not persisted: %v", err)
	}
}

func TestRunWritesNamedOutputWithTemplate(t *testing.T) {
	dir := t.TempDir()
	wf := &workflow.Workflow{
		SchemaVersion: "0.1.0",
		ID:            "welcome",
		Version:       "0.1.0",
		Steps: []workflow.Step{
			{ID: "shout", Capability: "text.uppercase", Needs: []string{}, With: map[string]any{"text": "hi"}},
			{
				ID:         "doc",
				Capability: "text.template",
				Needs:      []string{"shout"},
				With: map[string]any{
					"template": "value: {{v}}",
					// reference nested inside vars — exercises recursive resolution
					"vars": map[string]any{"v": map[string]any{"from_task": "shout"}},
				},
			},
		},
		Outputs: map[string]string{"out/welcome.md": "doc"},
	}
	opt := Options{ProjectDir: dir, RunID: "run-outfixed001", Now: fixedNow()}

	trace, err := Run(wf, opt)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if trace.Status != "succeeded" {
		t.Fatalf("status = %q, want succeeded", trace.Status)
	}
	if len(trace.Outputs) != 1 || trace.Outputs[0] != "out/welcome.md" {
		t.Fatalf("outputs = %v, want [out/welcome.md]", trace.Outputs)
	}

	content, err := os.ReadFile(filepath.Join(dir, "out", "welcome.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "value: HI" {
		t.Fatalf("welcome.md = %q, want %q", content, "value: HI")
	}
}
