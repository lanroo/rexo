// Package kernel executes a deterministic REXO workflow end to end: it resolves
// step inputs, runs providers, stores content-addressed artifacts with
// provenance, writes an execution trace and an append-only event log, and can
// replay a run to verify determinism. It matches the contracts under
// contracts/ (task-envelope, artifact-manifest, execution-trace).
package kernel

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/lanroo/rexo/internal/workflow"
)

// Options configures a run. Now and RunID are injectable so tests are
// deterministic.
type Options struct {
	ProjectDir string
	RunID      string
	Now        func() time.Time
}

// ProducedBy records which task produced an artifact.
type ProducedBy struct {
	RunID      string `json:"run_id"`
	TaskID     string `json:"task_id"`
	Capability string `json:"capability"`
}

// ArtifactManifest matches contracts/artifact-manifest.schema.json.
type ArtifactManifest struct {
	SchemaVersion     string     `json:"schema_version"`
	ID                string     `json:"id"`
	ContentHash       string     `json:"content_hash"`
	MediaType         string     `json:"media_type"`
	SizeBytes         int        `json:"size_bytes"`
	ProducedBy        ProducedBy `json:"produced_by"`
	InputsFingerprint string     `json:"inputs_fingerprint"`
	CreatedAt         string     `json:"created_at"`
}

// WorkflowRef identifies the workflow a run executed.
type WorkflowRef struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// TaskResult is one entry in an execution trace.
type TaskResult struct {
	TaskID            string   `json:"task_id"`
	Capability        string   `json:"capability"`
	Status            string   `json:"status"`
	Attempts          int      `json:"attempts"`
	Needs             []string `json:"needs,omitempty"`
	OutputArtifact    *string  `json:"output_artifact"`
	InputsFingerprint string   `json:"inputs_fingerprint,omitempty"`
	StartedAt         string   `json:"started_at,omitempty"`
	EndedAt           string   `json:"ended_at,omitempty"`
	Error             string   `json:"error,omitempty"`
}

// Trace matches contracts/execution-trace.schema.json.
type Trace struct {
	SchemaVersion string       `json:"schema_version"`
	RunID         string       `json:"run_id"`
	Workflow      WorkflowRef  `json:"workflow"`
	Status        string       `json:"status"`
	StartedAt     string       `json:"started_at"`
	EndedAt       string       `json:"ended_at,omitempty"`
	Tasks         []TaskResult `json:"tasks"`
	Outputs       []string     `json:"outputs,omitempty"`
	EventLog      string       `json:"event_log,omitempty"`
}

