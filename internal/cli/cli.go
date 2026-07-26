package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"runtime"

	"github.com/lanroo/rexo/internal/doctor"
	"github.com/lanroo/rexo/internal/project"
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
	fmt.Fprintf(stdout, "Next: cd %s\n", result.Path)
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

  rexo doctor              Check that your machine is ready
  rexo init my-project     Create a new REXO project
  rexo version             Show version details
  rexo help                Show the full command list

On Windows, type it as:  .\rexo.exe doctor

Docs & downloads: https://github.com/lanroo/rexo

Note: this %s public foundation does not call an AI model yet.
`, version, version)
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, `REXO — Runtime for Execution & eXchange Orchestration

Usage:
  rexo <command> [options]

Commands:
  version   Show version and platform information
  doctor    Validate this machine for REXO
  init      Create a new REXO project
  help      Show this help

This is the v0.0.1 public foundation. It does not call an LLM.`)
}
