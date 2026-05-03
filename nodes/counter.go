package nodes

import (
	"context"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// Counter is deterministic: step = input.value + 1, total = step * step.
// Useful as the modify-state target so the inspector has a meaningful payload
// to override during fork scenarios.
func Counter(ctx context.Context, ax axiom.Context, input *gen.Message) (*gen.CounterState, error) {
	_, _ = ctx, ax
	step := input.GetValue() + 1
	return &gen.CounterState{Step: step, Total: step * step, LastText: input.GetText()}, nil
}
