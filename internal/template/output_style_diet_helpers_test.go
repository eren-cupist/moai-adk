// Helpers shared by the output-style budget guard and the agent-description budget guard.
package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf16"
)

// dietStyleFiles lists the deployed output styles in a fixed order.
var dietStyleFiles = []string{"moai.md", "moai-easy.md"}

// dietStyleTemplateDir is the template source directory, relative to the package directory.
const dietStyleTemplateDir = "templates/.claude/output-styles/moai"

// dietUTF16Len returns the length of s in UTF-16 code units (the JavaScript string length).
func dietUTF16Len(s string) int {
	return len(utf16.Encode([]rune(s)))
}

// dietSplitFrontmatter splits a file into its frontmatter (through the closing --- line and
// its newline) and its body (everything after).
func dietSplitFrontmatter(text string) (head, body string, err error) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", "", fmt.Errorf("no frontmatter")
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			return strings.Join(lines[:i+1], "\n") + "\n", strings.Join(lines[i+1:], "\n"), nil
		}
	}
	return "", "", fmt.Errorf("frontmatter not closed")
}

// dietReadDeployed returns the whole text of each deployed (template) output-style file.
func dietReadDeployed(t *testing.T) map[string]string {
	t.Helper()
	out := make(map[string]string, len(dietStyleFiles))
	for _, f := range dietStyleFiles {
		raw, err := os.ReadFile(filepath.Join(dietStyleTemplateDir, f))
		if err != nil {
			t.Fatalf("read deployed %s: %v", f, err)
		}
		out[f] = string(raw)
	}
	return out
}
