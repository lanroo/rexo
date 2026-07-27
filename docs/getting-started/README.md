# Getting started

REXO ships as a single self-contained executable. It needs **no** Python, Node,
Docker, or LLM account to run. Deterministic workflows (`rexo run`) need no AI
account at all. The optional `rexo demo` command adds AI orchestration by
driving an AI CLI you install separately — see
[Try AI orchestration](#try-ai-orchestration-rexo-demo).

There are two tracks below. Pick the one that matches you:

- **[Track A — Non-developers](#track-a--non-developers)**: install a ready-made
  program and run it. No coding.
- **[Track B — Developers](#track-b--developers)**: build from source and hack
  on it.

---

## Track A — Non-developers

You do **not** need to know how to code.

### Easiest: a package manager

If you have (or don't mind installing) a package manager, this is the cleanest
route — one command, and updates are handled for you:

- **macOS / Linux** ([Homebrew](https://brew.sh)):

  ```bash
  brew install lanroo/tap/rexo
  ```

- **Windows** ([Scoop](https://scoop.sh)): installs in user space and avoids the
  SmartScreen "unrecognized app" prompt.

  ```powershell
  scoop bucket add rexo https://github.com/lanroo/rexo
  scoop install rexo
  ```

Then run `rexo doctor`. If you prefer to just download a file instead, keep
reading.

### Or: download the program

1. Open the project's **Releases** page on GitHub.
2. Under the latest release, download the archive that matches your computer:

   | Your computer | File to download |
   |---|---|
   | Windows | `rexo_*_windows_amd64.zip` |
   | Mac (Apple Silicon: M1/M2/M3/M4) | `rexo_*_darwin_arm64.tar.gz` |
   | Mac (older Intel) | `rexo_*_darwin_amd64.tar.gz` |
   | Linux (most PCs) | `rexo_*_linux_amd64.tar.gz` |

   > Not sure which Mac you have? Click the Apple menu →  **About This Mac**. If
   > it says "Apple M...", pick Apple Silicon.

3. Also download `checksums.txt` (optional, to verify the file is intact).

### 2. Unzip it

- **Windows**: right-click the `.zip` → **Extract All**. You get `rexo.exe`.
- **Mac / Linux**: double-click the `.tar.gz` (or run
  `tar -xzf rexo_*_.tar.gz`). You get a file named `rexo`.

### 3. Open a terminal in that folder

- **Windows**: open the folder in File Explorer, click the address bar, type
  `powershell`, press Enter.
- **Mac**: open **Terminal**, type `cd ` (with a space), then drag the folder
  onto the window, press Enter.
- **Linux**: right-click the folder → **Open Terminal Here**.

### 4. Run it

Windows (PowerShell):

```powershell
.\rexo.exe doctor
.\rexo.exe init my-first-project
```

Mac / Linux:

```bash
./rexo doctor
./rexo init my-first-project
```

- `doctor` checks that your machine is compatible and prints a green report.
- `init my-first-project` creates a new folder called `my-first-project` with a
  ready-to-use project structure — including an example `workflow.json`.

Then run your first workflow:

```bash
cd my-first-project
rexo run workflow.json
```

It executes the example (build "hello world" → uppercase it to "HELLO WORLD"),
writing outputs and a trace under `.rexo/runs/`.

> **macOS security note**: the first time, macOS may say the app is from an
> unidentified developer. Right-click `rexo` → **Open** → **Open**, or go to
> **System Settings → Privacy & Security** and click **Open Anyway**. This is
> normal for programs not yet notarized by Apple.

That's it — you've installed and run REXO. 🎉

---

## Track B — Developers

Requirements: **Go 1.24+** and **Git**.

### Clone and build

```bash
git clone https://github.com/lanroo/rexo.git
cd rexo
go test ./...
go build -o rexo ./cmd/rexo   # produces rexo.exe on Windows
```

### Run

```bash
./rexo version   # build + target platform
./rexo doctor    # machine compatibility check
./rexo init course-builder
cd course-builder
```

### Install onto your PATH (optional)

```bash
go install github.com/lanroo/rexo/cmd/rexo@latest
```

This puts a `rexo` binary in `$(go env GOPATH)/bin`. Make sure that directory is
on your `PATH`, then you can run `rexo` from anywhere.

---

## What `init` creates

Running `rexo init <name>` scaffolds a portable project:

- `rexo.project.json` — versioned project contract;
- `workflow.json` — a runnable example workflow;
- `AGENTS.md` — local instructions for coding agents;
- `REXO_BOOTSTRAP.md` — stable activation sequence;
- `REXO_STATE.md` — compact project state;
- `.rexo/artifacts` — content-addressed outputs;
- `.rexo/runs` — one folder per run (execution trace + event log);
- `.rexo/memory` — curated project memory.

`init` refuses to overwrite an existing path, so it is safe to run.

## Commands at a glance

| Command | What it does |
|---|---|
| `rexo version` | Prints the version and your OS/CPU target. |
| `rexo doctor` | Verifies machine compatibility and lists available AI providers. |
| `rexo init <dir>` | Creates a new portable REXO project in `<dir>`. |
| `rexo run <workflow.json>` | Executes a deterministic workflow; `--replay <run-id>` verifies determinism. |
| `rexo demo "<topic>"` | Generates an AI mini-lesson on a topic through your installed AI CLI. |

`rexo run` is **deterministic** and replayable. `rexo demo` is the first
**probabilistic** capability — it calls a language model. The two are kept
separate on purpose; see the [roadmap](../roadmap/core-v1.md).

---

## Try AI orchestration: `rexo demo`

This is where REXO stops being a deterministic toy and starts orchestrating
real AI. `rexo demo` runs a four-step pipeline — summary → learning objectives
→ slide outline → quiz question — where each step feeds the next. That chaining
is the orchestration; it is not a single chat call.

### 1. Install one AI CLI

REXO does not ship or manage any API key. It drives a command-line tool you
already trust. Install **any one** of these:

| Provider | Install | Cost |
|---|---|---|
| **Claude Code** (`claude`) | [claude.com/claude-code](https://claude.com/claude-code) | Anthropic account |
| **Codex** (`codex`) | OpenAI Codex CLI | OpenAI account |
| **Ollama** (`ollama`) | [ollama.com](https://ollama.com) + `ollama pull <model>` | Free, offline |

Then confirm REXO can see it:

```bash
rexo doctor
```

The `ai-providers` line lists whichever CLIs are on your PATH.

### 2. Generate a mini-lesson

```bash
rexo demo "REST APIs"
```

REXO picks an available provider, runs the four steps, and writes a Markdown
lesson to `.rexo/demo/rest-apis.md`. To force a specific provider:

```bash
rexo demo "REST APIs" --provider ollama
```

With Ollama, REXO auto-detects the first model from `ollama list`. Output
quality tracks the model: a small model rambles, a larger one reads like a real
lesson. Pick one explicitly with `--model`:

```bash
rexo demo "REST APIs" --provider ollama --model gpt-oss:20b
```

REXO reaches Ollama two ways: `ollama-api` (the HTTP server, preferred — it
allows a low temperature for focused output) and `ollama` (the CLI, a fallback
when only the binary is running). When the server is up, the API path is chosen
automatically. Tune it with environment variables:

| Variable | Effect | Default |
|---|---|---|
| `REXO_OLLAMA_MODEL` | Model to use (same as `--model`) | autodetect |
| `REXO_OLLAMA_TEMPERATURE` | Lower = more focused, less rambling | `0.3` |
| `REXO_OLLAMA_HOST` | Ollama server URL | `http://localhost:11434` |

### 3. See the Economy Engine work

Run the **exact same command again**:

```bash
rexo demo "REST APIs"
```

Every step now reports *reused from cache — 0 model calls*, and the result is
instant. REXO content-addresses each generation, so an identical request never
pays twice. That is the Economy Engine: reuse before generating.
