# Repository Guidelines

## Project Structure & Module Organization

Devopin is a Go service with an embedded Vue 3/Vuetify frontend. Backend entry points and Cobra commands live in `cmd/`; application internals are organized under `internal/` by layer, including HTTP handlers, services, repositories, models, workers, monitoring, database access, and log parsing. Reusable public helpers belong in `pkg/`. The frontend is in `frontend/`, with pages, components, stores, assets, and Vite configuration under `frontend/src/`. Product and implementation notes are in `docs/`, and `configs/config.yaml.example` documents local configuration. The build copies frontend output into `internal/web/dist/` before compiling the Go binary.

## Build, Test, and Development Commands

- `CGO_ENABLED=0 go test ./...` runs all Go package tests without a C toolchain.
- `go vet ./...` performs basic Go correctness checks.
- `make build_frontend` installs frontend dependencies with pnpm, builds the SPA, and copies its assets into the Go project.
- `make build_core` builds the current-platform binary at `build/godevopin`.
- `make build_all` cleans embedded assets and builds the complete application.
- `cd frontend && pnpm dev` starts the Vite development server; `pnpm type-check` runs Vue/TypeScript checks and `pnpm lint` runs ESLint (with its configured fixes).
- `./build/godevopin serve --port 9000` runs the built server. Local development requires PostgreSQL with TimescaleDB enabled.

## Coding Style & Naming Conventions

Format Go changes with `gofmt` and follow idiomatic Go naming: exported identifiers use `PascalCase`, locals and unexported functions use `camelCase`, and packages use lowercase names. Keep new backend code in the appropriate `internal/` layer. Frontend components use `PascalCase.vue`; pages, stores, and utilities follow the existing lowercase or kebab-case filenames. Use TypeScript types rather than untyped values, and run the frontend type check and ESLint before submitting UI changes.

## Testing Guidelines

Go tests should be colocated with the package they cover and named `*_test.go`; run them with `go test ./...`. There is currently no frontend test framework configured, so validate frontend changes with `pnpm type-check`, `pnpm lint`, `pnpm build`, and focused manual checks. No repository-wide coverage threshold is currently defined.

## Commit & Pull Request Guidelines

Use concise, imperative commit subjects in the established style, such as `feat: add ...`, `fix(frontend): ...`, or `refactor: ...`. Pull requests should explain the behavior change, list validation commands, and call out database or configuration changes. Include screenshots or short recordings for UI changes and link a related issue when one exists. Never commit credentials, JWT secrets, database passwords, or local `config.yaml` files.
