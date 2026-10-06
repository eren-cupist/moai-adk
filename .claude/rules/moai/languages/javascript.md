---
paths: "**/*.js,**/*.mjs,**/*.cjs,**/package.json"
---

# JavaScript

The package-manager, workspace and turbo commands in [typescript.md](typescript.md) apply to JavaScript projects as well.

- Module format follows package.json `"type"` and the file extension (`.mjs` ESM, `.cjs` CommonJS); do not mix `require` and `import` in one module.
- Lint with whichever linter the project configures. `moai gate` runs eslint, biome or oxlint only when that tool's config file is present, and runs `npm test -- --passWithNoTests` as its test step.
- When editing package.json, keep dependency changes in the package that uses them and update the lockfile with the project's package manager in the same change.
