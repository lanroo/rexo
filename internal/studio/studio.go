// Package studio serves a small local web UI for REXO: describe a goal, watch
// the AI pipeline run live, and read the result — no terminal knowledge needed.
// It uses only the standard library and embeds its single page, so the binary
// stays self-contained and offline.
package studio

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lanroo/rexo/internal/demo"
	"github.com/lanroo/rexo/internal/pipeline"
	"github.com/lanroo/rexo/internal/providers"
)

//go:embed index.html templates/catalog.json
var assets embed.FS

// Options configures the Studio server.
type Options struct {
	ProjectDir string
	Addr       string // host:port, e.g. 127.0.0.1:4747
	Provider   string
	Model      string
	Open       bool // open the browser once the server is listening
}

// Serve starts the Studio HTTP server and blocks until the process stops.
func Serve(opt Options, stdout io.Writer) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		page, err := assets.ReadFile("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(page)
	})

	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		res := providers.DefaultResolver(opt.ProjectDir, opt.Provider)
		writeJSON(w, map[string]any{
			"providers":    res.Available(),
			"ollamaModels": providers.OllamaModels(),
			"capabilities": []string{"text.generate"},
		})
	})

	mux.HandleFunc("/api/templates", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			built, err := builtinTemplates()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			all := append(built, loadUserTemplates()...)
			writeJSON(w, map[string]any{"templates": all})
		case http.MethodPost:
			var t map[string]any
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			name, _ := t["name"].(string)
			steps, _ := t["steps"].([]any)
			if strings.TrimSpace(name) == "" || len(steps) == 0 {
				http.Error(w, "template needs a name and at least one step", http.StatusBadRequest)
				return
			}
			id, _ := t["id"].(string)
			if strings.TrimSpace(id) == "" {
				id = name
			}
			id = safeID(id)
			dir := userTemplatesDir()
			if dir == "" {
				http.Error(w, "cannot resolve home directory", http.StatusInternalServerError)
				return
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t["id"] = id
			t["source"] = "user"
			data, _ := json.MarshalIndent(t, "", "  ")
			if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0o644); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, map[string]any{"ok": true, "id": id})
		default:
			http.Error(w, "GET or POST", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/demo", func(w http.ResponseWriter, r *http.Request) {
		topic := r.URL.Query().Get("topic")
		lang := r.URL.Query().Get("lang")
		provider := r.URL.Query().Get("provider")
		if provider == "" {
			provider = opt.Provider
		}
		model := r.URL.Query().Get("model")
		if model == "" {
			model = opt.Model
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		send := func(event string, payload any) {
			b, _ := json.Marshal(payload)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b)
			flusher.Flush()
		}

		result, err := demo.Run(r.Context(), demo.Options{
			Topic:      topic,
			ProjectDir: opt.ProjectDir,
			Provider:   provider,
			Model:      model,
			Lang:       lang,
			Progress:   func(s demo.StepResult) { send("step", s) },
		})
		if err != nil {
			send("error", map[string]string{"message": err.Error()})
			return
		}
		send("done", map[string]any{
			"outputPath":  result.OutputPath,
			"fullyCached": result.FullyCached,
		})
	})

	mux.HandleFunc("/api/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Topic    string          `json:"topic"`
			Lang     string          `json:"lang"`
			Provider string          `json:"provider"`
			Model    string          `json:"model"`
			Steps    []pipeline.Step `json:"steps"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-ndjson")

		send := func(event string, data any) {
			line, _ := json.Marshal(map[string]any{"event": event, "data": data})
			_, _ = w.Write(append(line, '\n'))
			flusher.Flush()
		}

		provider := body.Provider
		if provider == "" {
			provider = opt.Provider
		}
		model := body.Model
		if model == "" {
			model = opt.Model
		}

		result, err := pipeline.Run(r.Context(), pipeline.Options{
			Topic:      body.Topic,
			Lang:       body.Lang,
			Provider:   provider,
			Model:      model,
			ProjectDir: opt.ProjectDir,
			Steps:      body.Steps,
			Progress:   func(s pipeline.StepResult) { send("step", s) },
		})
		if err != nil {
			send("error", map[string]string{"message": err.Error()})
			return
		}
		send("done", map[string]any{"outputPath": result.OutputPath, "fullyCached": result.FullyCached})
	})

	ln, err := net.Listen("tcp", opt.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", opt.Addr, err)
	}
	url := "http://" + ln.Addr().String()
	fmt.Fprintf(stdout, "REXO Studio is running at %s\n", url)
	fmt.Fprintln(stdout, "Open it in your browser. Press Ctrl+C to stop.")
	if opt.Open {
		openBrowser(url)
	}
	return http.Serve(ln, mux)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// builtinTemplates reads the curated catalog embedded in the binary and tags
// each entry as built-in so the UI can distinguish it from a user's own.
func builtinTemplates() ([]map[string]any, error) {
	data, err := assets.ReadFile("templates/catalog.json")
	if err != nil {
		return nil, err
	}
	var list []map[string]any
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, t := range list {
		t["source"] = "built-in"
	}
	return list, nil
}

// userTemplatesDir is where imported/saved templates live, next to the rest of
// REXO's user state. Empty if the home directory cannot be resolved.
func userTemplatesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".rexo", "templates")
}

// loadUserTemplates reads every *.json template the user has imported or saved.
// Unreadable or malformed files are skipped rather than failing the catalog.
func loadUserTemplates() []map[string]any {
	dir := userTemplatesDir()
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []map[string]any
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var t map[string]any
		if json.Unmarshal(data, &t) != nil {
			continue
		}
		t["source"] = "user"
		out = append(out, t)
	}
	return out
}

// safeID reduces an arbitrary string to a filename-safe slug ([a-z0-9-]). This
// also prevents path traversal when a saved template becomes <id>.json.
func safeID(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
		} else if !dash {
			b.WriteByte('-')
			dash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		out = "template"
	}
	if len(out) > 60 {
		out = strings.Trim(out[:60], "-")
	}
	return out
}
