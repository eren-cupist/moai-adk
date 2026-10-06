// agentless_audit_test.go: Audit suite for the Agentless contract of the
// utility workflow (fix.md).
//
// @MX:NOTE - TestAgentlessUtilityNoLLMControlFlow (REQ-WF004-013). The --mode
// flag and its MODE_* sentinels were retired from the workflows; no non-test Go
// code matched them, so no sentinel test remains.
package template

import (
	"io/fs"
	"path"
	"regexp"
	"strings"
	"testing"
)

// utilitySkillPaths lists the utility skill files subject to Agentless classification.
// Path separator is forward-slash (embedded FS convention).
// (coverage.md removed by SPEC-SUBCOMMAND-RETIRE-001, 2026-07-01.)
var utilitySkillPaths = []string{
	".claude/skills/moai/workflows/fix.md",
}

// forbiddenControlFlowPatterns are regex patterns whose presence in utility skill bodies
// (outside code blocks) indicates LLM-driven control flow — a violation of the Agentless
// contract (REQ-WF004-013). See research.md §6.2.
//
// @MX:ANCHOR fan_in=4 - SPEC-V3R2-WF-004 REQ-WF004-013 enforcer; guards 4 utility
// skills against LLM-dispatch regression. Touching this regex set affects the contract
// for fix/mx/codemaps/clean. (coverage retired by SPEC-SUBCOMMAND-RETIRE-001.)
var forbiddenControlFlowPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Use the .* subagent to (decide|determine|choose|select|orchestrate|route|dispatch)`),
	regexp.MustCompile(`(?i)Use the .* subagent to (plan|design) the (pipeline|workflow|next phase|sequence)`),
	regexp.MustCompile(`(?i)delegate to .* (orchestrator|router|dispatcher|controller)`),
	regexp.MustCompile(`(?i)manager-strategy.*subagent.*(branch|fork|conditional)`),
}

// TestAgentlessUtilityNoLLMControlFlow verifies that none of the 5 utility workflow
// skills contain LLM-driven control-flow patterns (REQ-WF004-013).
//
// This test is a regression guard: at M1 all 5 subtests pass because no utility skill
// currently violates. The test will turn red if a future PR introduces one of the
// forbidden patterns into a utility skill body.
//
// @MX:ANCHOR fan_in=5 - SPEC-V3R2-WF-004 REQ-WF004-013 enforcer.
func TestAgentlessUtilityNoLLMControlFlow(t *testing.T) {
	t.Parallel()

	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates() error: %v", err)
	}

	for _, skillPath := range utilitySkillPaths {
		t.Run(path.Base(skillPath), func(t *testing.T) {
			t.Parallel()

			data, readErr := fs.ReadFile(fsys, skillPath)
			if readErr != nil {
				t.Fatalf("ReadFile(%q) error: %v", skillPath, readErr)
			}

			lines := strings.Split(string(data), "\n")
			inCodeBlock := false
			for lineIdx, line := range lines {
				// Toggle code block state on fence open/close.
				if strings.HasPrefix(strings.TrimSpace(line), "```") {
					inCodeBlock = !inCodeBlock
					continue
				}
				if inCodeBlock {
					continue
				}

				// Check each forbidden pattern against the non-code-block line.
				// Per acceptance.md AC-WF004-12 Failure Scenario, the error message
				// MUST contain the literal sentinel "AGENTLESS_CONTROL_FLOW_VIOLATION"
				// so CI log parsers (grep) can detect regressions.
				for patIdx, re := range forbiddenControlFlowPatterns {
					if match := re.FindString(line); match != "" {
						t.Errorf(
							"AGENTLESS_CONTROL_FLOW_VIOLATION: %s line %d matches forbidden pattern #%d %q (matched: %q)",
							skillPath, lineIdx+1, patIdx, re.String(), match,
						)
					}
				}
			}
		})
	}
}
