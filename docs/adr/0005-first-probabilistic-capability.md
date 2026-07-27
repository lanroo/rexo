# ADR 0005: First probabilistic capability via CLI providers

- Status: accepted
- Date: 2026-07-26

## Context

Phase 1 delivered a deterministic walking skeleton: the kernel runs a workflow
DAG, stores content-addressed artifacts, and can replay a run to prove
determinism. Its providers are pure functions — "no network, no clock, no
randomness."

Phase 2 introduces the first capability whose output comes from a language
model, so a user can install REXO and immediately experience AI orchestration,
not just a deterministic demo.

## Decision

**One capability, `text.generate@1`,** defined as a versioned manifest before
any implementation (`capabilities/text.generate.json`), conforming to the
existing `contracts/capability.schema.json` meta-schema.

**Providers are local tools the user already trusts, not a REXO-managed API
key.** The adapters drive `claude`, `codex`, and Ollama — the latter both via
its HTTP API (`ollama-api`, preferred: it exposes temperature and returns clean
responses) and its CLI (`ollama`, fallback). This follows the constitution
("REXO abstracts MCPs, APIs, CLIs, local apps and cloud services as capability
providers") and has three practical benefits: REXO never handles API keys (each
tool owns its own auth), zero friction for users who already have these tools,
and it positions REXO as the orchestrator of these tools rather than a
competitor to them. Two Ollama transports behind one capability is itself a
proof of provider independence — the workflow never changes.

**Probabilistic execution is kept separate from deterministic execution.**
The deterministic `Provider` registry in `internal/kernel` is unchanged and
keeps its replay guarantee. Generation lives in a new `internal/providers`
package with its own `Generator` interface. A workflow step is deterministic or
probabilistic; the two never share a registry.

**The Economy Engine is the reason to reuse, not just an optimization.**
Generation outputs are cached content-addressed by
`hash(capability, provider, prompt, model_hint)`. A cache hit skips the model
call entirely. The demo proves this: the second run of the same topic costs
zero calls and returns instantly. Replay of a probabilistic run verifies
against the cached artifact, never by re-calling the model.

## Consequences

- `rexo run` gains awareness of probabilistic steps; `--replay` must compare
  against cached outputs rather than re-invoking a non-deterministic provider.
- `rexo doctor` reports which provider CLIs are present on PATH.
- Provider independence is testable from day one: switching provider must not
  change a workflow definition.
- Cost, latency, provider id, and cache hits are recorded in the trace so the
  economic behavior is observable, per ADR 0003.
- Name/really-probabilistic boundary is explicit, so Phase 1's determinism
  invariant is never silently violated.
