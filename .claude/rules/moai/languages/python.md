---
paths: "**/*.py,**/pyproject.toml,**/requirements*.txt"
---

# Python

Commands:

- Run tools through the project's environment manager: `uv run <tool>` when `uv.lock` exists, `poetry run <tool>` with `poetry.lock`, otherwise the active virtualenv. Add dependencies with that manager (`uv add`, `poetry add`), not bare `pip install`.
- Lint and format: `ruff check` and `ruff format` (or the formatter the project configures).
- Type check: `mypy` or `pyright`, whichever pyproject.toml configures.
- Test: `pytest` scoped to the changed modules while iterating, then the full suite; coverage with `pytest --cov=<package>`.
- `moai gate` runs `ruff check .` and `mypy .` (each only when its config or pyproject.toml exists and a .py file is staged), then `pytest`.

Conventions:

- Configuration lives in pyproject.toml; follow its ruff, mypy and pytest settings rather than adding per-tool config files.
- Public functions carry type hints; match the project's existing async or sync style and do not call blocking I/O from async code.
