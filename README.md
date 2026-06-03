# agent-cluster-backend

Go backend for the agent_cluster_project platform.

**Owned by this repo:** Go implementation, GraphQL API, SSE event surface, local
runtime, generated-contract consumption, backend-side authorization and
validation adapters.

**Not owned by this repo:** any domain vocabulary, lifecycle, event, policy,
view, or other DSL concept — those are defined in
[`agent-cluster-contracts`](https://github.com/kimjooyoon/agent-cluster-contracts)
and consumed here via generated code (see `contracts/dsl/`, `contracts/ir/`).

This repo must not redefine anything that lives in the contracts repo. The
cross-repo vocab-scan guard (a future contracts-owned tool) will enforce this.

## Status

Bootstrap-minimum only. No GraphQL or SSE code yet; the first vertical slice
arrives with the next decision record after `001-initial-agreement` and
`002-dumb-agent-role`.

## What works today

- `go build ./...` produces a single binary that prints a placeholder line
  identifying the build and the contracts repo it points at.
- `.github/workflows/security.yml` runs secretscan on every push/PR (the
  binary is sourced from this repo; the contracts repo's secretscan is the
  canonical implementation).

## Layout (future)

```
.
├── cmd/
│   └── agent-cluster-backend/    server entrypoint
├── internal/
│   ├── graphql/                  consumes ../contracts/ir/* (generated)
│   ├── sse/
│   └── ...
├── contracts/                    git submodule or vendored generated client
└── .github/workflows/
```

## How to run (local-first)

```sh
go build -o bin/ ./cmd/...
./bin/agent-cluster-backend
```
