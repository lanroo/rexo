package providers

import (
	"context"
	"strings"
	"testing"
)

// fakeGen is a controllable Generator for tests: it counts calls so we can
// prove the cache prevents a second model call.
type fakeGen struct {
	id        string
	available bool
	calls     int
	text      string
}

func (f *fakeGen) ID() string   { return f.id }
func (f *fakeGen) Detect() bool { return f.available }
func (f *fakeGen) Generate(ctx context.Context, req Request) (*Result, error) {
	f.calls++
	return &Result{Text: f.text, ProviderID: f.id}, nil
}

func TestResolverCachesGeneration(t *testing.T) {
	dir := t.TempDir()
	fake := &fakeGen{id: "fake", available: true, text: "hello world"}
	r := NewResolver(dir, "", fake)

	req := Request{Prompt: "say hi"}

	out1, res1, err := r.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("first generate: %v", err)
	}
	if res1.CacheHit {
		t.Error("first call should be a cache miss")
	}
	if out1.Text != "hello world" {
		t.Errorf("got %q", out1.Text)
	}
	if fake.calls != 1 {
		t.Fatalf("expected 1 model call, got %d", fake.calls)
	}

	// Same request again: must be served from cache, zero new model calls.
	_, res2, err := r.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("second generate: %v", err)
	}
	if !res2.CacheHit {
		t.Error("second identical call should hit the cache")
	}
	if fake.calls != 1 {
		t.Errorf("cache miss: model was called %d times, expected 1", fake.calls)
	}
	if res1.CacheKey != res2.CacheKey {
		t.Errorf("cache keys differ for identical requests: %s vs %s", res1.CacheKey, res2.CacheKey)
	}
}

func TestResolverPrefersRequestedProvider(t *testing.T) {
	dir := t.TempDir()
	a := &fakeGen{id: "a", available: true, text: "from a"}
	b := &fakeGen{id: "b", available: true, text: "from b"}
	r := NewResolver(dir, "b", a, b)

	out, _, err := r.Generate(context.Background(), Request{Prompt: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Text != "from b" {
		t.Errorf("preferred provider ignored: got %q", out.Text)
	}
}

func TestResolverFallsBackWhenNoneAvailable(t *testing.T) {
	dir := t.TempDir()
	a := &fakeGen{id: "a", available: false}
	r := NewResolver(dir, "", a)

	if _, _, err := r.Generate(context.Background(), Request{Prompt: "x"}); err == nil {
		t.Error("expected error when no provider is available")
	}
}

func TestANSIStrippedFromOutput(t *testing.T) {
	// Simulate a CLI that leaks cursor-movement codes into stdout.
	dirty := "REST\x1b[2D\x1b[K API summary\x1b[1G done"
	clean := ansiEscape.ReplaceAllString(dirty, "")
	if strings.ContainsRune(clean, '\x1b') {
		t.Fatalf("escape byte survived: %q", clean)
	}
	if clean != "REST API summary done" {
		t.Errorf("got %q", clean)
	}
}

func TestResolverDifferentPromptsDifferentKeys(t *testing.T) {
	dir := t.TempDir()
	fake := &fakeGen{id: "fake", available: true, text: "out"}
	r := NewResolver(dir, "", fake)

	_, r1, _ := r.Generate(context.Background(), Request{Prompt: "one"})
	_, r2, _ := r.Generate(context.Background(), Request{Prompt: "two"})
	if r1.CacheKey == r2.CacheKey {
		t.Error("different prompts must produce different cache keys")
	}
	if fake.calls != 2 {
		t.Errorf("expected 2 model calls for 2 distinct prompts, got %d", fake.calls)
	}
}
