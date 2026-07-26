# REXO state

- Product name: REXO
- Expansion: Runtime for Execution & eXchange Orchestration
- Naming history: AIOS (codename) → AIOREN → REXO; AIOS/AIOREN survive only in historical documents. See ADR 0004.
- Current release target: v0.0.2
- Current phase: Phase 0 — Public Foundation
- Branch: `codex/foundation-v0.0.1`
- Runtime status: not implemented
- CLI status: `version`, `doctor`, and `init` implemented
- Supported targets: Windows amd64, macOS amd64/arm64, Linux amd64/arm64
- License decision: Apache-2.0
- Repository: `github.com/lanroo/rexo` (default branch `main`)
- Publication status: **public** on 2026-07-26 (started private the same day, then made public by owner)
- CI status: green on Windows, macOS, and Linux (test, vet, build, smoke)
- Release status: `v0.0.2` public (latest) with binaries + checksums for all 5 targets; version-less asset names (`rexo_<os>_<arch>`) so /releases/latest/download links are stable; anonymous download verified. v0.0.2 adds a friendly no-arg welcome screen and keeps the Windows console open on Explorer double-click.
- Name clearance: NOT cleared. Preliminary search (2026-07-26) found "REXO" is commercially active in software/IT (REXO Solutions LLC, Rexoit, Rexo Group) and REX-based AI-agent names nearby; no exact REXO trademark confirmed, but USPTO/INPI not formally searched. Formal INPI (BR, classes 9/42) + USPTO search still required before any commercial/brand investment.

## Next acceptance gate

1. Contracts validate.
2. Tests and vet pass.
3. Cross-platform CI is configured.
4. Release archives and checksums can be generated.
5. Windows smoke test passes locally.
6. GitHub-hosted macOS and Linux tests pass before the first public release.
