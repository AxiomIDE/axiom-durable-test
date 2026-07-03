package nodes

import (
	"context"
	"fmt"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// StreamFrameCount is the fixed number of frames StreamEcho emits per input
// message. Exported so e2e-gate fixtures can assert the exact observed frame
// count without hardcoding a magic number on both sides.
const StreamFrameCount = 4

// StreamEcho is the streaming pipeline witness node: for a single input
// Message it emits StreamFrameCount StreamFrame frames carrying
// "<input.text>-frame<i>", the last of which sets is_final=true and carries
// "<input.text>-FINAL". type: pipeline nodes take an input channel + an emit
// callback instead of a single request/response — the platform auto-appends
// the transport-level terminal NodeResponse (empty payload, IsFinal=true)
// once this function returns, so the "final frame" the caller cares about is
// a business-level field on the last emitted message, not a transport flag.
// inputs: stream of gen.Message frames from the previous node.
// For the start node of a pipeline flow, the channel yields exactly one item.
func StreamEcho(ctx context.Context, ax axiom.Context, in <-chan *gen.Message, emit func(*gen.StreamFrame) error) error {
	for input := range in {
		text := input.GetText()
		for i := 0; i < StreamFrameCount; i++ {
			isFinal := i == StreamFrameCount-1
			frameText := fmt.Sprintf("%s-frame%d", text, i)
			if isFinal {
				frameText = fmt.Sprintf("%s-FINAL", text)
			}
			ax.Log().Info("StreamEcho frame", "index", i, "is_final", isFinal)
			if err := emit(&gen.StreamFrame{
				Index:   int32(i),
				Text:    frameText,
				IsFinal: isFinal,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
