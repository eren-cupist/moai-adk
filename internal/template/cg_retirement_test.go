package template

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCGEmbeddedRetirementPreservesRoutingAndAudit(t *testing.T) {
	embedded, err := EmbeddedTemplates()
	if err != nil {
		t.Fatal(err)
	}
	renderer := NewRenderer(embedded)
	render := func(path string) string {
		t.Helper()
		body, err := renderer.Render(path, nil)
		if err != nil {
			t.Fatalf("actual embedded render %s: %v", path, err)
		}
		return string(body)
	}
	claude := render("CLAUDE.md")
	if strings.Contains(claude, "60-70% cost reduction") {
		t.Error("CG cost guarantee survived")
	}
	var cfg struct {
		Harness struct {
			Modes map[string]string `yaml:"mode_defaults"`
			Audit struct {
				Always bool `yaml:"always_enabled"`
			} `yaml:"plan_audit_global"`
		} `yaml:"harness"`
	}
	if err := yaml.Unmarshal([]byte(render(".moai/config/sections/harness.yaml")), &cfg); err != nil {
		t.Fatal(err)
	}
	if _, exists := cfg.Harness.Modes["cg"]; exists {
		t.Error("CG still selects harness depth")
	}
	if cfg.Harness.Modes["solo"] != "auto" || cfg.Harness.Modes["team"] != "auto" || !cfg.Harness.Audit.Always {
		t.Error("ordinary mode or mandatory plan audit lost")
	}
	for _, path := range []string{".claude/skills/moai/workflows/run/task-decomposition.md", ".claude/skills/moai/workflows/run/phase-execution.md", ".claude/agents/moai/sync-auditor.md", ".codex/agents/moai/sync-auditor.toml"} {
		body := render(path)
		if !strings.Contains(body, "sync-auditor") {
			t.Errorf("independent auditor missing %s", path)
		}
		if strings.Contains(body, "CG mode:") || strings.Contains(body, "**CG mode**:") {
			t.Errorf("CG inline audit exception survives %s", path)
		}
	}
}
