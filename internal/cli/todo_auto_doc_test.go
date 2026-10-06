// todo_auto_doc_test.go — the `--auto` flag help and its refusal text do not
// assert a pick order.
package cli

import (
	"strings"
	"testing"
)

// autoDocNormalize collapses every whitespace run to one space and drops the
// backticks, so a phrase compares equal across a reflow and across the
// markdown code spans a flag help string cannot carry.
func autoDocNormalize(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "`", "")), " ")
}

// TestAutoRankHelpAndRefusalDoNotAssertPickOrder (REQ-TAP-013 wording): the
// `--auto` flag help and the card-argument refusal no longer say the cycle
// consumes "the queue in queue order" — the cycle ranks its queued candidates,
// so neither string may assert the pick order. The test name sits outside the
// TestAutoRank prefix so the planned sweep count of that selector is unchanged.
func TestAutoHelpAndRefusalDoNotAssertPickOrder(t *testing.T) {
	help := newTodoCmd().Flags().Lookup("auto")
	if help == nil {
		t.Fatal("the todo command has no --auto flag")
	}
	if strings.Contains(autoDocNormalize(help.Usage), "queue order") {
		t.Errorf("the --auto flag help still asserts the pick order: %q", help.Usage)
	}

	todoFixture(t)
	_, _, err := runTodo(t, "--auto", "please", "accept")
	if err == nil {
		t.Fatal("--auto with card arguments must be refused")
	}
	msg := err.Error()
	if !strings.Contains(msg, "--auto takes no card arguments") {
		t.Fatalf("unexpected refusal text %q", msg)
	}
	if strings.Contains(autoDocNormalize(msg), "queue order") {
		t.Errorf("the refusal still asserts the pick order: %q", msg)
	}
}
