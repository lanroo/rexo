# REXO Core v1 Roadmap

The architecture is frozen. Each phase must leave a usable, testable system and
must not implement later-phase concerns prematurely.

## Phase 0 — Public Foundation

**Objective:** create a trustworthy, portable project that people can inspect,
build, and install.

**Components:** constitution, ADRs, schemas, repository governance, CLI
`version/doctor/init`, cross-platform CI, release packaging.

**Public interfaces:** CLI surface and Project Manifest v0.1.

**Acceptance:** tests on Windows/macOS/Linux; five release targets; checksums;
no mandatory runtime dependency; `init` never overwrites an existing path.

**Dependencies:** none.

**Risks:** brand collision, platform-specific paths, promising features that do
not exist.

**Tests:** unit, vet, smoke initialization, schema parsing, release build.

**Not now:** LLM calls, scheduler, MCP, marketplace, Studio.

## Phase 1 — Walking Skeleton

**Objective:** execute one deterministic workflow end to end.

**Components:** Kernel API, run ID, task envelope, artifact store, local event
log, deterministic provider, minimal CLI run command.

**Public interfaces:** Task Envelope, Artifact Manifest, Execution Trace.

**Acceptance:** execute and replay a two-step workflow; interruption produces a
diagnosable state; artifacts include provenance.

**Dependencies:** Phase 0 contracts.

**Risks:** leaking implementation details into public schemas.

**Tests:** golden workflow, crash recovery boundary, idempotency.

**Not now:** distributed queue, semantic memory, LLM provider selection.

## Phase 2 — Durable Runtime

**Objective:** resume workflows safely after interruption.

**Components:** state machine, checkpoints, retry policy, timeout, cancellation,
compensation hooks.

**Public interfaces:** Workflow Definition v0.1 and run lifecycle API.

**Acceptance:** resume without duplicating completed side effects; bounded
retries; deterministic terminal states.

**Dependencies:** Phase 1 event and artifact model.

**Risks:** ambiguous idempotency and hidden side effects.

**Tests:** kill/restart, retry exhaustion, cancellation, migration fixture.

**Not now:** general distributed consensus.

## Phase 3 — Capability Registry and Tool Gateway

**Objective:** decouple workflow intent from providers.

**Components:** capability catalog, provider catalog, compatibility checks,
resolver, Gateway interface, one deterministic and one MCP adapter.

**Public interfaces:** Capability and Provider Manifests.

**Acceptance:** swap providers without changing a workflow; deny providers that
violate policy; explain each resolution.

**Dependencies:** durable task execution.

**Risks:** capability names becoming vendor-shaped.

**Tests:** resolution matrix, contract conformance, unavailable-provider
fallback.

**Not now:** public marketplace ranking.

## Phase 4 — Economy Engine

**Objective:** minimize cost while preserving a declared quality floor.

**Components:** Budget Manager, cache lookup, artifact fingerprints,
model/provider tiers, incremental build planner, usage ledger.

**Public interfaces:** Budget Policy and Economy Decision.

**Acceptance:** hard budget enforcement; cache hit avoids provider call;
decisions are explainable; partial changes rebuild only affected nodes.

**Dependencies:** resolver, provenance, artifacts.

**Risks:** false cache hits and quality degradation.

**Tests:** budget exhaustion, cache invalidation, model downgrade, incremental
rebuild.

**Not now:** speculative global optimization.

## Phase 5 — Layered Memory

**Objective:** retrieve only the smallest relevant memory layer.

**Components:** global/project/execution/task stores, retrieval policy,
promotion workflow, compression, retention.

**Public interfaces:** Memory Query and Memory Record.

**Acceptance:** task workers cannot enumerate unauthorized layers; promotion is
audited; large documents use structured summaries when sufficient.

**Dependencies:** policy, artifacts, Economy Engine.

**Risks:** sensitive-data leakage and stale summaries.

**Tests:** access boundaries, provenance, expiry, summary invalidation.

**Not now:** autonomous global knowledge merging.

## Phase 6 — Quality Engine

**Objective:** enforce measurable acceptance criteria.

**Components:** evaluator registry, rubrics, deterministic validators,
LLM-as-judge adapter, retry/fallback controller.

**Public interfaces:** Evaluation Report and Quality Policy.

**Acceptance:** failed outputs cannot silently pass; evaluator versions are
recorded; loops always terminate.

**Dependencies:** durable runtime and Economy Engine.

**Risks:** biased evaluators and self-grading.

**Tests:** threshold boundaries, conflicting evaluators, max attempts.

**Not now:** universal quality score.

## Phase 7 — Education MVP

**Objective:** validate Core with a commercial vertical.

**Flow:** briefing → research reuse → curriculum → lesson plan → script → PDF →
QA → course package.

**Acceptance:** one complete course package; unchanged research is reused;
editing one lesson rebuilds only affected outputs; budget and QA reports ship
with the package.

**Dependencies:** Phases 1–6.

**Risks:** education-specific logic leaking into Core.

**Tests:** reference course, incremental edit, provider failure, pedagogical
rubrics.

**Not now:** video generation, LMS analytics, automatic publishing.

## Phase 8 — Packs and SDK

**Objective:** let third parties create extensions without Core changes.

**Components:** Pack Manifest, pack lifecycle, scaffolding, validators,
TypeScript/Python SDK priorities based on user evidence.

**Acceptance:** create, validate, install, update, and remove a sample pack;
compatibility failures are safe and clear.

**Dependencies:** stable Core contracts proven by Education.

**Risks:** supply-chain attacks and premature SDK breadth.

**Tests:** malicious archive, migration, rollback, conformance suite.

**Not now:** unmoderated marketplace.

## Phase 9 — Distributions and Private Registry

**Objective:** compose governed REXO variants.

**Components:** distribution manifest, lockfile, signed registry metadata,
private catalog, reproducible installer.

**Acceptance:** rebuild an identical distribution from lockfile and verified
artifacts.

**Dependencies:** pack lifecycle and signing.

**Risks:** dependency resolution complexity.

**Not now:** public commercial marketplace.

## Phase 10 — Studio

**Objective:** visually author workflows over stable contracts.

**Components:** graph editor, schema forms, validation, run inspection, artifact
canvas.

**Acceptance:** Studio-created workflows execute identically through CLI/SDK.

**Dependencies:** mature workflow schemas and observed user needs.

**Risks:** UI defining semantics that do not exist in Core.

**Not now:** Creator-generated ecosystems.

## Phase 11 — Creator

**Objective:** generate a proposed vertical platform from a domain brief.

**Components:** pack/distribution scaffolding, policy and knowledge proposals,
test generation, human approval gates.

**Acceptance:** generated outputs pass the same conformance tests as human-built
packs and never self-publish.

**Dependencies:** proven SDK, registry, quality, and Studio contracts.

**Risks:** unsafe policy generation and false confidence.

**Not now:** autonomous self-deployment.
