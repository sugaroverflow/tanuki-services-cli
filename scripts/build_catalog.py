#!/usr/bin/env python3
"""
Build catalog.json from registry/*.yml manifests.
Validates each manifest against catalog.schema.json.
Run from repo root: python scripts/build_catalog.py [--validate] [--output PATH]
"""
import argparse
import json
import sys
from pathlib import Path

import yaml
from jsonschema import Draft7Validator, ValidationError

REPO_ROOT = Path(__file__).resolve().parent.parent
REGISTRY_DIR = REPO_ROOT / "registry"
SCHEMA_PATH = REPO_ROOT / "catalog.schema.json"
DEFAULT_OUTPUT = REPO_ROOT / "dist" / "catalog.json"


def load_schema():
    with open(SCHEMA_PATH, encoding="utf-8") as f:
        return json.load(f)


def load_and_validate_manifests(schema):
    validator = Draft7Validator(schema)
    catalog = []
    yaml_files = sorted(REGISTRY_DIR.glob("*.yml")) + sorted(REGISTRY_DIR.glob("*.yaml"))
    if not yaml_files:
        print("No .yml/.yaml files in registry/", file=sys.stderr)
        sys.exit(1)

    for path in yaml_files:
        with open(path, encoding="utf-8") as f:
            data = yaml.safe_load(f)
        if data is None:
            print(f"Empty or invalid YAML: {path}", file=sys.stderr)
            sys.exit(1)
        errors = list(validator.iter_errors(data))
        if errors:
            for err in errors:
                print(f"{path}: {err.message}", file=sys.stderr)
            sys.exit(1)
        catalog.append(data)
    return catalog


def main():
    parser = argparse.ArgumentParser(description="Build catalog from registry YAMLs")
    parser.add_argument("--validate", action="store_true", help="Only validate; do not write catalog")
    parser.add_argument("--output", "-o", type=Path, default=DEFAULT_OUTPUT, help="Output path for catalog.json")
    args = parser.parse_args()

    if not REGISTRY_DIR.is_dir():
        print(f"Registry directory not found: {REGISTRY_DIR}", file=sys.stderr)
        sys.exit(1)
    if not SCHEMA_PATH.is_file():
        print(f"Schema not found: {SCHEMA_PATH}", file=sys.stderr)
        sys.exit(1)

    schema = load_schema()
    catalog = load_and_validate_manifests(schema)

    if args.validate:
        print("Validation passed.")
        return

    args.output.parent.mkdir(parents=True, exist_ok=True)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(catalog, f, indent=2)
    print(f"Wrote {len(catalog)} services to {args.output}")


if __name__ == "__main__":
    main()
