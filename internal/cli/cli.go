package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"runtime"

	"github.com/lanroo/rexo/internal/demo"
	"github.com/lanroo/rexo/internal/doctor"
	"github.com/lanroo/rexo/internal/kernel"
	"github.com/lanroo/rexo/internal/project"
	"github.com/lanroo/rexo/internal/studio"
	"github.com/lanroo/rexo/internal/workflow"
)

type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func Run(args []string, stdout, stderr io.Writer, build BuildInfo) int {
	if len(args) == 0 {
		printWelcome(stdout, build)
		return 0
	}

	switch args[0] {
	case "help", "-h", "--help":
		printHelp(stdout)
		return 0
	case "version":
		return runVersion(args[1:], stdout, stderr, build)
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "run":
		return runRun(args[1:], stdout, stderr)
	case "demo":
		return runDemo(args[1:], stdout, stderr)
	case "studio":
		return runStudio(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		printHelp(stderr)
		return 2
	}
}

func runVersion(args []string, stdout, stderr io.Writer, build BuildInfo) int {
	flags := flag.NewFlagSet("version", flag.ContinueOnError)
	flags.SetOutput(stderr)
	asJSON := flags.Bool("json", false, "print machine-readable JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "version does not accept positional arguments")
		return 2
	}

	if *asJSON {
		payload := struct {
			BuildInfo
			OS   string `json:"os"`
			Arch string `json:"arch"`
		}{build, runtime.GOOS, runtime.GOARCH}
		if err := json.NewEncoder(stdout).Encode(payload); err != nil {
			fmt.Fprintf(stderr, "write version: %v\n", err)
			return 1
		}
		return 0
	}

	fmt.Fprintf(stdout, "REXO %s (%s/%s, commit %s, built %s)\n",
		build.Version, runtime.GOOS, runtime.GOARCH, build.Commit, build.Date)
	return 0
}

func runDoctor(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("doctor", flag.ContinueOnError)
	flags.SetOutput(stderr)
	asJSON := flags.Bool("json", false, "print machine-readable JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "doctor does not accept positional arguments")
		return 2
	}

	report := doctor.Run()
	if *asJSON {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintf(stderr, "write report: %v\n", err)
			return 1
		}
	} else {
		doctor.Print(stdout, report)
	}
	if !report.Healthy {
		return 1
	}
	return 0
}

func runInit(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: rexo init <project-directory>")
		return 2
	}

	result, err := project.Init(flags.Arg(0))
	if err != nil {
		if errors.Is(err, project.ErrTargetExists) {
			fmt.Fprintf(stderr, "cannot initialize project: %v\n", err)
			return 2
		}
		fmt.Fprintf(stderr, "initialize project: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Created REXO project %q at %s\n", result.Name, result.Path)
	fmt.Fprintf(stdout, "Next:\n  cd %s\n  rexo run workflow.json\n", result.Path)
	return 0
}

