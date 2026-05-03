package nodes

import (
	"context"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// Identity passes its input through unchanged. Useful as a step counter and
// breakpoint anchor in the durable-engine manual test scenarios.
func Identity(ctx context.Context, ax axiom.Context, input *gen.Message) (*gen.Message, error) {
	_ = ctx
	ax.Log().Info("Identity", "text", input.GetText(), "value", input.GetValue())
	return &gen.Message{Text: input.GetText(), Value: input.GetValue(), Note: input.GetNote()}, nil
}
