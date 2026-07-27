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

	"github.com/lanroo/rexo/internal/demo"
	"github.com/lanroo/rexo/internal/providers"
)

//go:embed index.html
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
