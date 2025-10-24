# TypeScript SDK

https://github.com/kwila-cloud/simple-sync/issues/18

Build an offline-first TypeScript SDK for the Simple-Sync API and add an OpenAPI spec so we can validate and CI-check SDKs.

### docs: add OpenAPI spec and spec validation tests
- [ ] Add `specs/openapi.yaml` describing the public API endpoints used by SDK (`/events`, `/auth/setup-token/exchange`, `/acl` endpoints used by clients).
- [ ] Add contract tests that validate the OpenAPI spec is loadable and passes linting (e.g., `swagger-cli validate`) under `tests/contract/openapi_spec_test.go`.
- [ ] Add a CI job step to run OpenAPI lint/validate on push and pull requests.

### feat: add generated TypeScript client scaffold
- [ ] Create directory `clients/typescript` with basic TypeScript project
- [ ] Add prettier config
- [ ] Add eslint config
- [ ] Run prettier and eslint in CI/CD

### feat: implement full TypeScript SDK
- [ ] Generate typescript SDK from `specs/openapi.yaml`

### feat: add SDK validation contract tests against test server
- [ ] Add contract tests under `tests/contract/` that start the test server (the repo already has test helpers) and verify the TypeScript SDK behavior against real endpoints (authentication, event post/get, ACL checks).

- [ ] Ensure tests fail if the OpenAPI spec and server disagree.
- [ ] Wire these contract tests into CI (separate job that starts the test server and runs Node tests).

### docs: TypeScript SDK
- [ ] Add docs to `docs/` describing the TypeScript SDK example, how to run generation, and how to run the contract tests locally.
- [ ] Add a GitHub Actions workflow (`.github/workflows/sdk-validation.yml`) or extend existing CI to include a job that installs Node, generates the client, builds the example, and runs the contract tests.

