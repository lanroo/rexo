# ADR 0004: Rename the product to REXO

- Status: accepted provisionally
- Date: 2026-07-26
- Supersedes: [ADR 0001](0001-product-name.md)

## Context

ADR 0001 chose **AIOREN**. A follow-up availability sweep across the channels
that actually matter for a developer platform showed AIOREN and every pure
three-letter candidate collided with existing uses, and short candidates like
CREX and REXA collided in the AI/agent space specifically (multiple active
"Rexa" AI platforms, a "crypto exchange" PyPI package for CREX).

Requirements for the name:

- Nothing relevant in the AI / agents / orchestration / runtime category, where
  the project will actually compete for attention and search.
- Preserves the "rex" (king / T-rex / Latin *rex*) connotation of runtime power.
- Pronounceable identically in Portuguese and English.
- Channel collisions, if any, must be resolvable with standard namespacing.

## Decision

Use **REXO**, expanded as **Runtime for Execution & eXchange Orchestration**.

- **R** — Runtime
- **E** — Execution
- **X** — eXchange (interchangeable capabilities via the Tool Gateway)
- **O** — Orchestration

Channel strategy (industry-standard scoping, not a workaround):

- **PyPI:** `rexo` is free — claim it.
- **npm:** publish under a scope — `@rexo/core`, `@rexo/sdk`, `@rexo/cli`
  (register the `@rexo` org) or `@lanroo/rexo`. The unscoped `rexo` package
  belongs to an unrelated project; scopes are separate namespaces.
- **GitHub:** `github.com/lanroo/rexo` for the repo and Go module; an org such as
  `rexoproject` / `getrexo` / `rexoruntime` may be used if needed.
- **Domain:** `rexo.dev` / `rexo.sh` / `getrexo.com` / `rexoruntime.com` — to be
  checked with a registrar before any public launch.

## Consequences

- All source, docs, contracts, CLI (`rexo`), the Go module path, state and
  bootstrap files, tooling, and the architecture PDFs are renamed from AIOREN to
  REXO in a single pass. `AIOREN_STATE.md` → `REXO_STATE.md`,
  `AIOREN_BOOTSTRAP.md` → `REXO_BOOTSTRAP.md`, `cmd/aioren/` → `cmd/rexo/`.
- ADR 0001 is kept as a historical record and marked superseded.
- This is **not** formal trademark clearance. The name is only cleared by
  preliminary web/package/GitHub search; a professional search in target
  jurisdictions is still required before commercial use.
- No repository, npm org, domain, or package is created or published as part of
  this decision — per the activation contract, renaming does not authorize
  publishing or spending.