func runRun(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	flags.SetOutput(stderr)
	projectDir := flags.String("project", ".", "project directory that holds .rexo")
	replay := flags.String("replay", "", "verify determinism against a recorded run id")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: rexo run <workflow.json> [--project <dir>] [--replay <run-id>]")
		return 2
	}

	wf, err := workflow.Load(flags.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "load workflow: %v\n", err)
		return 1
	}
	opt := kernel.Options{ProjectDir: *projectDir}

	if *replay != "" {
		res, err := kernel.Replay(wf, opt, *replay)
		if err != nil {
			fmt.Fprintf(stderr, "replay failed: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "Replay OK: %d task(s) reproduced identical output for run %s\n", res.Checked, res.RunID)
		return 0
	}

	trace, err := kernel.Run(wf, opt)
	if err != nil {
		fmt.Fprintf(stderr, "run failed: %v\n", err)
		if trace != nil {
			fmt.Fprintf(stderr, "diagnosable trace: %s\n",
				filepath.Join(*projectDir, ".rexo", "runs", trace.RunID, "trace.json"))
		}
		return 1
	}

	fmt.Fprintf(stdout, "Run %s: %s (%d task(s))\n", trace.RunID, trace.Status, len(trace.Tasks))
	for _, path := range trace.Outputs {
		fmt.Fprintf(stdout, "  wrote:     %s\n", filepath.Join(*projectDir, filepath.FromSlash(path)))
	}
	fmt.Fprintf(stdout, "  trace:     %s\n", filepath.Join(*projectDir, ".rexo", "runs", trace.RunID, "trace.json"))
	fmt.Fprintf(stdout, "  artifacts: %s\n", filepath.Join(*projectDir, ".rexo", "artifacts"))
	fmt.Fprintf(stdout, "Replay to verify: rexo run %s --replay %s\n", flags.Arg(0), trace.RunID)
	return 0
}

func runDemo(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("demo", flag.ContinueOnError)
	flags.SetOutput(stderr)
	projectDir := flags.String("project", ".", "project directory that holds .rexo")
	provider := flags.String("provider", "", "preferred provider: claude-code, codex, or ollama")
	model := flags.String("model", "", "ollama model to use (default: autodetect the first installed)")
	lang := flags.String("lang", "en", "output language for the whole lesson, e.g. en, pt, es, fr")

	// Go's flag package stops parsing at the first positional argument, so
	// `demo "topic" --provider x` would drop the flags. Pull the topic out and
	// parse any flags that trail it, so both orderings work.
	if err := flags.Parse(args); err != nil {
		return 2
	}
	rest := flags.Args()
	if len(rest) == 0 {
		fmt.Fprintln(stderr, "usage: rexo demo <topic> [--provider <id>] [--model <name>] [--lang <code>] [--project <dir>]")
		fmt.Fprintln(stderr, "example: rexo demo \"REST APIs\"")
		return 2
	}
	topic := rest[0]
	if len(rest) > 1 {
		if err := flags.Parse(rest[1:]); err != nil {
			return 2
		}
		if flags.NArg() != 0 {
			fmt.Fprintf(stderr, "unexpected extra arguments: %v\n", flags.Args())
			return 2
		}
	}
	fmt.Fprintf(stdout, "Generating a mini-lesson on %q ...\n\n", topic)

	res, err := demo.Run(context.Background(), demo.Options{
		Topic:      topic,
		ProjectDir: *projectDir,
		Provider:   *provider,
		Model:      *model,
		Lang:       *lang,
	})
	if err != nil {
		fmt.Fprintf(stderr, "demo failed: %v\n", err)
		return 1
	}

	// Per-step provenance makes the Economy Engine visible: a cached step did
	// not call the model at all.
	calls := 0
	for _, s := range res.Steps {
		who := s.Provider
		if s.Model != "" {
			who = s.Provider + ":" + s.Model
		}
		status := fmt.Sprintf("generated via %s in %dms", who, s.DurationMS)
		if s.CacheHit {
			status = fmt.Sprintf("reused from cache (%s) — 0 model calls", who)
		} else {
			calls++
		}
		fmt.Fprintf(stdout, "  [%s] %s\n", s.Title, status)
	}

	fmt.Fprintf(stdout, "\nMini-lesson written to: %s\n", res.OutputPath)
	if res.FullyCached {
		fmt.Fprintln(stdout, "Every step was served from cache — this run cost 0 model calls.")
	} else {
		fmt.Fprintf(stdout, "%d of %d step(s) called a model; the rest were cached.\n", calls, len(res.Steps))
		fmt.Fprintln(stdout, "Run the same command again — it will be instant and free (Economy Engine).")
	}
	return 0
}

func runStudio(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("studio", flag.ContinueOnError)
	flags.SetOutput(stderr)
	projectDir := flags.String("project", ".", "project directory that holds .rexo")
	port := flags.String("port", "4747", "port to serve the Studio UI on")
	provider := flags.String("provider", "", "preferred provider: claude-code, codex, or ollama")
	model := flags.String("model", "", "ollama model to use (default: autodetect)")
	noOpen := flags.Bool("no-open", false, "do not open the browser automatically")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	err := studio.Serve(studio.Options{
		ProjectDir: *projectDir,
		Addr:       "127.0.0.1:" + *port,
		Provider:   *provider,
		Model:      *model,
		Open:       !*noOpen,
	}, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "studio: %v\n", err)
		return 1
	}
	return 0
}

func printWelcome(w io.Writer, build BuildInfo) {
	version := build.Version
	if version == "" {
		version = "dev"
	}
	fmt.Fprintf(w, `============================================================
  REXO — Runtime for Execution & eXchange Orchestration
  %s · public foundation
============================================================

REXO is a command-line tool, not an app with buttons. That is
why double-clicking only flashed a black window — this is normal
and safe, not a virus. You use REXO by typing a command below.

Try one of these:

  rexo studio              Open the visual Studio in your browser
  rexo demo "REST APIs"    Generate an AI mini-lesson (needs an AI CLI)
  rexo doctor              Check that your machine is ready
  rexo init my-project     Create a new REXO project
  rexo run workflow.json   Execute a deterministic workflow
  rexo version             Show version details
  rexo help                Show the full command list

On Windows, type it as:  .\rexo.exe demo "REST APIs"

Docs & downloads: https://github.com/lanroo/rexo

Tip: "rexo demo" uses your Claude Code, Codex, or Ollama CLI. Run the same
topic twice — the second run is instant and free (the Economy Engine caches it).
`, version)
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, `REXO — Runtime for Execution & eXchange Orchestration

Usage:
  rexo <command> [options]

Commands:
  version   Show version and platform information
  doctor    Validate this machine for REXO (lists available AI providers)
  init      Create a new REXO project
  run       Execute a deterministic workflow (--replay verifies determinism)
  demo      Generate an AI mini-lesson on a topic (needs Claude Code, Codex, or Ollama)
  studio    Open a local visual UI in your browser (describe a goal, watch it run)
  help      Show this help

"rexo run" is deterministic and replayable. "rexo demo" calls a language model
through your installed CLI and caches results so repeats are free.`)
}
