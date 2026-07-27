# REXO state

- Product name: REXO
- Expansion: Runtime for Execution & eXchange Orchestration
- Naming history: AIOS (codename) → AIOREN → REXO; AIOS/AIOREN survive only in historical documents. See ADR 0004.
- Current release target: v0.0.8
- Install channels: Homebrew (`brew install lanroo/tap/rexo`, tap repo `lanroo/homebrew-tap`) and Scoop (`scoop bucket add rexo <repo>` + `scoop install rexo`, bucket in `bucket/`), both auto-published by goreleaser (brews + scoops) each release; also `go install`. Homebrew needs the `HOMEBREW_TAP_TOKEN` repo secret (maintainer-only, one-time).
- Current phase: Phase 2 — First probabilistic capability (in progress, branch `phase-2-first-capability`)
- Branch: `codex/foundation-v0.0.1` (deterministic core); `phase-2-first-capability` (first AI capability)
- Runtime status: deterministic walking skeleton implemented (`rexo run` executes a workflow DAG, stores content-addressed artifacts + execution trace + JSONL event log, and `--replay` verifies determinism). Deterministic providers: text.constant/uppercase/concat/template; a workflow may declare an `outputs` map (path → step id) and the kernel writes real files on success (`from_task` references resolve recursively, e.g. inside a template's `vars`) — the example/init workflow now generates a real `welcome.md`. Phase 2 adds the first probabilistic capability `text.generate@1`, served by four providers (Claude Code CLI, Codex CLI, Ollama HTTP API, Ollama CLI) through a Resolver with a content-addressed generation cache (the Economy Engine). The `ollama-api` adapter is preferred over the `ollama` CLI when the server is reachable (temperature control + clean responses); Ollama model is autodetected from `ollama list`/`/api/tags` or set via `--model`/`REXO_OLLAMA_MODEL`. Cache key includes provider + variant (model/temperature) so switching never serves a stale result. Probabilistic execution lives in `internal/providers` + `internal/demo`, kept fully separate from the deterministic kernel. See docs/adr/0005. Verified live on 2026-07-26: `rexo demo` ran end-to-end against both Claude (logged-in CLI) and Ollama (gpt-oss:20b, qwen3:14b), 2nd run fully cache-served (0 calls).
- CLI status: `version`, `doctor` (now lists AI providers), `init`, `run`, `demo`, and `studio` implemented. `rexo demo "<topic>"` generates a 4-step mini-lesson; re-running the same topic is fully cache-served (0 model calls). `rexo studio` starts a local stdlib-only web UI (`internal/studio`, embedded single page) that streams the `text.generate` pipeline live via SSE — the first graphical/visual surface; opens the browser automatically.
- Contracts: project, workflow (now with step `with` inputs), task-envelope, artifact-manifest, execution-trace, and capability manifest `capabilities/text.generate.json`
- Verification note: Go is not installed on the author's machine, so Phase 2 Go code was written but not locally compiled/tested. CI (`go test`/`vet`/`build` on 3 OS) runs on pull_request — open a PR from `phase-2-first-capability` to verify before merge.
- Supported targets: Windows amd64, macOS amd64/arm64, Linux amd64/arm64
- License decision: Apache-2.0
- Repository: `github.com/lanroo/rexo` (default branch `main`)
- Publication status: **public** on 2026-07-26 (started private the same day, then made public by owner)
- CI status: green on Windows, macOS, and Linux (test, vet, build, smoke)
- Release status: `v0.0.3` public (latest) with binaries + checksums for all 5 targets; version-less asset names (`rexo_<os>_<arch>`) so /releases/latest/download links are stable. Scoop manifest (`scoop/rexo.json`) tracks the latest release. v0.0.2 added the friendly welcome screen + Windows double-click console-hold; v0.0.3 adds `rexo run`.
- Windows Defender: v0.0.2 exe flagged as false positive `Program:Script/Wacapew.A!ml` (ML heuristic); submitted to Microsoft WDSI 2026-07-26 (pending). Workaround while unsigned: Defender folder exclusion for the Scoop dir (`Add-MpPreference -ExclusionPath "$HOME\scoop"`). Code signing is the eventual real fix.
- Name clearance: NOT cleared. Preliminary search (2026-07-26) found "REXO" is commercially active in software/IT (REXO Solutions LLC, Rexoit, Rexo Group) and REX-based AI-agent names nearby; no exact REXO trademark confirmed, but USPTO/INPI not formally searched. Formal INPI (BR, classes 9/42) + USPTO search still required before any commercial/brand investment.

## Next acceptance gate

1. Contracts validate.
2. Tests and vet pass.
3. Cross-platform CI is configured.
4. Release archives and checksums can be generated.
5. Windows smoke test passes locally.
6. GitHub-hosted macOS and Linux tests pass before the first public release.
