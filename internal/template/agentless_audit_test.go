// agentless_audit_test.go: Audit suite for the Agentless contract of the
// utility workflow (fix.md) and the pipeline-rejection sentinel of the
// implementation workflows.
//
// @MX:NOTE - Two tests: TestAgentlessUtilityNoLLMControlFlow (REQ-WF004-013),
// TestImplementationSkillsContainPipelineRejectionSentinel (REQ-WF004-014).
// The run.md mode-dispatch sentinels (MODE_UNKNOWN, MODE_TEAM_UNAVAILABLE) and
// the utility MODE_FLAG_IGNORED_FOR_UTILITY sentinel were retired with the
// --mode flag; no non-test Go code matched them.
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

// implementationSkillPaths lists the 3 implementation skill files that must reject
// the --mode pipeline flag per REQ-WF004-014.
// (design.md removed by SPEC-SUBCOMMAND-RETIRE-001, 2026-07-01.)
var implementationSkillPaths = []string{
	".claude/skills/moai/workflows/plan.md",
	".claude/skills/moai/workflows/sync.md",
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

// TestImplementationSkillsContainPipelineRejectionSentinel verifies that each of the
// 4 implementation skills contains the literal sentinel string MODE_PIPELINE_ONLY_UTILITY
// (REQ-WF004-014). At M1 (RED), all 4 subtests fail because the sentinel has not yet
// been added to the skill files.
func TestImplementationSkillsContainPipelineRejectionSentinel(t *testing.T) {
	t.Parallel()

	fsys, err := EmbeddedTemplates()
	if err != nil {
		t.Fatalf("EmbeddedTemplates() error: %v", err)
	}

	const sentinel = "MODE_PIPELINE_ONLY_UTILITY"

	for _, skillPath := range implementationSkillPaths {
		t.Run(path.Base(skillPath), func(t *testing.T) {
			t.Parallel()

			data, readErr := fs.ReadFile(fsys, skillPath)
			if readErr != nil {
				t.Fatalf("ReadFile(%q) error: %v", skillPath, readErr)
			}

			if !strings.Contains(string(data), sentinel) {
				t.Errorf("file %s missing sentinel %s", skillPath, sentinel)
			}
		})
	}
}