// Run executes the workflow and returns its trace. A provider error stops the
// run, marks the trace failed, and leaves a diagnosable on-disk state; the
// returned error is non-nil in that case.
func Run(wf *workflow.Workflow, opt Options) (*Trace, error) {
	now := opt.Now
	if now == nil {
		now = time.Now
	}
	runID := opt.RunID
	if runID == "" {
		id, err := genRunID()
		if err != nil {
			return nil, err
		}
		runID = id
	}

	ordered, err := wf.Ordered()
	if err != nil {
		return nil, err
	}

	rexoDir := filepath.Join(opt.ProjectDir, ".rexo")
	artifactsDir := filepath.Join(rexoDir, "artifacts")
	runDir := filepath.Join(rexoDir, "runs", runID)
	runArtifactsDir := filepath.Join(runDir, "artifacts")
	if err := os.MkdirAll(artifactsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create artifacts dir: %w", err)
	}
	if err := os.MkdirAll(runArtifactsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create run dir: %w", err)
	}

	events, err := os.Create(filepath.Join(runDir, "events.jsonl"))
	if err != nil {
		return nil, fmt.Errorf("create event log: %w", err)
	}
	defer events.Close()

	trace := &Trace{
		SchemaVersion: "0.1.0",
		RunID:         runID,
		Workflow:      WorkflowRef{ID: wf.ID, Version: wf.Version},
		Status:        "running",
		StartedAt:     now().UTC().Format(time.RFC3339),
		EventLog:      "events.jsonl",
	}
	if err := writeJSON(filepath.Join(runDir, "trace.json"), trace); err != nil {
		return nil, err
	}
	appendEvent(events, now, "run_started", map[string]any{"run_id": runID, "workflow": wf.ID})

	outputs := make(map[string]string, len(ordered))
	var runErr error

	for _, step := range ordered {
		tr := TaskResult{
			TaskID:     step.ID,
			Capability: step.Capability,
			Needs:      step.Needs,
			StartedAt:  now().UTC().Format(time.RFC3339),
		}
		appendEvent(events, now, "task_started", map[string]any{"task": step.ID, "capability": step.Capability})

		provider, ok := providers[step.Capability]
		if !ok {
			runErr = failTask(&tr, now, fmt.Errorf("unknown capability %q", step.Capability))
			trace.Tasks = append(trace.Tasks, tr)
			appendEvent(events, now, "task_failed", map[string]any{"task": step.ID, "error": tr.Error})
			break
		}

		inputs, err := resolveInputs(step, outputs)
		if err != nil {
			runErr = failTask(&tr, now, err)
			trace.Tasks = append(trace.Tasks, tr)
			appendEvent(events, now, "task_failed", map[string]any{"task": step.ID, "error": tr.Error})
			break
		}

		tr.Attempts = 1
		output, err := provider(inputs)
		if err != nil {
			runErr = failTask(&tr, now, err)
			trace.Tasks = append(trace.Tasks, tr)
			appendEvent(events, now, "task_failed", map[string]any{"task": step.ID, "error": tr.Error})
			break
		}

		content := []byte(output)
		hash := sha256Hex(content)
		if err := writeArtifact(artifactsDir, hash, content); err != nil {
			runErr = failTask(&tr, now, err)
			trace.Tasks = append(trace.Tasks, tr)
			break
		}
		fingerprint := "sha256:" + inputsFingerprint(step.Capability, inputs)
		manifest := ArtifactManifest{
			SchemaVersion:     "0.1.0",
			ID:                step.ID,
			ContentHash:       "sha256:" + hash,
			MediaType:         "text/plain; charset=utf-8",
			SizeBytes:         len(content),
			ProducedBy:        ProducedBy{RunID: runID, TaskID: step.ID, Capability: step.Capability},
			InputsFingerprint: fingerprint,
			CreatedAt:         now().UTC().Format(time.RFC3339),
		}
		if err := writeJSON(filepath.Join(runArtifactsDir, step.ID+".json"), manifest); err != nil {
			runErr = failTask(&tr, now, err)
			trace.Tasks = append(trace.Tasks, tr)
			break
		}

		outputs[step.ID] = output
		id := step.ID
		tr.OutputArtifact = &id
		tr.InputsFingerprint = fingerprint
		tr.Status = "succeeded"
		tr.EndedAt = now().UTC().Format(time.RFC3339)
		trace.Tasks = append(trace.Tasks, tr)
		appendEvent(events, now, "task_succeeded", map[string]any{"task": step.ID, "content_hash": manifest.ContentHash})
	}

	if runErr == nil && len(wf.Outputs) > 0 {
		written, err := writeOutputs(opt.ProjectDir, wf.Outputs, outputs)
		if err != nil {
			runErr = err
		} else {
			trace.Outputs = written
			for _, p := range written {
				appendEvent(events, now, "output_written", map[string]any{"path": p})
			}
		}
	}

	if runErr != nil {
		trace.Status = "failed"
	} else {
		trace.Status = "succeeded"
	}
	trace.EndedAt = now().UTC().Format(time.RFC3339)
	appendEvent(events, now, "run_"+trace.Status, map[string]any{"run_id": runID})
	if err := writeJSON(filepath.Join(runDir, "trace.json"), trace); err != nil {
		return trace, err
	}
	return trace, runErr
}

// ReplayResult reports how a replay compared to a recorded run.
type ReplayResult struct {
	RunID   string
	Checked int
}

// Replay re-executes the workflow and asserts every task reproduces the
// content hash recorded for the given run. A mismatch is a hard error.
func Replay(wf *workflow.Workflow, opt Options, runID string) (*ReplayResult, error) {
	ordered, err := wf.Ordered()
	if err != nil {
		return nil, err
	}
	runArtifactsDir := filepath.Join(opt.ProjectDir, ".rexo", "runs", runID, "artifacts")
	if _, err := os.Stat(runArtifactsDir); err != nil {
		return nil, fmt.Errorf("no recorded run %q: %w", runID, err)
	}

	outputs := make(map[string]string, len(ordered))
	res := &ReplayResult{RunID: runID}
	for _, step := range ordered {
		provider, ok := providers[step.Capability]
		if !ok {
			return nil, fmt.Errorf("task %q: unknown capability %q", step.ID, step.Capability)
		}
		inputs, err := resolveInputs(step, outputs)
		if err != nil {
			return nil, fmt.Errorf("task %q: %w", step.ID, err)
		}
		output, err := provider(inputs)
		if err != nil {
			return nil, fmt.Errorf("task %q: %w", step.ID, err)
		}
		outputs[step.ID] = output

		var manifest ArtifactManifest
		data, err := os.ReadFile(filepath.Join(runArtifactsDir, step.ID+".json"))
		if err != nil {
			return nil, fmt.Errorf("task %q: read recorded manifest: %w", step.ID, err)
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("task %q: parse recorded manifest: %w", step.ID, err)
		}
		got := "sha256:" + sha256Hex([]byte(output))
		res.Checked++
		if got != manifest.ContentHash {
			return res, fmt.Errorf("replay mismatch on task %q: recorded %s, got %s (non-determinism detected)", step.ID, manifest.ContentHash, got)
		}
	}
	return res, nil
}

