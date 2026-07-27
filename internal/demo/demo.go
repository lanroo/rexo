// Package demo runs REXO's first probabilistic pipeline: a four-step "mini-aula"
// (mini-lesson) generated from a single topic. It is the experience a new user
// gets from `rexo demo <topic>` — real AI orchestration, not a deterministic
// placeholder. Each step is a text.generate call served through the Economy
// Engine, so a second run of the same topic is fully cache-served and free.
//
// This pipeline is deliberately separate from the deterministic kernel: its
// steps are probabilistic and are not replay-verified against a model. See
// docs/adr/0005.
package demo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lanroo/rexo/internal/providers"
)

// Options configures a demo run.
type Options struct {
	Topic      string
	ProjectDir string
	Provider   string // preferred provider id (claude-code, codex, ollama); optional
	Model      string // ollama model override; empty means autodetect
	Lang       string // output language for every step; empty means English
}

// StepResult is the outcome of one pipeline stage, including the Economy Engine
// signal (which provider ran and whether it was served from cache).
type StepResult struct {
	ID         string
	Title      string
	Provider   string
	Model      string
	CacheHit   bool
	DurationMS int64
	Text       string
}

// Result is the finished mini-lesson plus per-step provenance.
type Result struct {
	Topic       string
	OutputPath  string
	Steps       []StepResult
	FullyCached bool
}

type stage struct {
	id     string
	title  string
	system string
	build  func(topic string, prev map[string]string) string
}

// concise is appended to every system prompt. Small local models ramble and
// leak their reasoning; a firm instruction narrows the output and brings weaker
// providers closer to a strong one's directness.
const concise = " Be concise and direct. Output only what is asked — no preamble, no meta-commentary, and no explanation of your reasoning."

// languageDirective forces one output language across every step. Choosing the
// language once (rather than "answer in the language of the topic") is what
// stops a run from drifting — e.g. an English summary followed by Portuguese
// slides. Unknown codes pass through so any language name still works.
func languageDirective(lang string) string {
	return " Write your entire response in " + languageName(lang) + "."
}

func languageName(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "", "en", "en-us", "english":
		return "English"
	case "pt", "pt-br", "portuguese", "português":
		return "Portuguese"
	case "es", "spanish", "español":
		return "Spanish"
	case "fr", "french", "français":
		return "French"
	default:
		return lang
	}
}

// pipeline is the fixed four-step mini-lesson. Later stages consume earlier
// outputs, which is what makes this orchestration rather than a single call.
var pipeline = []stage{
	{
		id:     "summary",
		title:  "Summary",
		system: "You are a concise technical educator.",
		build: func(topic string, _ map[string]string) string {
			return fmt.Sprintf("Write a 3-sentence, plain-language summary of the topic %q for a complete beginner. Output only the summary.", topic)
		},
	},
	{
		id:     "objectives",
		title:  "Learning objectives",
		system: "You are an instructional designer.",
		build: func(topic string, prev map[string]string) string {
			return fmt.Sprintf("Given this summary:\n\n%s\n\nWrite exactly 3 learning objectives for a short lesson on %q. Output a numbered list, one objective per line, and nothing else.", prev["summary"], topic)
		},
	},
	{
		id:     "outline",
		title:  "Slide outline",
		system: "You are a curriculum designer.",
		build: func(topic string, prev map[string]string) string {
			return fmt.Sprintf("Create a 5-slide outline for a lesson on %q. Each slide: a title and a one-line description. Cover these learning objectives:\n\n%s\n\nOutput a numbered list of exactly 5 slides and nothing else.", topic, prev["objectives"])
		},
	},
	{
		id:     "quiz",
		title:  "Quiz question",
		system: "You are an assessment writer.",
		build: func(topic string, prev map[string]string) string {
			return fmt.Sprintf("Write exactly 1 multiple-choice question (4 options labelled A-D) that tests these objectives:\n\n%s\n\nMark the correct option clearly. Output only the question, its options, and the correct answer.", prev["objectives"])
		},
	},
}

// Run executes the pipeline and writes the assembled mini-lesson to
// <ProjectDir>/.rexo/demo/<slug>.md, returning its path and per-step details.
func Run(ctx context.Context, opt Options) (*Result, error) {
	topic := strings.TrimSpace(opt.Topic)
	if topic == "" {
		return nil, fmt.Errorf("demo requires a topic, e.g. rexo demo \"REST APIs\"")
	}

	// An explicit --model choice is threaded to the ollama adapter through the
	// same env var it already honours, so the flag wins over autodetection.
	if strings.TrimSpace(opt.Model) != "" {
		os.Setenv("REXO_OLLAMA_MODEL", opt.Model)
	}

	resolver := providers.DefaultResolver(opt.ProjectDir, opt.Provider)
	if len(resolver.Available()) == 0 {
		return nil, fmt.Errorf("no AI provider found — install one of: claude (Claude Code), codex, or ollama, then run `rexo doctor`")
	}

	outputs := make(map[string]string, len(pipeline))
	result := &Result{Topic: topic, FullyCached: true}

	lang := languageDirective(opt.Lang)
	for _, st := range pipeline {
		req := providers.Request{System: st.system + lang + concise, Prompt: st.build(topic, outputs)}
		out, res, err := resolver.Generate(ctx, req)
		if err != nil {
			return result, fmt.Errorf("step %q: %w", st.id, err)
		}
		outputs[st.id] = out.Text
		if !res.CacheHit {
			result.FullyCached = false
		}
		result.Steps = append(result.Steps, StepResult{
			ID:         st.id,
			Title:      st.title,
			Provider:   res.Provider,
			Model:      out.ModelID,
			CacheHit:   res.CacheHit,
			DurationMS: res.DurationMS,
			Text:       out.Text,
		})
	}

	path, err := writeLesson(opt.ProjectDir, topic, result.Steps)
	if err != nil {
		return result, err
	}
	result.OutputPath = path
	return result, nil
}

func writeLesson(projectDir, topic string, steps []StepResult) (string, error) {
	dir := filepath.Join(projectDir, ".rexo", "demo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create demo dir: %w", err)
	}

	var b strings.Builder
	// The topic is the heading and the byline is English tool-provenance
	// metadata, so the file's framing never contradicts the chosen content
	// language (the old hardcoded "Mini-aula:" clashed with English lessons).
	fmt.Fprintf(&b, "# %s\n\n", topic)
	provider := "unknown"
	if len(steps) > 0 {
		provider = steps[0].Provider
	}
	fmt.Fprintf(&b, "_Mini-lesson generated by REXO — capability `text.generate@1`, provider `%s`._\n\n", provider)
	for _, s := range steps {
		fmt.Fprintf(&b, "## %s\n\n%s\n\n", s.Title, strings.TrimSpace(s.Text))
	}

	path := filepath.Join(dir, slug(topic)+".md")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return "", fmt.Errorf("write lesson: %w", err)
	}
	return path, nil
}

// slug turns a topic into a safe, short filename stem.
func slug(topic string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(topic) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	s := strings.Trim(b.String(), "-")
	if s == "" {
		s = "mini-aula"
	}
	if len(s) > 60 {
		s = strings.Trim(s[:60], "-")
	}
	return s
}
