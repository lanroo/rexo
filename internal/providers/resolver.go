package providers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Resolution records how the Economy Engine served one generation request:
// which provider was selected, and whether the model was actually called or the
// result came from cache. It is the auditable "why" behind every model call —
// or every avoided one.
type Resolution struct {
	Provider   string `json:"provider"`
	Model      string `json:"model,omitempty"`
	CacheHit   bool   `json:"cache_hit"`
	CacheKey   string `json:"cache_key"`
	DurationMS int64  `json:"duration_ms"`
}

// Resolver selects an available provider and enforces the Economy Engine's
// first rule: reuse before generating. Generations are cached content-addressed
// by (capability, provider, prompt, system, max_tokens), so an identical
// request never pays for a second model call.
type Resolver struct {
	generators []Generator
	cacheDir   string
	preferred  string
}

// NewResolver builds a resolver over the given generators. cacheDir is where
// cached generations live (content-addressed). preferred, when set and
// available, wins provider selection; otherwise the first available generator
// in registration order is used.
func NewResolver(cacheDir, preferred string, gens ...Generator) *Resolver {
	return &Resolver{generators: gens, cacheDir: cacheDir, preferred: preferred}
}

// DefaultResolver wires the built-in providers in priority order and points
// the cache at <projectDir>/.rexo/cache/generations. The Ollama HTTP adapter
// comes before the Ollama CLI adapter: when the server is reachable it is the
// higher-quality path (temperature control, clean responses), and the CLI is
// the fallback when only the binary — not the server — is available.
func DefaultResolver(projectDir, preferred string) *Resolver {
	cacheDir := filepath.Join(projectDir, ".rexo", "cache", "generations")
	return NewResolver(cacheDir, preferred,
		NewClaudeCode(),
		NewCodex(),
		NewOllamaAPI(),
		NewOllama(),
	)
}

// Available returns the ids of every provider usable on this machine.
func (r *Resolver) Available() []string {
	ids := make([]string, 0, len(r.generators))
	for _, g := range r.generators {
		if g.Detect() {
			ids = append(ids, g.ID())
		}
	}
	return ids
}

// selectGenerator picks the provider to use: the preferred one if set and
// available, else the first available in registration order.
func (r *Resolver) selectGenerator() (Generator, error) {
	var first Generator
	for _, g := range r.generators {
		if !g.Detect() {
			continue
		}
		if first == nil {
			first = g
		}
		if r.preferred != "" && g.ID() == r.preferred {
			return g, nil
		}
	}
	if r.preferred != "" && first != nil {
		// Preferred was requested but absent; fall back transparently.
		return first, nil
	}
	if first == nil {
		return nil, fmt.Errorf("no text.generate provider available (install one of: claude, codex, ollama)")
	}
	return first, nil
}

// Generate serves a request through the Economy Engine. It first computes the
// cache key from the selected provider and the request; on a hit it returns the
// cached text with zero model calls. On a miss it calls the model, stores the
// result, and reports the resolution either way.
func (r *Resolver) Generate(ctx context.Context, req Request) (*Result, *Resolution, error) {
	gen, err := r.selectGenerator()
	if err != nil {
		return nil, nil, err
	}

	variant := ""
	if v, ok := gen.(Variant); ok {
		variant = v.CacheVariant()
	}
	key := cacheKey("text.generate", gen.ID(), variant, req)
	res := &Resolution{Provider: gen.ID(), CacheKey: key}

	if cached, ok := r.readCache(key); ok {
		res.CacheHit = true
		res.Model = cached.ModelID
		return cached, res, nil
	}

	out, err := gen.Generate(ctx, req)
	if err != nil {
		return nil, res, err
	}
	res.Model = out.ModelID
	res.DurationMS = out.Duration.Milliseconds()
	if err := r.writeCache(key, out); err != nil {
		return out, res, fmt.Errorf("cache write: %w", err)
	}
	return out, res, nil
}

// cacheKey is a stable content hash of everything that could change the output.
// Provider id (and its variant, e.g. the model) are included so switching
// provider or model does not serve a stale result.
func cacheKey(capability, providerID, variant string, req Request) string {
	payload := struct {
		Capability string `json:"capability"`
		Provider   string `json:"provider"`
		Variant    string `json:"variant"`
		Prompt     string `json:"prompt"`
		System     string `json:"system"`
		MaxTokens  int    `json:"max_tokens"`
	}{capability, providerID, variant, req.Prompt, req.System, req.MaxTokens}
	data, _ := json.Marshal(payload)
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

type cacheEntry struct {
	Text       string `json:"text"`
	ProviderID string `json:"provider_id"`
	ModelID    string `json:"model_id,omitempty"`
}

func (r *Resolver) cachePath(key string) string {
	// Strip the "sha256:" prefix for a clean filename.
	name := key
	if i := len("sha256:"); len(key) > i {
		name = key[i:]
	}
	return filepath.Join(r.cacheDir, name+".json")
}

func (r *Resolver) readCache(key string) (*Result, bool) {
	data, err := os.ReadFile(r.cachePath(key))
	if err != nil {
		return nil, false
	}
	var e cacheEntry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, false
	}
	return &Result{Text: e.Text, ProviderID: e.ProviderID, ModelID: e.ModelID}, true
}

func (r *Resolver) writeCache(key string, out *Result) error {
	if err := os.MkdirAll(r.cacheDir, 0o755); err != nil {
		return err
	}
	e := cacheEntry{Text: out.Text, ProviderID: out.ProviderID, ModelID: out.ModelID}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.cachePath(key), append(data, '\n'), 0o644)
}
