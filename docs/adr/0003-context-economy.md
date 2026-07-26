# ADR 0003: Context efficiency is a system invariant

- Status: accepted
- Date: 2026-07-26

## Decision

Every probabilistic execution receives only objective, required artifacts,
relevant memory, budget, and success criteria. Before a model call, the runtime
must check reusable artifacts, cache, incremental rebuild options, and smaller
adequate providers.

Workers are disposable. Budgets bound tokens, calls, time, cost, and retries.
Observability records consumption and cache behavior.

## Consequences

Conversation transcripts are not the primary persistence mechanism. Artifact
and context schemas become first-class architecture, and optimization can be
measured across projects.
