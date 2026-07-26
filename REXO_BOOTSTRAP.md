# REXO bootstrap

This file is the stable entry point for humans and AI coding agents.

## Read order

1. `REXO_STATE.md`
2. `docs/architecture/constitution.md`
3. `docs/roadmap/core-v1.md`
4. Relevant files under `docs/adr/`
5. Only the contracts and source files required by the current task

## Activation contract

An activation request means:

- restore the current project state from versioned files;
- identify the active roadmap phase and its acceptance criteria;
- load the smallest useful context;
- reuse existing research, templates, prompts, and artifacts;
- state the execution budget when LLM calls or paid providers are involved;
- implement incrementally and verify on all affected platforms;
- update `REXO_STATE.md` when a milestone materially changes.

It does not authorize publishing, spending money, accessing secrets, or changing
external systems beyond the explicit request.

## Sources of truth

- Markdown and machine-readable contracts are authoritative.
- PDFs are human-oriented snapshots.
- Git history is the change record.
- Runtime memory and caches are never authoritative architecture.
