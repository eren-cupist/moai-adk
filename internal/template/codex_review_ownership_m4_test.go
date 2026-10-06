// codex_review_ownership_m4_test.go pins the M4 slice of the review-ownership
// work: which agents hold the on-demand self-review tools, where the lane's
// card-review stage sits in the distributed kanban doctrine, and that the
// distributed workflow template mentions the tree-scope policy key only as a
// commented example.
//
// Every checker takes document text and returns the list of problems it
// found, so the same function judges the real files (expecting none) and a
// table of deliberately broken variants (expecting a named problem). A check
// that has never been seen to fail on a known bad input proves nothing
// (verification-completeness §1.1); the mutant tables are that observation.
package template_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	coCodexReview = "mcp__moai__codex_review"
	coGLMReview   = "mcp__moai__glm_review"
	coCodexAudit  = "mcp__moai__codex_audit"
	coGLMAudit    = "mcp__moai__glm_audit"
)

var (
	coReviewHolders = []string{"manager-develop", "manager-docs"}
	coAuditHolders  = []string{"plan-auditor", "sync-auditor"}
)

// --- agent tool lists ---

// coReadAgents returns agent name -> file content for every *.md in dir.
func coReadAgents(t *testing.T, dir string) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		out[strings.TrimSuffix(e.Name(), ".md")] = string(raw)
	}
	if len(out) == 0 {
		t.Fatalf("no agent definitions found in %s — empty sweep", dir)
	}
	return out
}

// coToolSet parses the `tools:` line of an agent definition's frontmatter.
func coToolSet(md string) map[string]bool {
	set := map[string]bool{}
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "tools:") {
			for _, tok := range strings.Split(strings.TrimPrefix(line, "tools:"), ",") {
				if tok = strings.TrimSpace(tok); tok != "" {
					set[tok] = true
				}
			}
			break
		}
	}
	return set
}

func coHolders(agents map[string]string, tool string) []string {
	var holders []string
	for name, md := range agents {
		if coToolSet(md)[tool] {
			holders = append(holders, name)
		}
	}
	sort.Strings(holders)
	return holders
}

// coHolderProblems checks the two holder sets: the review tools sit with the
// three lane-side agents only, the audit tools with the two auditors only.
func coHolderProblems(agents map[string]string) []string {
	var problems []string
	check := func(tool string, want []string) {
		got := coHolders(agents, tool)
		if strings.Join(got, ",") != strings.Join(want, ",") {
			problems = append(problems, fmt.Sprintf("%s holders = %v, want exactly %v", tool, got, want))
		}
	}
	check(coCodexReview, coReviewHolders)
	check(coGLMReview, coReviewHolders)
	check(coCodexAudit, coAuditHolders)
	check(coGLMAudit, coAuditHolders)
	return problems
}

// coParityProblems requires the local and distributed copies of each agent to
// declare the same tool set.
func coParityProblems(c1, c2 map[string]string) []string {
	var problems []string
	names := map[string]bool{}
	for n := range c1 {
		names[n] = true
	}
	for n := range c2 {
		names[n] = true
	}
	for n := range names {
		a, aok := c1[n]
		b, bok := c2[n]
		if !aok || !bok {
			problems = append(problems, fmt.Sprintf("%s exists in only one of the local and distributed copies", n))
			continue
		}
		sa, sb := coToolSet(a), coToolSet(b)
		for tool := range sa {
			if !sb[tool] {
				problems = append(problems, fmt.Sprintf("%s: %s is listed locally but not in the distributed copy", n, tool))
			}
		}
		for tool := range sb {
			if !sa[tool] {
				problems = append(problems, fmt.Sprintf("%s: %s is listed in the distributed copy but not locally", n, tool))
			}
		}
	}
	sort.Strings(problems)
	return problems
}

// coWithTool returns agents with the tool appended to / removed from one
// agent's `tools:` line. ok is false when the edit changed nothing.
func coWithTool(agents map[string]string, agent, tool string, add bool) (map[string]string, bool) {
	out := map[string]string{}
	for k, v := range agents {
		out[k] = v
	}
	md, found := out[agent]
	if !found {
		return out, false
	}
	lines := strings.Split(md, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, "tools:") {
			continue
		}
		switch {
		case add && !coToolSet(md)[tool]:
			lines[i] = line + ", " + tool
		case !add && coToolSet(md)[tool]:
			lines[i] = strings.Replace(line, ", "+tool, "", 1)
		default:
			return out, false
		}
		out[agent] = strings.Join(lines, "\n")
		return out, out[agent] != md
	}
	return out, false
}

// --- doctrine ---

func coRead(t *testing.T, root, rel string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(raw)
}

const (
	coTplPrefix = "internal/template/templates/"
)

// --- template workflow.yaml ---

