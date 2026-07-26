# REXO state

- Product name: REXO
- Expansion: Runtime for Execution & eXchange Orchestration
- Naming history: AIOS (codename) → AIOREN → REXO; AIOS/AIOREN survive only in historical documents. See ADR 0004.
- Current release target: v0.0.3
- Current phase: Phase 1 — Walking Skeleton (in progress)
- Branch: `codex/foundation-v0.0.1`
- Runtime status: deterministic walking skeleton implemented (`rexo run` executes a workflow DAG, stores content-addressed artifacts + execution trace + JSONL event log, and `--replay` verifies determinism). No LLM/network yet. See docs/roadmap/phase-1-walking-skeleton.md.
- CLI status: `version`, `doctor`, `init`, and `run` implemented
- Contracts: project, workflow (now with step `with` inputs), task-envelope, artifact-manifest, execution-trace
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
