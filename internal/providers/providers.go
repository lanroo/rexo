// Package providers implements probabilistic capabilities: capabilities whose
// output comes from a language model and therefore is NOT reproducible byte for
// byte. This is deliberately separate from the deterministic providers in
// internal/kernel, which keep the replay guarantee. A workflow step is one or
// the other; the two never share a registry. See docs/adr/0005.
package providers

import (
	"context"
	"time"
)

// Request is the input to a text.generate capability. It matches the input
// schema in capabilities/text.generate.json.
type Request struct {
	Prompt    string
	System    string
	MaxTokens int
}

// Result is the output of a generation, plus the observability the Economy
// Engine and the trace need (which provider ran, how long it took).
type Result struct {
	Text       string
	ProviderID string
	ModelID    string
	Duration   time.Duration
}

// Generator is a probabilistic capability implementation. Detect reports
// whether this provider is usable on the current machine (e.g. its CLI is on
// PATH); Generate performs the model call. Implementations must not panic when
// their backing tool is absent — Detect gates that.
type Generator interface {
	ID() string
	Detect() bool
	Generate(ctx context.Context, req Request) (*Result, error)
}

// Variant is an optional Generator capability: it contributes an extra
// cache-key component beyond the provider id (e.g. the specific model), so two
// configurations of the same provider don't collide in the cache. Providers
// that don't implement it are cached on provider id alone.
type Variant interface {
	CacheVariant() string
}