// coWorkflowYAMLProblems judges the distributed workflow template: the
// codex.review_gate block parses to the single live key `enabled`, and the
// tree-scope key appears only inside comment lines (at least once).
func coWorkflowYAMLProblems(text string) []string {
	var problems []string
	var doc struct {
		Workflow struct {
			Codex struct {
				ReviewGate map[string]any `yaml:"review_gate"`
			} `yaml:"codex"`
		} `yaml:"workflow"`
	}
	if err := yaml.Unmarshal([]byte(text), &doc); err != nil {
		return []string{"workflow.yaml does not parse: " + err.Error()}
	}
	keys := make([]string, 0, len(doc.Workflow.Codex.ReviewGate))
	for k := range doc.Workflow.Codex.ReviewGate {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if strings.Join(keys, ",") != "enabled" {
		problems = append(problems, fmt.Sprintf("workflow.codex.review_gate keys = %v, want exactly [enabled]", keys))
	}
	if v, ok := doc.Workflow.Codex.ReviewGate["enabled"]; !ok || v != false {
		problems = append(problems, fmt.Sprintf("workflow.codex.review_gate.enabled = %v, want false", v))
	}
	commented := 0
	for i, line := range strings.Split(text, "\n") {
		if !strings.Contains(line, "tree_scope") {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			commented++
			continue
		}
		problems = append(problems, fmt.Sprintf("line %d carries tree_scope outside a comment: %q", i+1, line))
	}
	if commented == 0 {
		problems = append(problems, "no commented example of tree_scope in the template")
	}
	return problems
}

// --- tests ---

// TestReviewOwnership_ToolHolders — AC-012 (a)(b): the self-review tools sit
// with manager-develop, manager-docs and manager-lead only; the audit tools
// stay with plan-auditor and sync-auditor only; the local and distributed
// copies declare the same tools.
func TestReviewOwnership_ToolHolders(t *testing.T) {
	root := findProjectRootForMirrorTest(t)
	c1 := coReadAgents(t, filepath.Join(root, ".claude/agents/moai"))
	c2 := coReadAgents(t, filepath.Join(root, coTplPrefix+".claude/agents/moai"))
	for label, agents := range map[string]map[string]string{"local (.claude/agents/moai)": c1, "distributed (templates)": c2} {
		for _, p := range coHolderProblems(agents) {
			t.Errorf("%s: %s", label, p)
		}
	}
	for _, p := range coParityProblems(c1, c2) {
		t.Errorf("local/distributed tools drift: %s", p)
	}
}

// TestReviewOwnership_ToolHoldersCodexEmission — AC-012 (c), the emitted layer:
// the audit tools are named in the two auditors' emitted definitions only, and
// no emitted definition outside the three lane-side agents names a review tool.
// The byte-level regeneration check is the agentemit golden test; this pins who
// is allowed to mention what.
func TestReviewOwnership_ToolHoldersCodexEmission(t *testing.T) {
	root := findProjectRootForMirrorTest(t)
	dir := filepath.Join(root, coTplPrefix+".codex/agents/moai")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	holders := map[string][]string{}
	n := 0
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		n++
		body := coRead(t, root, coTplPrefix+".codex/agents/moai/"+e.Name())
		name := strings.TrimSuffix(e.Name(), ".toml")
		for _, tool := range []string{coCodexAudit, coGLMAudit, coCodexReview, coGLMReview} {
			if strings.Contains(body, tool) {
				holders[tool] = append(holders[tool], name)
			}
		}
	}
	if n == 0 {
		t.Fatal("no emitted agent definitions found — empty sweep")
	}
	for _, tool := range []string{coCodexAudit, coGLMAudit} {
		if got := strings.Join(holders[tool], ","); got != strings.Join(coAuditHolders, ",") {
			t.Errorf("emitted definitions naming %s = %v, want exactly %v", tool, holders[tool], coAuditHolders)
		}
	}
	allowed := map[string]bool{}
	for _, h := range coReviewHolders {
		allowed[h] = true
	}
	for _, tool := range []string{coCodexReview, coGLMReview} {
		for _, h := range holders[tool] {
			if !allowed[h] {
				t.Errorf("emitted definition %s names %s but is not a lane-side holder", h, tool)
			}
		}
	}
}

