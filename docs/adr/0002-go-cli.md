# ADR 0002: Implement the portable CLI in Go

- Status: accepted
- Date: 2026-07-26

## Context

The first public artifact must run on Windows, macOS, and Linux without asking
users to install a language runtime or Docker.

## Decision

Implement the CLI and initial local control-plane components in Go. Avoid CGO in
Core v1. Keep public contracts language-neutral with JSON Schema.

## Consequences

Releases can contain standalone binaries for amd64 and arm64. Python and
TypeScript remain valid future SDK and provider languages, but they are not
installation requirements for the base CLI.
