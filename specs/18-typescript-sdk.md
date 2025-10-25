# TypeScript SDK

https://github.com/kwila-cloud/simple-sync/issues/18

Build an offline-first TypeScript SDK for the Simple-Sync API and add an OpenAPI spec so we can validate and CI-check SDKs.

### docs: add OpenAPI spec and spec validation tests
- [x] Add `specs/openapi.yaml` describing the full public API surface, as described in `docs/src/content/docs/api/v1.md`
- [ ] Add contract tests (written in bash) that validate the OpenAPI spec is loadable and passes linting (e.g., `swagger-cli validate`).
- [ ] Add a CI job step to run OpenAPI lint/validate on push and pull requests.
- [ ] Add contract tests (written in bash) that validate that a locally running instance of the full API matches the specification in `specs/openapi.yaml`.
- [ ] Add a CI job step to run the API validation on push and pull requests.

### feat: generate TypeScript SDK
- [ ] Generate typescript SDK from `specs/openapi.yaml`, using [openapi-ts](https://github.com/hey-api/openapi-ts)
- [ ] Add prettier config to the typescript SDK
- [ ] Add eslint config to the typescript SDK
- [ ] Run prettier and eslint in CI/CD

### feat: add SDK validation contract tests against test server
- [ ] Add contract tests (in bash) that start the local server and verify the TypeScript SDK behavior against real endpoints (authentication, event post/get, ACL checks).
- [ ] Ensure tests fail if the TypeScript SDK and server disagree.
- [ ] Wire these contract tests into CI (separate job that starts the test server and runs Node tests).

### docs: TypeScript SDK
- [ ] Add docs to `docs/src/content/docs/sdk/typescript.md` describing the TypeScript SDK

