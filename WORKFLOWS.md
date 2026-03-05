# Workflows — step by step

What each workflow does, job by job and step by step.

---

## Shared piece: `setup-tanuki-env`

Used by **CI** (test + build) and **Publish** (build-binary, smoke-test). It’s a composite action that:

1. **Set up Go** — Installs the Go version you pass in (default 1.21), uses `go.sum` for caching.
2. **Cache Go modules** — Caches `~/go/pkg/mod` and `~/.cache/go-build` so later runs are faster.

So whenever a job says “Setup Tanuki env,” it’s: checkout + this (Go + cache), then that job’s own steps.

---

## 1. CI workflow

**Runs when:** Any push to any branch, or any pull request.

**Order:** Jobs run in sequence where there’s a dependency; `sast` runs in parallel with the rest.

### Job 1: `validate`

- **Runner:** Ubuntu.
- **Steps:**
  1. Checkout the repo.
  2. Set up Python 3.11.
  3. `pip install -r requirements.txt`.
  4. `python scripts/build_catalog.py --validate` — checks every `registry/*.yml` against `catalog.schema.json`. Fails the job if any file is invalid.
- **Purpose:** Ensure the registry YAMLs are valid before we run tests or build.

### Job 2: `test`

- **Runs after:** `validate` (`needs: validate`).
- **Runs on:** Three runners in parallel — Ubuntu, macOS, Windows (matrix).
- **Steps:**
  1. Checkout the repo.
  2. Run the **setup-tanuki-env** action (Go + cache).
  3. `go test -v ./...` — run all Go tests on that OS.
- **Purpose:** Confirm the code passes tests on all supported platforms.

### Job 3: `build`

- **Runs after:** `test` (`needs: test`).
- **Runs on:** Same matrix — Ubuntu, macOS, Windows (three jobs).
- **Steps:**
  1. Checkout the repo.
  2. Run **setup-tanuki-env**.
  3. Build the binary with `GOOS`/`GOARCH` set per runner (e.g. `linux/amd64`, `darwin/amd64`, `windows/amd64`), output e.g. `tanuki-linux-amd64` or `tanuki-windows-amd64.exe`.
  4. **Upload artifact** — the binary is saved as a workflow artifact (name like `tanuki-linux-amd64`).
- **Purpose:** Produce the `tanuki` CLI for each OS so you can download them from the Actions run.

### Job 4: `sast`

- **Runs:** In parallel with the above (no `needs`).
- **Runner:** Ubuntu.
- **Steps:**
  1. Checkout the repo.
  2. **Initialize CodeQL** for Go.
  3. **CodeQL analyze** — runs the security analysis and uploads results (e.g. to the Security tab).
- **Purpose:** Security scan on the Go code; no dependency on validate/test/build.

**Summary:** CI = validate registry → test on 3 OSes → build binary on 3 OSes (and upload) + SAST in parallel.

---

## 2. Publish workflow

**Runs when:** A push to the `main` branch (merge or direct push).

**Order:** Strict sequence: build-catalog → build-binary and deploy-staging (both need build-catalog; deploy-staging also needs build-binary) → smoke-test → deploy-production.

### Job 1: `build-catalog`

- **Runner:** Ubuntu.
- **Steps:**
  1. Checkout the repo.
  2. Set up Python 3.11, install `requirements.txt`.
  3. `python scripts/build_catalog.py -o dist/catalog.json` — builds the full catalog from `registry/*.yml` and writes `dist/catalog.json`.
  4. **Upload artifact** named `catalog` containing `dist/catalog.json`.
- **Purpose:** Produce the single catalog JSON that the CLI and “released” catalog use.

### Job 2: `build-binary`

- **Runs after:** `build-catalog`.
- **Runner:** Ubuntu.
- **Steps:**
  1. Checkout the repo.
  2. **Setup Tanuki env** (Go + cache).
  3. `go build -o tanuki-linux-amd64 ./cmd/tanuki` — build the Linux CLI.
  4. **Upload artifact** named `tanuki-binary` with that binary.
- **Purpose:** Produce the Linux binary for “release” (staging/production in the demo).

### Job 3: `deploy-staging`

- **Runs after:** Both `build-catalog` and `build-binary`.
- **Runner:** Ubuntu; uses GitHub environment **staging**.
- **Steps:**
  1. **Download artifact** `catalog` (the `catalog.json`).
  2. **Download artifact** `tanuki-binary` (the Linux binary).
  3. A demo step that echoes that staging is ready and moves the catalog file if needed (in a real setup you’d push these to a staging bucket or Pages).
- **Purpose:** Represent “deploy to staging”; in this repo it only prepares the artifacts.

### Job 4: `smoke-test`

- **Runs after:** `deploy-staging`.
- **Runner:** Ubuntu.
- **Steps:**
  1. Checkout the repo.
  2. **Download artifact** `catalog` so the built catalog is in the workspace.
  3. **Setup Tanuki env**, then `go build -o tanuki ./cmd/tanuki` to build the CLI in this job.
  4. Run `./tanuki list` and then `./tanuki list | grep -q .` — the CLI reads the downloaded catalog; the pipeline fails if `tanuki list` fails or returns nothing.
- **Purpose:** Check that the catalog we built is valid and the CLI can list services from it.

### Job 5: `deploy-production`

- **Runs after:** `smoke-test`.
- **Runner:** Ubuntu; uses GitHub environment **production** (typically configured for manual approval in the repo settings).
- **Steps:**
  1. A single step that echoes “Production deploy complete” (in a real setup you’d upload catalog + binary to production Pages or a release).
- **Purpose:** Represent “deploy to production” after staging and smoke-test; in this repo it’s a demo step and the gate is the manual approval on the production environment.

**Summary:** Publish = build catalog → build Linux binary → “deploy” to staging → smoke-test with `tanuki list` → “deploy” to production (with approval).

---

## 3. Nightly validate workflow

**Runs when:** On a schedule at 02:00 UTC every day, or when you click “Run workflow” in the Actions tab.

### Job: `validate-registry`

- **Runner:** Ubuntu.
- **Steps:**
  1. Checkout the repo.
  2. Set up Python 3.11, install `requirements.txt`.
  3. `python scripts/build_catalog.py -o dist/catalog.json` — build the catalog (same as Publish).
  4. Run a small script that loads the JSON and prints how many services are in the catalog (and a message that in production you could add health-url checks here).
- **Purpose:** Daily check that the registry still builds; you can extend it later to hit each service’s `health_url`.

---

## 4. Dependabot

This isn’t a workflow you edit; GitHub runs it.

- **When:** Weekly, and when Dependabot opens a PR.
- **What:** Dependabot opens pull requests to update:
  - Go modules (`go.mod` / `go.sum`),
  - GitHub Actions versions in your workflow files,
  - pip dependencies (`requirements.txt`).
- **Where:** You’ll see “Dependabot Updates” in the Actions list and PRs from `dependabot` or `dependabot[bot]`.

---

## Quick reference

| Workflow   | Trigger        | Flow |
|-----------|----------------|------|
| **CI**    | Push / PR      | validate → test (3 OSes) → build (3 OSes) ∥ sast |
| **Publish** | Push to `main` | build-catalog → build-binary → deploy-staging → smoke-test → deploy-production |
| **Nightly** | 02:00 UTC / manual | validate-registry (build catalog + print service count) |
| **Dependabot** | Weekly / PRs | Opens dependency update PRs |
