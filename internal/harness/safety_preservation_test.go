// Package harness — 5-Layer Safety Architecture preservation tests (T-C5).
//
// This file holds the architectural assertion tests for SPEC-V3R4-HARNESS-002 Wave C.
// The two tests verify architectural invariants rather than code:
//  1. The frozenPrefixes slice must contain exactly 4 canonical entries.
//
// REQ-HRN-FND-006: guarantee immutability of the FROZEN path protection list.
// REQ-HRN-OBS-002: guarantee immutability of the 5-Layer Safety Architecture.
package harness

import (
	"testing"
)

// TestSafetyArchitecture_FrozenZoneUnchanged verifies the frozenPrefixes slice
// contains exactly the 4 canonical entries.
// REQ-HRN-FND-006: immutability of the FROZEN path protection list.
//
// Note: tasks.md T-C5 lists `.moai/project/brand/` as the 4th entry, but the
// current frozen_guard.go implementation includes `.claude/skills/moai/`
// instead. This test verifies actual code state; the discrepancy with
// tasks.md is recorded as an SPEC-V3R4-HARNESS-002 implementation note.
// Future changes must update this test together with the code.
func TestSafetyArchitecture_FrozenZoneUnchanged(t *testing.T) {
	t.Parallel()

	// Actual entries currently present in frozen_guard.go (order included)
	wantPrefixes := []string{
		".claude/agents/moai/",
		".claude/skills/moai-",
		".claude/skills/moai/",
		".claude/rules/moai/",
	}

	// Verify entry count
	if len(frozenPrefixes) != len(wantPrefixes) {
		t.Errorf("frozenPrefixes entry count: got=%d, want=%d", len(frozenPrefixes), len(wantPrefixes))
		t.Logf("actual entries: %v", frozenPrefixes)
		return
	}

	// Verify each entry order and value
	for i, want := range wantPrefixes {
		if frozenPrefixes[i] != want {
			t.Errorf("frozenPrefixes[%d]: got=%q, want=%q", i, frozenPrefixes[i], want)
		}
	}

	// Integration check: FROZEN paths must actually be blocked.
	// The constitution.md path must be blocked by the .claude/rules/moai/ prefix.
	constitutionPath := ".claude/rules/moai/design/constitution.md"
	_, err := IsAllowedPath(constitutionPath)
	if err == nil {
		t.Errorf("IsAllowedPath(%q) must return FrozenViolationError", constitutionPath)
	} else {
		var frozenErr *FrozenViolationError
		if !isFrozenViolationError(err, &frozenErr) {
			t.Errorf("IsAllowedPath(%q) error type: got=%T, want=*FrozenViolationError", constitutionPath, err)
		}
	}
}

// isFrozenViolationError checks whether the error is of type *FrozenViolationError.
func isFrozenViolationError(err error, out **FrozenViolationError) bool {
	if fve, ok := err.(*FrozenViolationError); ok {
		if out != nil {
			*out = fve
		}
		return true
	}
	return false
}
