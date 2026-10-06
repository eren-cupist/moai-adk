// red_now_cell_test.go — repository-local verifier for the RED-now cell
// adoption gate (card t343).
//
// It checks the rule clause in verification-completeness.md (both mirrors)
// and runs the RED-now command form check over the fixtures.
//
// Nothing here ships. The test is repository-local and its fixtures live under
// internal/spec/testdata/red_now, outside the distributed template tree.
package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

const (
	redNowRulePath       = ".claude/rules/moai/development/verification-completeness.md"
	redNowRuleMirrorPath = "internal/template/templates/.claude/rules/moai/development/verification-completeness.md"

	redNowFixtureDir = "testdata/red_now"

	redNowRuleSectionHead = "## 2. Two-cell adoption discipline"
)

// redNowForbidden is the machine-checkable half of the RED-now command form:
// a cited command carrying one of these tokens unquoted is more than one shell
// invocation. It used to be extracted from a sentinel span in the
// plan-auditor agent; that agent now states the rule in prose ("a single
// read-only invocation"), so the token list lives here.
var redNowForbidden = []string{"|", "&&", ";", ">", "<", "$(", "("}

// ---------------------------------------------------------------------------
// span extraction
// ---------------------------------------------------------------------------

func redNowRead(t *testing.T, rel string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(raw)
}

// redNowHeadingLevel returns the markdown heading level of a line, or 0.
func redNowHeadingLevel(line string) int {
	n := 0
	for n < len(line) && line[n] == '#' {
		n++
	}
	if n == 0 || n >= len(line) || line[n] != ' ' {
		return 0
	}
	return n
}

// redNowExtractSectionSpan returns the body of the section introduced by the
// given heading line, up to the next heading of the same or a higher level.
//
// Exactly-one is asserted before the body is returned: a heading matching zero
// or two lines would otherwise let every containment assertion pass vacuously.
func redNowExtractSectionSpan(content, heading string) (string, error) {
	lines := strings.Split(content, "\n")
	var starts []int
	for i, ln := range lines {
		if strings.TrimSpace(ln) == heading {
			starts = append(starts, i)
		}
	}
	if len(starts) != 1 {
		return "", fmt.Errorf("expected exactly one %q heading, got %d", heading, len(starts))
	}
	level := redNowHeadingLevel(heading)
	if level == 0 {
		return "", fmt.Errorf("%q is not a markdown heading", heading)
	}
	end := len(lines)
	for i := starts[0] + 1; i < len(lines); i++ {
		if l := redNowHeadingLevel(lines[i]); l != 0 && l <= level {
			end = i
			break
		}
	}
	body := strings.Join(lines[starts[0]:end], "\n")
	if strings.TrimSpace(strings.TrimPrefix(body, heading)) == "" {
		return "", fmt.Errorf("section %q is empty", heading)
	}
	return body, nil
}

func redNowMustExtractSection(t *testing.T, content, heading string) string {
	t.Helper()
	body, err := redNowExtractSectionSpan(content, heading)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return body
}

// ---------------------------------------------------------------------------
// the extracted form contract
// ---------------------------------------------------------------------------

// redNowUnquotedMask marks the byte positions of a command that sit outside any
// quoted span. A GFM table cell and a fenced ledger entry both routinely carry
// a literal `|` inside a quoted regex; treating that as a shell pipe would
// refuse commands that are in fact single invocations.
func redNowUnquotedMask(cmd string) []bool {
	mask := make([]bool, len(cmd))
	var inSingle, inDouble, escaped bool
	for i := 0; i < len(cmd); i++ {
		c := cmd[i]
		switch {
		case escaped:
			escaped = false
		case c == '\\' && !inSingle:
			escaped = true
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"' && !inSingle:
			inDouble = !inDouble
		}
		mask[i] = !inSingle && !inDouble
	}
	return mask
}

// redNowFormViolations returns the forbidden metacharacters that appear
// UNQUOTED in the command, i.e. the ones that would actually make it more than
// one shell invocation.
func redNowFormViolations(cmd string, forbidden []string) []string {
	mask := redNowUnquotedMask(cmd)
	var hits []string
	for _, tok := range forbidden {
		for i := 0; i+len(tok) <= len(cmd); i++ {
			if !strings.HasPrefix(cmd[i:], tok) {
				continue
			}
			allBare := true
			for j := i; j < i+len(tok); j++ {
				if !mask[j] {
					allBare = false
					break
				}
			}
			if allBare {
				hits = append(hits, tok)
				break
			}
		}
	}
	return hits
}

