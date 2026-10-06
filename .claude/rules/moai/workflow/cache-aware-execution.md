# Cache-Aware Execution

Prompt caching is a prefix match over the rendered request: a cache read costs a fraction of a write, and the cache expires after about five minutes idle. These points decide when and in what order you act; they never change what an approval gate requires.

- Ask the questions you know you will need early, while the context is small. A long wait on a question late in a big session can let the cache expire over the whole prefix; when several late questions are unavoidable, ask them together.
- When fanning out several agents of the same type, start one first and the rest once it is producing output. Requests that start together cannot read a cache entry still being written, so all of them pay the full write.
- Edits to files loaded at session start (`CLAUDE.md`, `AGENTS.md`, `.claude/rules/`, output styles) invalidate the cached prefix for every later turn. Batch them at the end of a task.
- Before a large multi-agent batch in a context bloated with finished, unrelated work, a `/clear` plus the resume message (`session-handoff.md`) is cheaper than carrying the bloat through every spawn. When only short follow-up work remains, keep the warm cache.
- Caches are per model, and a mid-session model or effort switch discards them; spawn agents without a model override and switch at a natural boundary.
- Pass a file's content with an `@`-mention or a Read rather than naming it for the model to fetch.
