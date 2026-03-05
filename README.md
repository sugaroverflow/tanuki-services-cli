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

## Local setup

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

## CI / Pipeline

### GitHub Actions (`.github/workflows/`)

- **CI** (`ci.yml`) — on every push: validate registry, test (matrix: Linux/macOS/Windows), build binaries, CodeQL SAST.
- **Publish** (`publish.yml`) — on push to `main`: build catalog, deploy staging, smoke-test, deploy production (with approval).
- **Nightly** (`nightly-validate.yml`) — scheduled: validate registry health.

---

## License

MIT (or your choice) — demo use.