// ---------------------------------------------------------------------------
// carrier-independent command collection
// ---------------------------------------------------------------------------

type redNowCommand struct {
	Carrier string // "table-cell" | "ledger-entry" | "fenced-block"
	Command string
	Line    int
}

var (
	redNowLedgerRe    = regexp.MustCompile(`^E-(\d+)\s\s*(\S.*)$`)
	redNowDollarRe    = regexp.MustCompile(`^\s*\$ (\S.*)$`)
	redNowInlineRe    = regexp.MustCompile("`([^`]+)`")
	redNowLedgerRefRe = regexp.MustCompile(`E-(\d+)`)
	redNowShaRe       = regexp.MustCompile("`[0-9a-f]{7,40}`")
)

// redNowCollectCommands walks an acceptance.md and returns every command it
// carries, whatever the carrier. Three carriers are recognised, and the scan
// is the same scan for all three — which is what makes the predicate
// carrier-independent by construction rather than by promise.
func redNowCollectCommands(content string) []redNowCommand {
	var out []redNowCommand
	inFence := false
	for i, ln := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(ln), "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			if m := redNowLedgerRe.FindStringSubmatch(ln); m != nil {
				out = append(out, redNowCommand{Carrier: "ledger-entry", Command: m[2], Line: i + 1})
				continue
			}
			if m := redNowDollarRe.FindStringSubmatch(ln); m != nil {
				out = append(out, redNowCommand{Carrier: "fenced-block", Command: m[1], Line: i + 1})
			}
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(ln), "|") {
			for _, m := range redNowInlineRe.FindAllStringSubmatch(ln, -1) {
				body := m[1]
				if strings.HasPrefix(body, "$ ") {
					out = append(out, redNowCommand{Carrier: "table-cell", Command: strings.TrimPrefix(body, "$ "), Line: i + 1})
				}
			}
			continue
		}
		if m := redNowDollarRe.FindStringSubmatch(ln); m != nil {
			out = append(out, redNowCommand{Carrier: "fenced-block", Command: m[1], Line: i + 1})
		}
	}
	return out
}

