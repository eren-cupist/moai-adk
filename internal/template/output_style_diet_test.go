// Output-style size guard. An output style is loaded on every turn of every session that selects
// it, so its size is a standing cost; each budget may only go down.
package template

import "testing"

// Whole-file UTF-16 budgets of the deployed output styles.
const (
	dietBudgetMoai     = 2500
	dietBudgetMoaiEasy = 2500
)

var dietBudgets = []struct {
	file   string
	name   string
	budget int
}{
	{"moai.md", "moai", dietBudgetMoai},
	{"moai-easy.md", "moai-easy", dietBudgetMoaiEasy},
}

func TestOutputStylesCharBudget(t *testing.T) {
	deployed := dietReadDeployed(t)
	for _, b := range dietBudgets {
		t.Run(b.name, func(t *testing.T) {
			size := dietUTF16Len(deployed[b.file])
			t.Logf("output-style=%s %d", b.name, size)
			if size > b.budget {
				t.Errorf("output style %s is %d UTF-16 units, over its budget %d", b.file, size, b.budget)
			}
			if _, _, err := dietSplitFrontmatter(deployed[b.file]); err != nil {
				t.Errorf("output style %s: %v", b.file, err)
			}
		})
	}
}
