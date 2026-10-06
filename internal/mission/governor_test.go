package mission

import (
	"testing"
)

func TestAutoMissionQuestionBoundary(t *testing.T) {
	inside := SuppressAutoMissionQuestions(true, "")
	if inside.Blocked || inside.AskUserQuestion || !inside.Proceed {
		t.Fatalf("inside boundary = %+v", inside)
	}
	outside := SuppressAutoMissionQuestions(false, "new authority required")
	if !outside.Blocked || outside.AskUserQuestion || outside.Proceed || outside.Report == "" {
		t.Fatalf("outside boundary = %+v", outside)
	}
}
