#!/usr/bin/env bash
set -euo pipefail

echo "Running OpenAPI validation for specs/openapi.yaml"
npx @redocly/cli lint specs/openapi.yaml
echo "OpenAPI validation completed successfully"
