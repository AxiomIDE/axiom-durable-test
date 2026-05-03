package nodes_test

import (
	"context"
	"reflect"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestSplitItems_PassesThroughExplicitItems(t *testing.T) {
	got, err := nodes.SplitItems(context.Background(), newTestContext(t),
		&gen.SplitRequest{Text: "go", Items: []string{"a", "b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetText() != "go" || !reflect.DeepEqual(got.GetItems(), []string{"a", "b"}) {
		t.Errorf("unexpected output: %+v", got)
	}
}

func TestSplitItems_DefaultsWhenEmpty(t *testing.T) {
	got, err := nodes.SplitItems(context.Background(), newTestContext(t), &gen.SplitRequest{Text: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got.GetItems(), []string{"alpha", "bravo", "charlie"}) {
		t.Errorf("expected default trio, got %v", got.GetItems())
	}
}
