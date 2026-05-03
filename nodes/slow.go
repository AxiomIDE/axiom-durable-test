package nodes

import (
	"context"
	"time"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// Slow sleeps sleep_ms milliseconds (or 500ms default) before returning, so
// breakpoints and pauses can land deterministically without racing execution.
func Slow(ctx context.Context, ax axiom.Context, input *gen.SlowRequest) (*gen.Message, error) {
	d := time.Duration(input.GetSleepMs()) * time.Millisecond
	if d <= 0 {
		d = 500 * time.Millisecond
	}
	ax.Log().Info("Slow sleeping", "ms", d.Milliseconds())
	select {
	case <-time.After(d):
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return &gen.Message{Text: input.GetText(), Value: int32(d.Milliseconds())}, nil
}
