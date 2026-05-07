#!/usr/bin/env bash
# Validate a .gitlab-ci.yml file against GitLab's CI Lint via glab.
# Usage: ./validate.sh [path-to-file]
# Default: .gitlab-ci.yml in the current directory.
#
# Always runs `glab ci lint --repo <PROJECT>` so it works from any directory,
# including cloned GitHub repos with no GitLab remote.
#
# Exit codes:
#   0 - file is syntactically valid
#   1 - file is invalid (errors printed)
#   2 - file not found, or glab missing

set -uo pipefail

FILE="${1:-.gitlab-ci.yml}"

# Repository to lint against. Any GitLab project the authenticated user can read works;
# this is just the "context repo" glab needs. Override with VALIDATE_PROJECT env var.
PROJECT="${VALIDATE_PROJECT:-gitlab-org/ci-cd/github-actions-to-gitlab-ci}"

if [ ! -f "$FILE" ]; then
  echo "ERROR: $FILE not found" >&2
  exit 2
fi

if ! command -v glab >/dev/null 2>&1; then
  echo "ERROR: glab is required. Install from https://gitlab.com/gitlab-org/cli" >&2
  exit 2
fi

echo "Validating $FILE via glab ci lint (repo: $PROJECT)..."
if glab ci lint --repo "$PROJECT" "$FILE"; then
  exit 0
else
  exit 1
fi
