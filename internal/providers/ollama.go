package providers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// Ollama generates text through a local Ollama install (`ollama run <model>`).
// It needs no API key and runs offline, at the cost of the user having pulled a
// model first. The model is chosen, in order: an explicit choice
// (--model / REXO_OLLAMA_MODEL), else the first model `ollama list` reports.
type Ollama struct {
	cli   cliRunner
	model string // explicit model; empty means "autodetect from `ollama list`"
}

func NewOllama() *Ollama {
	return &Ollama{cli: cliRunner{bin: "ollama"}, model: os.Getenv("REXO_OLLAMA_MODEL")}
}

func (o *Ollama) ID() string { return "ollama" }

func (o *Ollama) Detect() bool { return o.cli.available() }

// CacheVariant returns the explicitly-chosen model (or "") so the cache keys
// different ollama models apart without an extra `ollama list` call on hits.
// Autodetected runs share a key, which is fine: autodetection is deterministic
// per machine.
func (o *Ollama) CacheVariant() string { return o.model }

func (o *Ollama) Generate(ctx context.Context, req Request) (*Result, error) {
	model, err := o.resolveModel(ctx)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	// --nowordwrap stops ollama from streaming live cursor-movement escape
	// codes into stdout; --hidethinking drops the reasoning trace that models
	// like deepseek-r1 emit, leaving just the answer. Prompt goes on stdin so
	// long prompts don't hit argument limits.
	args := []string{"run", model, "--hidethinking", "--nowordwrap"}
	out, err := o.cli.run(ctx, args, foldSystem(req.System, req.Prompt))
	if err != nil {
		return nil, err
	}
	return &Result{
		Text:       out,
		ProviderID: o.ID(),
		ModelID:    model,
		Duration:   time.Since(start),
	}, nil
}

// resolveModel returns the explicit model if one was chosen, otherwise the
// first model reported by `ollama list`. Autodetection lets a fresh install
// work without the user knowing a model name; an explicit choice always wins.
func (o *Ollama) resolveModel(ctx context.Context) (string, error) {
	if strings.TrimSpace(o.model) != "" {
		return o.model, nil
	}
	out, err := o.cli.run(ctx, []string{"list"}, "")
	if err != nil {
		return "", fmt.Errorf("could not list ollama models: %w", err)
	}
	models := parseOllamaList(out)
	if len(models) == 0 {
		return "", fmt.Errorf("ollama has no models installed — run e.g. `ollama pull llama3.2`, or pass --model")
	}
	return models[0], nil
}

// parseOllamaList extracts model names (the first column) from `ollama list`
// output, skipping the header row and any blank lines.
func parseOllamaList(out string) []string {
	var models []string
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] == "NAME" {
			continue
		}
		models = append(models, fields[0])
	}
	return models
}