// redNowCollectCellCommandsOnly is the DELIBERATELY narrowed predicate the
// carrier-relocation mutant (M-5) walks through. It exists so the mutant can be
// observed surviving it while the real predicate catches it — rule §5: observe
// the two forms diverge, do not grep for the fixed form.
func redNowCollectCellCommandsOnly(content string) []redNowCommand {
	var out []redNowCommand
	for _, c := range redNowCollectCommands(content) {
		if c.Carrier == "table-cell" {
			out = append(out, c)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// the four-element check
// ---------------------------------------------------------------------------

type redNowLedgerEntry struct {
	Command   string
	HasStdout bool
	HasExit   bool
}

func redNowParseLedger(content string) map[string]redNowLedgerEntry {
	out := map[string]redNowLedgerEntry{}
	inFence := false
	current := ""
	for _, ln := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(ln), "```") {
			inFence = !inFence
			current = ""
			continue
		}
		if !inFence {
			continue
		}
		if m := redNowLedgerRe.FindStringSubmatch(ln); m != nil {
			current = "E-" + m[1]
			out[current] = redNowLedgerEntry{Command: m[2]}
			continue
		}
		if current == "" {
			continue
		}
		e := out[current]
		switch {
		case strings.HasPrefix(strings.TrimSpace(ln), "stdout:"):
			e.HasStdout = true
		case strings.HasPrefix(strings.TrimSpace(ln), "exit:"):
			e.HasExit = true
		}
		out[current] = e
	}
	return out
}

type redNowRow struct {
	ID      string
	Class   string
	RedCell string
	Line    int
}

var redNowRowRe = regexp.MustCompile(`^\|\s*\*\*(AC-[A-Z0-9-]+)\*\*\s*\|`)

// redNowParseRows returns the acceptance matrix rows. The RED-now proof column
// is the second-from-last cell of the six-column matrix.
func redNowParseRows(content string) []redNowRow {
	var out []redNowRow
	for i, ln := range strings.Split(content, "\n") {
		m := redNowRowRe.FindStringSubmatch(ln)
		if m == nil {
			continue
		}
		cells := strings.Split(strings.Trim(strings.TrimSpace(ln), "|"), "|")
		if len(cells) < 6 {
			continue
		}
		out = append(out, redNowRow{
			ID:      m[1],
			Class:   strings.TrimSpace(cells[1]),
			RedCell: strings.TrimSpace(cells[len(cells)-2]),
			Line:    i + 1,
		})
	}
	return out
}

func redNowIsReleaseBlocking(class string) bool {
	return strings.Contains(class, "release-blocking") && !strings.Contains(class, "regression-guard")
}

// redNowElementFindings reports every release-blocking row whose RED-now cell
// does not carry the four elements — the command, its verbatim stdout, its exit
// code, and a tree SHA. The elements may live in the cell or in a fenced
// evidence-ledger entry the cell cites by id; the check follows them to
// whichever carrier holds them.
func redNowElementFindings(content string) []string {
	ledger := redNowParseLedger(content)
	docPinned := redNowShaRe.MatchString(content)
	var out []string
	for _, row := range redNowParseRows(content) {
		if !redNowIsReleaseBlocking(row.Class) {
			continue
		}
		refs := redNowLedgerRefRe.FindAllString(row.RedCell, -1)
		resolved := false
		for _, ref := range refs {
			e, ok := ledger[ref]
			if !ok {
				out = append(out, fmt.Sprintf("%s (line %d): cites %s, which no ledger entry defines", row.ID, row.Line, ref))
				continue
			}
			if e.Command == "" || !e.HasStdout || !e.HasExit {
				out = append(out, fmt.Sprintf("%s (line %d): ledger entry %s is missing the command, its stdout, or its exit code", row.ID, row.Line, ref))
				continue
			}
			resolved = true
		}
		if !resolved {
			out = append(out, fmt.Sprintf("%s (line %d): release-blocking RED-now cell carries no command, stdout and exit code", row.ID, row.Line))
			continue
		}
		if !docPinned && !redNowShaRe.MatchString(row.RedCell) {
			out = append(out, fmt.Sprintf("%s (line %d): RED-now cell has no tree SHA and the document carries no pin to inherit", row.ID, row.Line))
		}
	}
	return out
}

// redNowFormFindings reports every command in the file — any carrier, any
// class — that is not a single shell invocation.
func redNowFormFindings(content string, forbidden []string) []string {
	var out []string
	for _, c := range redNowCollectCommands(content) {
		if hits := redNowFormViolations(c.Command, forbidden); len(hits) > 0 {
			out = append(out, fmt.Sprintf("line %d (%s): %q carries unquoted %v", c.Line, c.Carrier, c.Command, hits))
		}
	}
	return out
}

func redNowFixture(t *testing.T, name string) string {
	t.Helper()
	rel := filepath.Join("internal/spec", redNowFixtureDir, name, "acceptance.md")
	return redNowRead(t, rel)
}

// ===========================================================================
// L1 — the rule clause (AC-RNT-001, AC-RNT-003)
// ===========================================================================

// TestRuleClauseEnumeratesFourElements covers AC-RNT-001. Every predicate is
// asserted INSIDE the extracted §2 span, so a token pasted elsewhere in the
// file does not satisfy it. It does not defeat a token pasted inside the span
// — that residual is mutant M-3, recorded in acceptance.md §D.2.
func TestRuleClauseEnumeratesFourElements(t *testing.T) {
	span := redNowMustExtractSection(t, redNowRead(t, redNowRulePath), redNowRuleSectionHead)
	for _, want := range []string{
		"RED-now cell content",
		"command", "stdout", "exit", "SHA",
		"read-only", "single invocation", "raw file",
	} {
		if !strings.Contains(span, want) {
			t.Errorf("§2 span does not name %q", want)
		}
	}
	for _, carrier := range []string{"table cell", "evidence-ledger"} {
		if !strings.Contains(span, carrier) {
			t.Errorf("§2 span does not admit the carrier %q", carrier)
		}
	}
}

// TestRuleClauseStatesDemotionNotPass covers AC-RNT-003.
func TestRuleClauseStatesDemotionNotPass(t *testing.T) {
	span := redNowMustExtractSection(t, redNowRead(t, redNowRulePath), redNowRuleSectionHead)
	if !strings.Contains(span, "regression-guard") {
		t.Errorf("§2 span does not name the regression-guard disposition")
	}
	notPass := regexp.MustCompile(`not\s+(?:be\s+)?record(?:ed)?\s+as\s+a\s+pass`)
	if !notPass.MatchString(span) {
		t.Errorf("§2 span does not state that an undecidable RED is not recorded as a pass")
	}
	if !strings.Contains(span, "loses release-blocking eligibility") {
		t.Errorf("§2 span does not state the demotion as the disposition")
	}
}

// TestRuleClauseIsStructuralNotLexical covers REQ-RNT-002 on the rule surface:
// the clause must not key on the tense, mood, or a word list of the prose.
func TestRuleClauseIsStructuralNotLexical(t *testing.T) {
	content := redNowRead(t, redNowRulePath)
	for _, banned := range []string{"tense", "mood", "counterfactual"} {
		if strings.Contains(content, banned) {
			t.Errorf("%s carries the lexical discriminator %q", redNowRulePath, banned)
		}
	}
}

// ===========================================================================
// L2 / L3 — the MP-8 clause (AC-RNT-004..007, -013, -014, -015)
// ===========================================================================

// ===========================================================================
// L1 — the form check over fixtures (AC-RNT-008, -009a, -009b)
// ===========================================================================

// TestRedNowViolatingFixtureIsReported covers AC-RNT-009a.
func TestRedNowViolatingFixtureIsReported(t *testing.T) {
	forbidden := redNowForbidden
	content := redNowFixture(t, "violating")

	elements := redNowElementFindings(content)
	if len(elements) == 0 {
		t.Errorf("violating fixture: the prose-only release-blocking RED cell was not reported")
	}
	form := redNowFormFindings(content, forbidden)
	if len(form) == 0 {
		t.Errorf("violating fixture: the piped command was not reported")
	}
	t.Logf("violating fixture findings: elements=%v form=%v", elements, form)
}

// TestRedNowLegitimateFixtureIsClean covers AC-RNT-009b. Confirming only the
// fail direction is indistinguishable from a check that reports everything.
func TestRedNowLegitimateFixtureIsClean(t *testing.T) {
	forbidden := redNowForbidden
	content := redNowFixture(t, "legitimate")

	if got := redNowElementFindings(content); len(got) != 0 {
		t.Errorf("legitimate fixture reported element findings: %v", got)
	}
	if got := redNowFormFindings(content, forbidden); len(got) != 0 {
		t.Errorf("legitimate fixture reported form findings: %v", got)
	}
	if n := len(redNowCollectCommands(content)); n == 0 {
		t.Errorf("legitimate fixture yielded zero commands — a clean report over an empty set asserts nothing")
	}
}

// TestCommandScopeIsCarrierIndependent covers AC-RNT-008 and observes mutant
// M-5. The ledger fixture carries its malformed command in NO table cell, so
// the deliberately narrowed cell-scoped predicate finds nothing to check. The
// two predicates are run against the same input and observed to diverge.
func TestCommandScopeIsCarrierIndependent(t *testing.T) {
	forbidden := redNowForbidden
	content := redNowFixture(t, "ledger")

	carrierIndependent := redNowFormFindings(content, forbidden)
	if len(carrierIndependent) == 0 {
		t.Fatalf("carrier-independent scope missed the malformed ledger command")
	}

	var cellScoped []string
	for _, c := range redNowCollectCellCommandsOnly(content) {
		if hits := redNowFormViolations(c.Command, forbidden); len(hits) > 0 {
			cellScoped = append(cellScoped, c.Command)
		}
	}
	if len(cellScoped) != 0 {
		t.Fatalf("fixture no longer isolates the mutant: the cell-scoped predicate found %v", cellScoped)
	}
	t.Logf("M-5 observed: cell-scoped=0 findings, carrier-independent=%d findings %v",
		len(carrierIndependent), carrierIndependent)

	// The carriers themselves are asserted so a scan that silently stopped
	// recognising one of them would surface here.
	seen := map[string]int{}
	for _, c := range redNowCollectCommands(content) {
		seen[c.Carrier]++
	}
	if seen["ledger-entry"] == 0 {
		t.Errorf("ledger carrier not recognised at all: %v", seen)
	}
}

// TestClassLaunderingMutantIsDetected observes mutant M-4: a malformed command
// moved into a regression-guard criterion. The class-scoped predicate misses it
// by construction; the class-independent one — the one this SPEC adopts —
// reports it. Both are run on the same input.
func TestClassLaunderingMutantIsDetected(t *testing.T) {
	forbidden := redNowForbidden
	content := redNowFixture(t, "violating")

	classIndependent := redNowFormFindings(content, forbidden)
	if len(classIndependent) == 0 {
		t.Fatalf("class-independent scope missed the laundered command")
	}

	// The narrowed predicate: only commands cited by a release-blocking row.
	ledger := redNowParseLedger(content)
	var classScoped []string
	for _, row := range redNowParseRows(content) {
		if !redNowIsReleaseBlocking(row.Class) {
			continue
		}
		for _, ref := range redNowLedgerRefRe.FindAllString(row.RedCell, -1) {
			e, ok := ledger[ref]
			if !ok {
				continue
			}
			if hits := redNowFormViolations(e.Command, forbidden); len(hits) > 0 {
				classScoped = append(classScoped, e.Command)
			}
		}
	}
	if len(classScoped) != 0 {
		t.Fatalf("fixture no longer isolates the mutant: the class-scoped predicate found %v", classScoped)
	}
	t.Logf("M-4 observed: class-scoped=0 findings, class-independent=%d findings %v",
		len(classIndependent), classIndependent)
}

// TestRedNowFormCheckDivergesOnQuoting is the rule §5 audit-verification form:
// a quoted `|` is a literal and must NOT be refused, while an unquoted one must
// be. Without observing both directions, a checker that refuses everything and
// a checker that refuses nothing are indistinguishable.
func TestRedNowFormCheckDivergesOnQuoting(t *testing.T) {
	forbidden := redNowForbidden
	cases := []struct {
		cmd  string
		want bool
	}{
		{`grep -c "^| \*\*AC-RNT-" acceptance.md`, false},
		{`grep -c 'a|b' file.md`, false},
		{`grep -c alpha file.md | wc -l`, true},
		{`ls target.txt && echo ok`, true},
		{`ls target.txt ; echo done`, true},
		{`grep -c alpha file.md > out.txt`, true},
		{`ls internal/spec/red_now_cell_test.go`, false},
	}
	for _, tc := range cases {
		got := len(redNowFormViolations(tc.cmd, forbidden)) > 0
		if got != tc.want {
			t.Errorf("form check on %q = %v, want %v (forbidden=%v)", tc.cmd, got, tc.want, forbidden)
		}
	}
}

// ===========================================================================
// mirrors and neutrality (AC-RNT-010, -011, -012)
// ===========================================================================

// TestRuleMirrorIsByteIdentical covers REQ-RNT-010 for the rule pair, which was
// byte-identical before this work and must stay so.
func TestRuleMirrorIsByteIdentical(t *testing.T) {
	if redNowRead(t, redNowRulePath) != redNowRead(t, redNowRuleMirrorPath) {
		t.Errorf("%s and %s are not byte-identical", redNowRulePath, redNowRuleMirrorPath)
	}
}

// TestRedNowArtifactsDoNotShip covers AC-RNT-012.
func TestRedNowArtifactsDoNotShip(t *testing.T) {
	root := filepath.Join(repoRoot(t), "internal/template/templates")
	var hits []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), "red_now") {
			hits = append(hits, path)
			return nil
		}
		raw, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		if strings.Contains(string(raw), "red_now") {
			hits = append(hits, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk template tree: %v", err)
	}
	sort.Strings(hits)
	if len(hits) != 0 {
		t.Errorf("the repository-local test or its fixtures leaked into the template tree: %v", hits)
	}
}
