package nodes_test

import (
	"context"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestDecide_PassesTextAndValueThroughUnchanged(t *testing.T) {
	got, err := nodes.Decide(context.Background(), newTestContext(t), &gen.ReviewDecision{
		Text:  "reviewer-typed-this",
		Value: 42,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// text passes through so a downstream Echo renders `echo:<submitted text>`
	// — the e2e assertion that the HITL resume value reached the node.
	if got.GetText() != "reviewer-typed-this" {
		t.Errorf("Decide must pass text through unchanged; got %q", got.GetText())
	}
	if got.GetValue() != 42 {
		t.Errorf("Decide must pass value through unchanged; got %d", got.GetValue())
	}
}

func TestDecide_SummarisesEveryNonStringFieldIntoNote(t *testing.T) {
	got, err := nodes.Decide(context.Background(), newTestContext(t), &gen.ReviewDecision{
		Text:     "t",
		Note:     "looks fine",
		Approved: true,
		Items:    []string{"alpha", "bravo", "charlie"},
		Reviewer: &gen.Reviewer{Name: "alice", Role: "lead"},
		Status:   gen.ApprovalStatus_NEEDS_CHANGES,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "approved=true status=NEEDS_CHANGES items=alpha|bravo|charlie reviewer=alice/lead note=looks fine"
	if got.GetNote() != want {
		t.Errorf("note summary mismatch:\n got %q\nwant %q", got.GetNote(), want)
	}
}

// The all-defaults case is what an AUTO_APPROVE/AUTO_REJECT timeout resumes
// with when default_resume_value is `{}`, so it must not panic and must stay
// greppable.
func TestDecide_ZeroValueInputIsTotal(t *testing.T) {
	got, err := nodes.Decide(context.Background(), newTestContext(t), &gen.ReviewDecision{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "approved=false status=APPROVAL_STATUS_UNSPECIFIED items=- reviewer=- note="
	if got.GetNote() != want {
		t.Errorf("zero-value note mismatch:\n got %q\nwant %q", got.GetNote(), want)
	}
	if got.GetText() != "" {
		t.Errorf("zero-value text should stay empty; got %q", got.GetText())
	}
}

// A nil Reviewer submessage (the shape when the reviewer never expands the
// nested form) must render as "-" rather than panicking.
func TestDecide_NilReviewerSubmessage(t *testing.T) {
	got, err := nodes.Decide(context.Background(), newTestContext(t), &gen.ReviewDecision{
		Text:     "t",
		Reviewer: nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetNote() != "approved=false status=APPROVAL_STATUS_UNSPECIFIED items=- reviewer=- note=" {
		t.Errorf("nil reviewer must render as '-'; got %q", got.GetNote())
	}
}
