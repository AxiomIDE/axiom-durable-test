package nodes_test

import (
	"context"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestIdentity(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)
	in := &gen.Message{Text: "hello", Value: 42, Note: "n"}

	got, err := nodes.Identity(ctx, ax, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetText() != "hello" || got.GetValue() != 42 || got.GetNote() != "n" {
		t.Errorf("Identity should pass through unchanged; got %+v", got)
	}
}
