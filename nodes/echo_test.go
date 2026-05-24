package nodes_test

import (
	"context"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestEcho_PrefixesText(t *testing.T) {
	got, err := nodes.Echo(context.Background(), newTestContext(t), &gen.Message{Text: "hi", Value: 9, Note: "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetText() != "echo:hi" {
		t.Errorf("expected echo:hi, got %q", got.GetText())
	}
	if got.GetValue() != 9 {
		t.Errorf("Echo should preserve value; got %d", got.GetValue())
	}
	if got.GetNote() != "world" {
		t.Errorf("Echo should preserve note; got %q", got.GetNote())
	}
}