func failTask(tr *TaskResult, now func() time.Time, err error) error {
	tr.Status = "failed"
	tr.Error = err.Error()
	tr.EndedAt = now().UTC().Format(time.RFC3339)
	tr.OutputArtifact = nil
	return fmt.Errorf("task %q: %w", tr.TaskID, err)
}

func resolveInputs(step workflow.Step, outputs map[string]string) (map[string]any, error) {
	resolved := make(map[string]any, len(step.With))
	for key, value := range step.With {
		rv, err := resolveValue(value, outputs)
		if err != nil {
			return nil, fmt.Errorf("input %q: %w", key, err)
		}
		resolved[key] = rv
	}
	return resolved, nil
}

// resolveValue replaces every { "from_task": "<id>" } reference anywhere in the
// input tree with that task's output, so references work inside nested objects
// (e.g. a template's vars) and arrays.
func resolveValue(value any, outputs map[string]string) (any, error) {
	switch v := value.(type) {
	case map[string]any:
		if from, ok := v["from_task"].(string); ok && len(v) == 1 {
			out, ok := outputs[from]
			if !ok {
				return nil, fmt.Errorf("references task %q which produced no output", from)
			}
			return out, nil
		}
		m := make(map[string]any, len(v))
		for k, val := range v {
			rv, err := resolveValue(val, outputs)
			if err != nil {
				return nil, err
			}
			m[k] = rv
		}
		return m, nil
	case []any:
		arr := make([]any, len(v))
		for i, val := range v {
			rv, err := resolveValue(val, outputs)
			if err != nil {
				return nil, err
			}
			arr[i] = rv
		}
		return arr, nil
	default:
		return value, nil
	}
}

// writeOutputs writes each declared output file (path -> task id) into the
// project directory after a successful run, returning the paths in a stable
// order.
func writeOutputs(projectDir string, spec, outputs map[string]string) ([]string, error) {
	paths := make([]string, 0, len(spec))
	for path := range spec {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	written := make([]string, 0, len(paths))
	for _, path := range paths {
		taskID := spec[path]
		content, ok := outputs[taskID]
		if !ok {
			return written, fmt.Errorf("output %q references task %q which produced no output", path, taskID)
		}
		full := filepath.Join(projectDir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return written, fmt.Errorf("create output dir for %q: %w", path, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return written, fmt.Errorf("write output %q: %w", path, err)
		}
		written = append(written, path)
	}
	return written, nil
}

func inputsFingerprint(capability string, inputs map[string]any) string {
	payload := struct {
		Capability string         `json:"capability"`
		Inputs     map[string]any `json:"inputs"`
	}{capability, inputs}
	data, _ := json.Marshal(payload) // map keys are sorted, so this is stable
	return sha256Hex(data)
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func writeArtifact(dir, hash string, content []byte) error {
	path := filepath.Join(dir, hash)
	if _, err := os.Stat(path); err == nil {
		return nil // content-addressed and immutable; already stored
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write artifact: %w", err)
	}
	return nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", filepath.Base(path), err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filepath.Base(path), err)
	}
	return nil
}

func appendEvent(w io.Writer, now func() time.Time, event string, fields map[string]any) {
	record := map[string]any{
		"ts":    now().UTC().Format(time.RFC3339Nano),
		"event": event,
	}
	for k, v := range fields {
		record[k] = v
	}
	data, err := json.Marshal(record)
	if err != nil {
		return
	}
	_, _ = w.Write(append(data, '\n'))
}

func genRunID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate run id: %w", err)
	}
	return "run-" + hex.EncodeToString(b), nil
}
