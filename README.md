# Tanuki Services CLI

CLI to query a service catalog. Manifests are YAML in this repo; CI validates, builds, and publishes the catalog and CLI. Catalog is source of truth and only updates via CI.

## Commands

```bash
tanuki list                    # All services
tanuki status <name>           # Health, version, owner, last deploy
tanuki owners <name>           # Owner and on-call info
tanuki search --team <team>    # Filter by team
tanuki validate                # Validate registry against schema
```

## Setup & run

**Build catalog (one-time):**
```bash
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python scripts/build_catalog.py -o catalog.json
```

**Build and run CLI:**
```bash
go build -o tanuki ./cmd/tanuki
./tanuki list
```

Catalog is read from `TANUKI_CATALOG_URL`, or `./catalog.json`, or `./dist/catalog.json`.

**Run CI:** Push to GitHub; see [Actions](https://github.com/sugaroverflow/tanuki-services-cli/actions). CI runs on every push; Publish runs on push to `main`.

## Adding a service

1. Add `registry/<name>.yml` following `catalog.schema.json` (required: name, version, owner, team, health_url, repo_url).
2. Run `tanuki validate`.
3. Merge to main; CI builds and publishes the catalog.

## GitHub Actions

| Workflow | Trigger | Jobs |
|----------|---------|------|
| **CI** | Push / PR | validate → test (matrix) → build (matrix) → SAST |
| **Publish** | Push to `main` | build-catalog → build-binary → deploy-staging → smoke-test → deploy-production |
| **Nightly validate** | Daily 02:00 UTC | Build catalog, validate registry |
| **Dependabot** | Weekly | Dependency update PRs (Go, Actions, pip) |

Workflows: [`.github/workflows/`](.github/workflows/) · Runs: [Actions](https://github.com/sugaroverflow/tanuki-services-cli/actions) · Step-by-step: [WORKFLOWS.md](WORKFLOWS.md)

## License

MIT
