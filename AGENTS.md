# AGENTS.md — Warehouse Management System (WMS)

## What this repo is

- Warehouse Management System with: (backend) Go + (frontend) Vue3.
- Key domains: Inventory, Locations (aisle/bay/shelf/bin), Tasks, Products and Users.

## Repo map (high level)

- cmd/: Go entry point
- internal/: Go backend APIs, services, providers and utils
- frontend/web: Vue3 web app
- scripts/: dev scripts, migrations, tooling

## Style

- Go: gofmt, golangci-lint clean, no new lint warnings.
- Vue: follow existing component patterns; avoid breaking routes.

## Safety

- Do not run destructive commands (drop DB, reset prod configs) without asking.
