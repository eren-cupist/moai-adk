# Quality Reference

Lookup detail behind the TRUST 5 checklist in `SKILL.md`: where the thresholds live, how evaluator profiles are chosen, and the toolchain for languages beyond the main four.

## Configuration

| Setting | File | Meaning |
|---------|------|---------|
| `test_coverage_target` | `.moai/config/sections/quality.yaml` | module-level coverage target (default 85) |
| `development_mode` | `quality.yaml` | the run-phase cycle (`tdd` or `ddd`) |
| `lsp_quality_gates.<phase>` | `quality.yaml` | LSP thresholds per phase; sync: `max_errors: 0`, `max_warnings: 10`, `require_clean_lsp: true` |
| `harness.default_profile` | `.moai/config/sections/harness.yaml` | evaluator profile when the SPEC names none |
| `levels.<level>.evaluator` | `harness.yaml` | whether `sync-auditor` runs at that harness level (`minimal`: false) |
| `evaluator_profile` | SPEC frontmatter | profile for this SPEC: `.moai/config/evaluator-profiles/<name>.md` |

The shipped profiles are `default`, `strict`, `lenient`, and `frontend`. Each names its pass thresholds, must-pass areas, and hard thresholds; the default makes Functionality (every acceptance criterion met) and Security (no critical or high finding) must-pass, and fails Craft below 85% coverage. A project may name critical modules that need 90% or more.

A coverage figure counts only when the coverage command ran on the current tree and its output was observed — not a number carried over from an earlier run.

## More toolchains

Use whatever the project has configured; the common tools are:

| Language | Lint | Format | Test |
|----------|------|--------|------|
| Rust | `cargo clippy` | `rustfmt` | `cargo test`, `cargo llvm-cov` |
| Java | checkstyle / spotbugs | google-java-format | JUnit |
| Kotlin | detekt | ktlint | JUnit |
| C# | analyzers | `dotnet format` | xUnit / NUnit |
| Ruby | rubocop | rubocop | rspec |
| PHP | phpstan, phpcs | php-cs-fixer | pest, phpunit |
| Elixir | credo | `mix format` | `mix test` |
| C++ | clang-tidy | clang-format | gtest / catch2 |
| Scala | scalafix | scalafmt | scalatest |
| Swift | swiftlint | swift-format | XCTest |
| Dart/Flutter | `dart analyze` | `dart format` | `flutter test` |
| R | lintr | styler | testthat |

## Related

- `moai-ref-owasp-checklist` — the Secured items in detail
- `moai-ref-testing-pyramid` — test strategy and coverage targets
- `.claude/rules/moai/core/verification-claim-integrity.md` — claims of verification must be observed
