package nodes_test

import (
	"context"
	"fmt"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

// TestStreamEcho drives the pipeline signature directly: seed the input
// channel with one Message, collect every emitted StreamFrame via the emit
// callback, then assert the exact contract the e2e-gate spec relies on —
// frame count, per-frame index/text, and is_final set on (only) the last one.
func TestStreamEcho(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	in := make(chan *gen.Message, 1)
	in <- &gen.Message{Text: "hello"}
	close(in)

	var frames []*gen.StreamFrame
	emit := func(f *gen.StreamFrame) error {
		frames = append(frames, f)
		return nil
	}

	if err := nodes.StreamEcho(ctx, ax, in, emit); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(frames) != nodes.StreamFrameCount {
		t.Fatalf("expected %d frames, got %d", nodes.StreamFrameCount, len(frames))
	}

	for i, f := range frames {
		if f.GetIndex() != int32(i) {
			t.Errorf("frame %d: expected index %d, got %d", i, i, f.GetIndex())
		}
		isLast := i == nodes.StreamFrameCount-1
		if f.GetIsFinal() != isLast {
			t.Errorf("frame %d: expected is_final=%v, got %v", i, isLast, f.GetIsFinal())
		}
		if isLast {
			if f.GetText() != "hello-FINAL" {
				t.Errorf("final frame: expected text %q, got %q", "hello-FINAL", f.GetText())
			}
		} else {
			expected := fmt.Sprintf("hello-frame%d", i)
			if f.GetText() != expected {
				t.Errorf("frame %d: expected text %q, got %q", i, expected, f.GetText())
			}
		}
	}
}

// TestStreamEcho_EmptyInput proves a zero-item input channel (no seed frame)
// returns cleanly with no frames emitted, rather than hanging or panicking —
// guards the `for input := range in` loop against an empty-channel edge case.
func TestStreamEcho_EmptyInput(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	in := make(chan *gen.Message)
	close(in)

	var frames []*gen.StreamFrame
	emit := func(f *gen.StreamFrame) error {
		frames = append(frames, f)
		return nil
	}

	if err := nodes.StreamEcho(ctx, ax, in, emit); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) != 0 {
		t.Fatalf("expected 0 frames for empty input, got %d", len(frames))
	}
}
