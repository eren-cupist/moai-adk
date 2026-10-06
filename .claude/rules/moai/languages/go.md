---
paths: "**/*.go,**/go.mod,**/go.sum"
---

# Go

Commands:

- Format: `gofmt -w` (or `goimports -w` if the project uses it) on changed files.
- Vet and lint: `go vet ./...`; `golangci-lint run` when the repo has a `.golangci.*` config.
- Test: `go test ./<changed packages>/...` while iterating, then `go test ./...`. Add `-race` for concurrent code and `-count=1` to bypass the test cache when verifying. Coverage: `go test -cover ./<pkg>/...` or `-coverprofile`.
- `moai gate` runs `go vet ./...`, then golangci-lint when configured, then `go test ./...`.

Conventions:

- Follow the module's existing package layout and error style. Wrap errors with `%w` where context is added and test them with `errors.Is` or `errors.As`, not string matching.
- Tests are table-driven with `t.Run` subtests where cases vary only by input; use `t.TempDir()` and `t.Setenv()` for isolation (a test that calls `t.Setenv` cannot also call `t.Parallel`).
- Run `go mod tidy` after adding or removing imports and commit `go.mod` and `go.sum` together.
