# Getting started

REXO ships as a single self-contained executable. It needs **no** Python, Node,
Docker, or LLM account to run. The first release is intentionally small and does
not call any AI model yet.

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
| `rexo doctor` | Verifies basic machine compatibility. |
| `rexo init <dir>` | Creates a new portable REXO project in `<dir>`. |
| `rexo run <workflow.json>` | Executes a workflow; `--replay <run-id>` verifies determinism. |

Workflows are currently **deterministic** — no LLM or provider selection yet.
That comes in later phases; see the [roadmap](../roadmap/core-v1.md) and the
[Phase 1 design](../roadmap/phase-1-walking-skeleton.md).
