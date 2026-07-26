
package workflow

import (
	"encoding/json"
	"fmt"
	"os"
)

// Step is one node of the workflow DAG.
type Step struct {
	ID          string         `json:"id"`
	Capability  string         `json:"capability"`
	Needs       []string       `json:"needs"`
	With        map[string]any `json:"with,omitempty"`
	MaxAttempts int            `json:"max_attempts,omitempty"`
}

// Workflow is a deterministic DAG of steps.
type Workflow struct {
	SchemaVersion string `json:"schema_version"`
	ID            string `json:"id"`
	Version       string `json:"version"`
	Steps         []Step `json:"steps"`
}

// Load reads and parses a workflow definition from disk.
func Load(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workflow: %w", err)
	}
	var wf Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("parse workflow: %w", err)
	}
	if err := wf.validate(); err != nil {
		return nil, err
	}
	return &wf, nil
}

func (w *Workflow) validate() error {
	if w.ID == "" {
		return fmt.Errorf("workflow id is required")
	}
	if len(w.Steps) == 0 {
		return fmt.Errorf("workflow %q has no steps", w.ID)
	}
	seen := make(map[string]bool, len(w.Steps))
	for _, s := range w.Steps {
		if s.ID == "" {
			return fmt.Errorf("workflow %q has a step with no id", w.ID)
		}
		if seen[s.ID] {
			return fmt.Errorf("duplicate step id %q", s.ID)
		}
		seen[s.ID] = true
		if s.Capability == "" {
			return fmt.Errorf("step %q has no capability", s.ID)
		}
	}
	for _, s := range w.Steps {
		for _, need := range s.Needs {
			if !seen[need] {
				return fmt.Errorf("step %q needs unknown step %q", s.ID, need)
			}
			if need == s.ID {
				return fmt.Errorf("step %q cannot depend on itself", s.ID)
			}
		}
	}
	return nil
}

// Ordered returns the steps in a deterministic topological order. Ties are
// broken by the step's position in the file so runs are reproducible. It fails
// if the graph contains a cycle.
func (w *Workflow) Ordered() ([]Step, error) {
	index := make(map[string]int, len(w.Steps))
	remaining := make(map[string]int, len(w.Steps))
	for i, s := range w.Steps {
		index[s.ID] = i
		remaining[s.ID] = len(s.Needs)
	}

	ordered := make([]Step, 0, len(w.Steps))
	done := make(map[string]bool, len(w.Steps))
	for len(ordered) < len(w.Steps) {
		// Pick the earliest-declared step whose needs are all satisfied.
		next := -1
		for i, s := range w.Steps {
			if done[s.ID] || remaining[s.ID] != 0 {
				continue
			}
			if next == -1 || i < next {
				next = i
			}
		}
		if next == -1 {
			return nil, fmt.Errorf("workflow %q has a dependency cycle", w.ID)
		}
		s := w.Steps[next]
		done[s.ID] = true
		ordered = append(ordered, s)
		for _, other := range w.Steps {
			for _, need := range other.Needs {
				if need == s.ID {
					remaining[other.ID]--
				}
			}
		}
	}
	return ordered, nil
}
