package nodes

import (
	"context"
	"fmt"
	"strings"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// Decide is the HITL resume anchor: its INPUT message carries one field of
// every kind the console's resume form renders, so a flow can pause on a
// Decide placement, declare `hitl.resume_input_schema: ReviewDecision`, and
// have the reviewer's submitted value actually become this node's input.
//
// It echoes `text` through unchanged so a downstream Echo renders
// `echo:<what the reviewer typed>` — the assertion that proves the resume
// value reached the node rather than the flow's pre-pause input
// (AxiomIDE/poc#1354). `note` is rewritten to a deterministic, greppable
// summary of the NON-string fields so a spec can prove the bool, repeated,
// nested and enum halves of the form survived the wire round-trip too,
// through the same single `Message.text`-shaped surface the SPA already
// renders. `value` passes through so the int32 is observable as itself.
//
// Deterministic by construction: no clock, no randomness, no I/O.
func Decide(ctx context.Context, ax axiom.Context, input *gen.ReviewDecision) (*gen.Message, error) {
	_, _ = ctx, ax

	reviewer := "-"
	if r := input.GetReviewer(); r != nil && (r.GetName() != "" || r.GetRole() != "") {
		reviewer = fmt.Sprintf("%s/%s", r.GetName(), r.GetRole())
	}

	items := "-"
	if len(input.GetItems()) > 0 {
		items = strings.Join(input.GetItems(), "|")
	}

	return &gen.Message{
		Text:  input.GetText(),
		Value: input.GetValue(),
		Note: fmt.Sprintf(
			"approved=%t status=%s items=%s reviewer=%s note=%s",
			input.GetApproved(),
			input.GetStatus().String(),
			items,
			reviewer,
			input.GetNote(),
		),
	}, nil
}
