package providers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllamaAPIGenerate(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"model":"test-model","response":"  a clean answer  "}`))
	}))
	defer srv.Close()

	o := newOllamaAPI(srv.URL, "test-model", 0.3)
	out, err := o.Generate(context.Background(), Request{Prompt: "hi", System: "be nice"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Text != "a clean answer" {
		t.Errorf("text = %q (should be trimmed)", out.Text)
	}
	if out.ProviderID != "ollama-api" || out.ModelID != "test-model" {
		t.Errorf("provider/model = %s/%s", out.ProviderID, out.ModelID)
	}
	// The request must carry the temperature option and the system prompt.
	opts, _ := gotBody["options"].(map[string]any)
	if opts == nil || opts["temperature"] != 0.3 {
		t.Errorf("temperature not sent: %v", gotBody["options"])
	}
	if gotBody["system"] != "be nice" {
		t.Errorf("system not sent: %v", gotBody["system"])
	}
	if gotBody["stream"] != false {
		t.Errorf("stream should be false: %v", gotBody["stream"])
	}
}

func TestOllamaAPIAutodetectModel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tags":
			_, _ = w.Write([]byte(`{"models":[{"name":"first:latest"},{"name":"second:latest"}]}`))
		case "/api/generate":
			_, _ = w.Write([]byte(`{"model":"first:latest","response":"ok"}`))
		}
	}))
	defer srv.Close()

	o := newOllamaAPI(srv.URL, "", 0.3) // no explicit model → autodetect
	out, err := o.Generate(context.Background(), Request{Prompt: "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if out.ModelID != "first:latest" {
		t.Errorf("autodetect picked %q, want first:latest", out.ModelID)
	}
}

func TestOllamaAPIDetect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	if !newOllamaAPI(srv.URL, "", 0.3).Detect() {
		t.Error("Detect should be true for a reachable server")
	}
	// An unroutable host must not hang or panic.
	if newOllamaAPI("http://127.0.0.1:1", "", 0.3).Detect() {
		t.Error("Detect should be false for an unreachable server")
	}
}

func TestOllamaAPICacheVariantDistinguishes(t *testing.T) {
	a := newOllamaAPI("http://x", "m1", 0.3)
	b := newOllamaAPI("http://x", "m1", 0.9)
	c := newOllamaAPI("http://x", "m2", 0.3)
	if a.CacheVariant() == b.CacheVariant() {
		t.Error("different temperatures must yield different cache variants")
	}
	if a.CacheVariant() == c.CacheVariant() {
		t.Error("different models must yield different cache variants")
	}
	if strings.TrimSpace(a.CacheVariant()) == "" {
		t.Error("cache variant should not be empty")
	}
}
