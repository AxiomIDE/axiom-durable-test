package nodes_test

import (
	"context"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestCounter_DeterministicValues(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	cases := []struct {
		in    int32
		step  int32
		total int32
	}{
		{0, 1, 1},
		{4, 5, 25},
		{9, 10, 100},
	}
	for _, tc := range cases {
		got, err := nodes.Counter(ctx, ax, &gen.Message{Value: tc.in, Text: "k"})
		if err != nil {
			t.Fatalf("unexpected error for value=%d: %v", tc.in, err)
		}
		if got.GetStep() != tc.step || got.GetTotal() != tc.total || got.GetLastText() != "k" {
			t.Errorf("Counter(%d) = %+v, want step=%d total=%d", tc.in, got, tc.step, tc.total)
		}
	}
}
