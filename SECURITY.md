# Security policy

## Reporting

Do not open a public issue for a suspected vulnerability. Until a dedicated
security address is published, use GitHub's private vulnerability reporting for
this repository.

## Supported versions

Before the first stable release, only the latest published version receives
security fixes.

## Repository hygiene

Never commit API keys, access tokens, `.env` files, private memory, runtime
transcripts, user data, or generated caches. REXO manifests reference secrets
by logical name; they do not contain secret values.

Packs, providers, model output, and MCP servers must be treated as untrusted
inputs. The current foundation release does not install or execute third-party
packs.
