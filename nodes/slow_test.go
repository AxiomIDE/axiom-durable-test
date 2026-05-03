package nodes_test

import (
	"context"
	"testing"
	"time"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestSlow_HonoursSleepMs(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)
	start := time.Now()
	got, err := nodes.Slow(ctx, ax, &gen.SlowRequest{Text: "x", SleepMs: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Errorf("Slow returned too quickly: %v", elapsed)
	}
	if got.GetText() != "x" || got.GetValue() != 50 {
		t.Errorf("unexpected output: %+v", got)
	}
}

func TestSlow_RespectsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := nodes.Slow(ctx, newTestContext(t), &gen.SlowRequest{Text: "x", SleepMs: 5000}); err == nil {
		t.Errorf("expected ctx.Err(), got nil")
	}
}
