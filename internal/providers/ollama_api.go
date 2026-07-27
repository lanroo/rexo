package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultOllamaHost = "http://localhost:11434"
	// A low default temperature keeps the mini-lesson focused; local models
	// ramble far more at the CLI default. Override with REXO_OLLAMA_TEMPERATURE.
	defaultOllamaTemp = 0.3
)

// OllamaAPI generates text through the Ollama HTTP API instead of the CLI.
// Talking to the API directly unlocks generation options the `ollama run`
// command does not expose — chiefly temperature — and returns a single clean
// JSON response with no streaming/cursor artefacts. It is the higher-quality
// path to a local model; the CLI adapter (ID "ollama") remains as a fallback
// for when the server is not reachable.
type OllamaAPI struct {
	host   string
	model  string // explicit model; empty means autodetect via /api/tags
	temp   float64
	client *http.Client
}

func NewOllamaAPI() *OllamaAPI {
	host := os.Getenv("REXO_OLLAMA_HOST")
	if strings.TrimSpace(host) == "" {
		host = defaultOllamaHost
	}
	temp := defaultOllamaTemp
	if v := os.Getenv("REXO_OLLAMA_TEMPERATURE"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			temp = f
		}
	}
	return newOllamaAPI(host, os.Getenv("REXO_OLLAMA_MODEL"), temp)
}

// newOllamaAPI is the injectable constructor used by tests.
func newOllamaAPI(host, model string, temp float64) *OllamaAPI {
	return &OllamaAPI{
		host:   strings.TrimRight(host, "/"),
		model:  model,
		temp:   temp,
		client: &http.Client{Timeout: 5 * time.Minute},
	}
}

func (o *OllamaAPI) ID() string { return "ollama-api" }

// OllamaModels returns the model names installed on the local Ollama server, or
// nil if the server is unreachable. It lets the UI offer a concrete choice of
// local model instead of only autodetecting the first one.
func OllamaModels() []string {
	host := os.Getenv("REXO_OLLAMA_HOST")
	if strings.TrimSpace(host) == "" {
		host = defaultOllamaHost
	}
	host = strings.TrimRight(host, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/api/tags", nil)
	if err != nil {
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var payload struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil
	}
	names := make([]string, 0, len(payload.Models))
	for _, m := range payload.Models {
		names = append(names, m.Name)
	}
	return names
}

// CacheVariant folds in both the model and the temperature, since either can
// change the output; without this the cache would serve a stale generation
// after a temperature or model change.
func (o *OllamaAPI) CacheVariant() string {
	return fmt.Sprintf("%s@t=%.2f", o.model, o.temp)
}

// Detect reports whether the Ollama server answers on its API port. Unlike the
// CLI adapter, this needs a running server, so we probe it with a short deadline
// rather than just checking PATH.
func (o *OllamaAPI) Detect() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.host+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (o *OllamaAPI) Generate(ctx context.Context, req Request) (*Result, error) {
	model, err := o.resolveModel(ctx)
	if err != nil {
		return nil, err
	}
	start := time.Now()

	payload := map[string]any{
		"model":   model,
		"prompt":  req.Prompt,
		"stream":  false,
		"think":   false, // suppress reasoning traces for models that support it
		"options": map[string]any{"temperature": o.temp},
	}
	if strings.TrimSpace(req.System) != "" {
		payload["system"] = req.System
	}
	if req.MaxTokens > 0 {
		payload["options"].(map[string]any)["num_predict"] = req.MaxTokens
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.host+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama api call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama api %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var out struct {
		Response string `json:"response"`
		Model    string `json:"model"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}
	return &Result{
		Text:       strings.TrimSpace(out.Response),
		ProviderID: o.ID(),
		ModelID:    model,
		Duration:   time.Since(start),
	}, nil
}

// resolveModel returns the explicit model, else the first model from /api/tags.
func (o *OllamaAPI) resolveModel(ctx context.Context) (string, error) {
	if strings.TrimSpace(o.model) != "" {
		return o.model, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.host+"/api/tags", nil)
	if err != nil {
		return "", err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama api not reachable at %s: %w", o.host, err)
	}
	defer resp.Body.Close()
	var payload struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode ollama tags: %w", err)
	}
	if len(payload.Models) == 0 {
		return "", fmt.Errorf("ollama has no models installed — run e.g. `ollama pull llama3.2`, or pass --model")
	}
	return payload.Models[0].Name, nil
}
