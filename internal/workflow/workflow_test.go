package workflow

import "testing"

func TestOrderedRespectsDependencies(t *testing.T) {
	w := &Workflow{
		ID:      "w",
		Version: "0.1.0",
		Steps: []Step{
			{ID: "b", Capability: "text.uppercase", Needs: []string{"a"}},
			{ID: "a", Capability: "text.constant", Needs: []string{}},
		},
	}
	ordered, err := w.Ordered()
	if err != nil {
		t.Fatalf("Ordered() error = %v", err)
	}
	if ordered[0].ID != "a" || ordered[1].ID != "b" {
		t.Fatalf("order = %s,%s; want a,b", ordered[0].ID, ordered[1].ID)
	}
}

func TestOrderedDetectsCycle(t *testing.T) {
	w := &Workflow{
		ID:      "w",
		Version: "0.1.0",
		Steps: []Step{
			{ID: "a", Capability: "text.x", Needs: []string{"b"}},
			{ID: "b", Capability: "text.x", Needs: []string{"a"}},
		},
	}
	if _, err := w.Ordered(); err == nil {
		t.Fatal("expected a cycle error")
	}
}

func TestValidateRejectsUnknownNeed(t *testing.T) {
	w := &Workflow{
		ID:      "w",
		Version: "0.1.0",
		Steps: []Step{
			{ID: "a", Capability: "text.x", Needs: []string{"ghost"}},
		},
	}
	if err := w.validate(); err == nil {
		t.Fatal("expected an error for an unknown need")
	}
}
