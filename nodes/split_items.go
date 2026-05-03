package nodes

import (
	"context"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// SplitItems echoes the items array as-is — falling back to three defaults
// when the input list is empty so a graph can fan out without ceremony. The
// graph's gate_config.dynamic_fan_out uses input_items_path = "value.items"
// to derive N branches from this output.
func SplitItems(ctx context.Context, ax axiom.Context, input *gen.SplitRequest) (*gen.ItemBatch, error) {
	_, _ = ctx, ax
	items := input.GetItems()
	if len(items) == 0 {
		items = []string{"alpha", "bravo", "charlie"}
	}
	return &gen.ItemBatch{Text: input.GetText(), Items: items}, nil
}
