# Phase 1 — Walking Skeleton (design)

**Goal:** execute one deterministic workflow end to end with `rexo run`, persist
artifacts and an execution trace, and replay it to an identical result. No LLM,
no network, no scheduler — just the smallest honest runtime.

This document is the design contract for Phase 1. It refines the roadmap entry
and pins the shapes before the Go code lands.

## Public interfaces (contracts)

Machine-readable schemas live in `contracts/`:

- [`workflow.schema.json`](../../contracts/workflow.schema.json) — the input: a
  DAG of steps, each with a `capability`, its `needs`, and `with` inputs.
- [`task-envelope.schema.json`](../../contracts/task-envelope.schema.json) — one
  self-contained unit of work handed to a provider.
- [`artifact-manifest.schema.json`](../../contracts/artifact-manifest.schema.json)
  — an immutable, content-addressed output with provenance.
- [`execution-trace.schema.json`](../../contracts/execution-trace.schema.json) —
  the durable record of a run.

## Deterministic providers (v0.1)

Pure functions only. Same inputs → same bytes, always:

| Capability | Input | Output |
|---|---|---|
| `text.constant` | `value` (string) | that string |
| `text.uppercase` | `text` (string) | uppercased string |
| `text.concat` | `parts` (array of strings) | joined string |

References: a `with` value of `{ "from_task": "<step-id>" }` is replaced by that
upstream step's output artifact content before the envelope is dispatched.

## Run directory layout

Each run is written under the project's `.rexo/` tree:

```
.rexo/
  artifacts/
    <content-hash>            # immutable bytes, addressed by sha256
  runs/
    <run-id>/
      trace.json             # execution-trace.schema.json
      events.jsonl           # append-only event log
      artifacts/
        <task-id>.json       # artifact-manifest per task output
```

Artifacts are content-addressed and never overwritten. The manifest ties a
task's output to its `inputs_fingerprint` (sha256 of capability + resolved
inputs), which is what makes replay verifiable and future caching possible.

## Execution model

1. Load and validate the workflow. Topologically sort the steps (fail on cycles
   or missing `needs`).
2. Generate a `run_id`. Write an initial `trace.json` with status `running`.
3. For each step in order: resolve `with` (inject upstream outputs), build a
   Task Envelope, run the deterministic provider, store the output as a
   content-addressed artifact + manifest, append events, update the trace.
4. On completion, set the trace status to `succeeded`. Any provider error sets
   the failing task to `failed`, the run to `failed`, and stops — the partial
   trace stays on disk and is diagnosable.

## Replay

`rexo run <workflow> --replay <run-id>` re-executes from the recorded envelopes
and asserts each task reproduces the same `content_hash`. A mismatch is a hard
error (non-determinism detected). This is the Phase 1 acceptance test.

## CLI surface added

```
rexo run <workflow.json> [--project <dir>]     # execute a workflow
rexo run <workflow.json> --replay <run-id>     # verify determinism
```

## Acceptance (from the roadmap)

- Execute and replay the two-step `examples/hello.workflow.json`.
- An interruption produces a diagnosable on-disk state (partial trace + events).
- Artifacts include provenance (`produced_by` + `inputs_fingerprint`).

## Explicitly not now

LLM calls, provider selection, distributed queue, semantic memory, retries with
backoff (Phase 2), caching (Phase 4).