// TestReviewOwnership_ToolHolderMutants — the holder checker must go red on
// each cheapest way of getting the grant wrong.
func TestReviewOwnership_ToolHolderMutants(t *testing.T) {
	root := findProjectRootForMirrorTest(t)
	base := coReadAgents(t, filepath.Join(root, coTplPrefix+".claude/agents/moai"))
	fourth := "manager-git"
	if _, ok := base[fourth]; !ok {
		t.Fatalf("fixture agent %s is not in the roster", fourth)
	}
	cases := []struct {
		name  string
		agent string
		tool  string
		add   bool
		want  string // substring the problem list must contain
	}{
		{"review tool granted to a fourth agent", fourth, coCodexReview, true, fourth},
		{"glm review granted to an auditor", "sync-auditor", coGLMReview, true, "sync-auditor"},
		{"audit tool granted to a lane agent", "manager-develop", coGLMAudit, true, "manager-develop"},
		{"review tool withheld from a lane agent", "manager-develop", coGLMReview, false, "manager-develop"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// A removal mutant needs the grant to exist; make the base carry it first.
			start := base
			if !c.add {
				full := base
				for _, h := range coReviewHolders {
					for _, tool := range []string{coCodexReview, coGLMReview} {
						full, _ = coWithTool(full, h, tool, true)
					}
				}
				start = full
			}
			mut, changed := coWithTool(start, c.agent, c.tool, c.add)
			if !changed {
				t.Fatalf("mutant did not apply: %s %v %s", c.agent, c.add, c.tool)
			}
			problems := strings.Join(coHolderProblems(mut), "\n")
			if !strings.Contains(problems, c.want) {
				t.Errorf("the holder checker did not flag %q; problems:\n%s", c.want, problems)
			}
		})
	}
	t.Run("distributed copy edited, local copy not", func(t *testing.T) {
		c1 := coReadAgents(t, filepath.Join(root, ".claude/agents/moai"))
		// Flip the grant in the distributed copy only: add it when absent, remove it when present.
		c2, changed := coWithTool(base, "manager-docs", coCodexReview, !coToolSet(base["manager-docs"])[coCodexReview])
		if !changed {
			t.Fatal("mutant did not apply")
		}
		if len(coParityProblems(c1, c2)) == 0 {
			t.Error("a one-sided edit passed the local/distributed parity check")
		}
	})
	t.Run("parity checker flags either direction", func(t *testing.T) {
		a := map[string]string{"x": "tools: Read, mcp__moai__codex_review"}
		b := map[string]string{"x": "tools: Read"}
		if len(coParityProblems(a, b)) == 0 || len(coParityProblems(b, a)) == 0 {
			t.Error("parity checker is blind to a one-sided tools edit")
		}
	})
}

// TestReviewOwnership_TemplateWorkflowKeyIsCommentOnly — AC-015 (iii)(iv): the
// distributed workflow template mentions tree_scope only in a comment, and the
// surfaces that enumerate live keys do not know it.
func TestReviewOwnership_TemplateWorkflowKeyIsCommentOnly(t *testing.T) {
	root := findProjectRootForMirrorTest(t)
	text := coRead(t, root, coTplPrefix+".moai/config/sections/workflow.yaml")
	for _, p := range coWorkflowYAMLProblems(text) {
		t.Error(p)
	}

	inv := coRead(t, root, "internal/config/testdata/shipped_key_inventory.yaml")
	// Positive control: the inventory is read and still enumerates the
	// neighbouring review_gate keys, so the absence check below sweeps a
	// non-empty set.
	if !strings.Contains(inv, "review_gate") {
		t.Errorf("shipped_key_inventory.yaml carries no review_gate key — the absence check below would sweep nothing")
	}
	for _, rel := range []string{
		"internal/config/testdata/shipped_key_inventory.yaml",
		"internal/settings/schema_sections.go",
		"internal/web/fieldsets.templ",
		"internal/web/assets/i18n.js",
	} {
		if strings.Contains(coRead(t, root, rel), "tree_scope") {
			t.Errorf("%s mentions tree_scope — the key ships as a template comment only", rel)
		}
	}
}

// TestReviewOwnership_TemplateWorkflowMutants — the template checker must go
// red on each way of shipping the key live.
func TestReviewOwnership_TemplateWorkflowMutants(t *testing.T) {
	root := findProjectRootForMirrorTest(t)
	text := coRead(t, root, coTplPrefix+".moai/config/sections/workflow.yaml")
	anchor := "        review_gate:\n            enabled: false\n"
	if !strings.Contains(text, anchor) {
		t.Fatal("template review_gate block not found at the expected shape")
	}
	cases := []struct {
		name string
		edit func(string) string
		want string
	}{
		{"live key added", func(s string) string {
			return strings.Replace(s, anchor, anchor+"            tree_scope: review\n", 1)
		}, "keys ="},
		{"live key beside a comment example", func(s string) string {
			return strings.Replace(s, anchor, anchor+"            # tree_scope: review\n            tree_scope: skip\n", 1)
		}, "outside a comment"},
		{"example removed", func(s string) string {
			var keep []string
			for _, l := range strings.Split(s, "\n") {
				if !strings.Contains(l, "tree_scope") {
					keep = append(keep, l)
				}
			}
			return strings.Join(keep, "\n")
		}, "no commented example"},
		{"enabled flipped on", func(s string) string {
			return strings.Replace(s, anchor, "        review_gate:\n            enabled: true\n", 1)
		}, "enabled"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mut := c.edit(text)
			if mut == text && c.name != "example removed" {
				t.Fatal("mutant did not apply")
			}
			got := strings.Join(coWorkflowYAMLProblems(mut), "\n")
			if !strings.Contains(got, c.want) {
				t.Errorf("not flagged (%q); problems:\n%s", c.want, got)
			}
		})
	}
}
