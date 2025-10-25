#!/usr/bin/env bash
set -euo pipefail

echo "Running OpenAPI validation for specs/openapi.yaml"

# Prefer using npx to avoid global installs. In CI this will fetch the package.
if command -v npx >/dev/null 2>&1; then
  npx @apidevtools/swagger-cli validate specs/openapi.yaml
else
  echo "npx not found; attempting to use docker if available"
  if command -v docker >/dev/null 2>&1; then
    docker run --rm -v "$(pwd)":/work -w /work node:18 bash -c "npm install --no-audit @apidevtools/swagger-cli && npx @apidevtools/swagger-cli validate specs/openapi.yaml"
  else
    echo "Neither npx nor docker available; skipping OpenAPI validation"
    exit 0
  fi
fi

echo "OpenAPI validation completed successfully"
