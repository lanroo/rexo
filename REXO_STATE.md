# REXO state

- Product name: REXO
- Expansion: Runtime for Execution & eXchange Orchestration
- Naming history: AIOS (codename) → AIOREN → REXO; AIOS/AIOREN survive only in historical documents. See ADR 0004.
- Current release target: v0.0.1
- Current phase: Phase 0 — Public Foundation
- Branch: `codex/foundation-v0.0.1`
- Runtime status: not implemented
- CLI status: `version`, `doctor`, and `init` implemented
- Supported targets: Windows amd64, macOS amd64/arm64, Linux amd64/arm64
- License decision: Apache-2.0
- Repository: `github.com/lanroo/rexo` (default branch `main`)
- Publication status: published as a **private** repo on 2026-07-26; public visibility still pending owner decision
- CI status: green on Windows, macOS, and Linux (test, vet, build, smoke)
- Release status: `v0.0.1` published with binaries + checksums for all 5 targets
- Name clearance: preliminary search only; formal legal review pending

## Next acceptance gate

1. Contracts validate.
2. Tests and vet pass.
3. Cross-platform CI is configured.
4. Release archives and checksums can be generated.
5. Windows smoke test passes locally.
6. GitHub-hosted macOS and Linux tests pass before the first public release.
