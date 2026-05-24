package nodes

import (
	"context"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// Echo logs and prefixes the input text with "echo:" so HITL anchors and
// log-trace assertions have an obvious signal to look for. Value and Note
// are passed through unchanged so wave-021 compose demos can author plans
// that bind different fields to different incident edges and verify each
// field was sourced from the right edge.
func Echo(ctx context.Context, ax axiom.Context, input *gen.Message) (*gen.Message, error) {
	_ = ctx
	ax.Log().Info("Echo", "text", input.GetText())
	return &gen.Message{
		Text:  "echo:" + input.GetText(),
		Value: input.GetValue(),
		Note:  input.GetNote(),
	}, nil
}
