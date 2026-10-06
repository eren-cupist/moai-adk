# Session Handoff — Appendix (examples)

> Reference appendix: the two illustrative examples. This location sits outside the loaded-instruction budget, so it costs no trigger budget; the rules file keeps a pointer.

## Example (Illustrative; substitute project-specific values when adapting)

```
✂──── 여기부터 복사 ────✂

ultrathink. SPEC-EXAMPLE-A implementation 진입.
applied lessons: <lesson-id-1>, <lesson-id-2>.
source_session_id: <not-available — environment-fallback, next session will backfill via /moai session register on activation>

전제 검증:
1) git log --oneline -1 → <commit-sha> 확인
2) ls .moai/specs/SPEC-EXAMPLE-A/ → N files

실행: /moai run SPEC-EXAMPLE-A

머지 후: SPEC-EXAMPLE-B → SPEC-EXAMPLE-C

✂──── 여기까지 복사 ────✂
```

> Block 5 carries the work-starting action.


## Example with Block 0 (Illustrative)

```
✂──── 여기부터 복사 ────✂

[New Terminal — START IN WORKTREE]
$ moai cc -w ~/.moai/worktrees/<project>/SPEC-EXAMPLE-A
   # (launcher -w accepts L2 absolute paths; or moai glm -w ...)

ultrathink. SPEC-EXAMPLE-A Epic N 진입.
applied lessons: <lesson-id-1>, <lesson-id-2>.

전제 검증:
0) git rev-parse --show-toplevel → ~/.moai/worktrees/<project>/SPEC-EXAMPLE-A (★ critical)
1) gh pr view <PR-number> → MERGED

실행: /moai run SPEC-EXAMPLE-A --team

후속: Milestone M<N+1> (single-SPEC next step) 또는 Epic N+1 (multi-SPEC next grouping)

✂──── 여기까지 복사 ────✂
```

---

