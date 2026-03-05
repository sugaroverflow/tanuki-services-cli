# Tanuki Services CLI

A CLI that gives developers a single command to query an internal **service catalog**. Service manifests live as YAML files in this repo; the CI pipeline validates, builds, and publishes both the catalog and the CLI binaries. The catalog is the source of truth and only updates through CI, never manually.

---

## Commands

```bash
tanuki list                      # All registered services
tanuki status <name>              # Health, version, owner, last deploy
tanuki owners <name>             # Owner and on-call info
tanuki search --team <team>      # Filter by team
tanuki validate                  # Validate local registry against schema
```

---

## Run the application locally

From the repo root:

```bash
# 1. Create Python venv and build the catalog (one-time)
python3 -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install -r requirements.txt
python scripts/build_catalog.py -o catalog.json

# 2. Build the CLI
go build -o tanuki ./cmd/tanuki

# 3. Run commands
./tanuki list
./tanuki status payments-api
./tanuki owners auth-service
./tanuki search --team platform
./tanuki validate
```

If you already have `catalog.json` or `dist/catalog.json` in the repo, you can skip step 1 and just run `go build -o tanuki ./cmd/tanuki` then `./tanuki list`.

---

## Test with GitHub Actions

1. Create a repo on GitHub and add it as `origin` (if you haven’t already):
   ```bash
   git remote add origin https://github.com/YOUR_ORG/tanuki-services-cli.git
   ```
2. Push a branch:
   ```bash
   git add .
   git commit -m "Add Tanuki CLI and registry"
   git push -u origin main
   ```
3. Open the repo on GitHub → **Actions**. You should see:
   - **CI** running on every push (validate → test → build → SAST).
   - **Publish** only when you push to `main` (build-catalog → deploy-staging → smoke-test → deploy-production).

---

## Local setup (detailed)

### 1. Build the catalog (optional, for local use)

```bash
python3 -m venv .venv
source .venv/bin/activate   # or .venv\Scripts\activate on Windows
pip install -r requirements.txt
python scripts/build_catalog.py -o catalog.json
```

Or use the default output path:

```bash
python scripts/build_catalog.py   # writes dist/catalog.json
```

### 2. Build and run the CLI

```bash
go build -o tanuki ./cmd/tanuki
./tanuki list
./tanuki status payments-api
./tanuki search --team platform
./tanuki validate   # uses .venv/bin/python3 if present
```

**Catalog source:** The CLI looks for the catalog in this order:

- Env var `TANUKI_CATALOG_URL` (e.g. your published catalog URL)
- `./catalog.json`
- `./dist/catalog.json`

---

## Adding a service

1. Add a new YAML file under `registry/` (e.g. `registry/my-service.yml`).
2. Follow the schema in `catalog.schema.json` (required: `name`, `version`, `owner`, `team`, `health_url`, `repo_url`).
3. Run `tanuki validate` or `python scripts/build_catalog.py --validate` to check.
4. Merge to main; CI builds the catalog and publishes it.

---

## GitHub Actions

Workflows live in [`.github/workflows/`](.github/workflows/). **Workflow runs:** [sugaroverflow/tanuki-services-cli → Actions](https://github.com/sugaroverflow/tanuki-services-cli/actions)

| Workflow | When it runs | What it does |
|----------|----------------|---------------|
| **CI** (`ci.yml`) | Every push and every pull request (any branch) | **validate** — Check all `registry/*.yml` against the schema. **test** — Run `go test ./...` on Linux, macOS, and Windows. **build** — Build the `tanuki` binary per OS and upload artifacts. **sast** — CodeQL security analysis on the Go code. |
| **Publish** (`publish.yml`) | Only when you push to `main` | **build-catalog** — Build `catalog.json` from the registry. **build-binary** — Build the Linux CLI. **deploy-staging** — Demo staging (prepares artifacts). **smoke-test** — Run `tanuki list` to confirm catalog works. **deploy-production** — After manual approval of the production environment (demo). |
| **Nightly validate** (`nightly-validate.yml`) | Daily at 02:00 UTC (or “Run workflow”) | Build catalog and validate registry; placeholder for health-url checks. |
| **Dependabot Updates** | Weekly + when Dependabot opens PRs | PRs to bump Go modules, GitHub Actions, and pip deps (managed by GitHub). |

**In short:** CI = validate + test + build + SAST on every change. Publish = build and release catalog + CLI on merge to main. Nightly = registry health check. Dependabot = dependency update PRs.

---

## License

MIT
