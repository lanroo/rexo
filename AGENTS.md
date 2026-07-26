# REXO repository instructions

## Mandatory bootstrap

Before changing this repository:

1. Read `REXO_BOOTSTRAP.md`.
2. Read `REXO_STATE.md`.
3. Read only the ADRs and contracts relevant to the task.
4. Inspect the current working tree and preserve unrelated user changes.

## Context economy

- Never send the whole repository to an agent.
- Provide only the task objective, required artifacts, relevant memory, budget,
  and success criteria.
- Prefer deterministic tools and cached artifacts before any LLM call.
- Workers are disposable. Persist artifacts and provenance, not conversations.
- Summarize large inputs before delegation when the summary preserves the facts
  needed for the task.

## Engineering rules

- Core is capability-first.
- Keep public contracts language-neutral.
- Do not add a dependency without an ADR-level reason.
- Do not make Docker, Python, Node, or a cloud account mandatory for the CLI.
- All filesystem behavior must work on Windows, macOS, and Linux.
- Never commit secrets, runtime memory, caches, build output, or user data.
- Run `go test ./...` and `go vet ./...` before proposing a commit.
- Architecture changes require a measurable, recurring problem and an ADR.
