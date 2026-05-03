package nodes

import (
	"context"
	"errors"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// MaybeFail returns input.error_message as an error when fail=true. When
// fail=false it echoes input.text. Drives gate failure-policy scenarios.
func MaybeFail(ctx context.Context, ax axiom.Context, input *gen.MaybeFailRequest) (*gen.Message, error) {
	_, _ = ctx, ax
	if input.GetFail() {
		msg := input.GetErrorMessage()
		if msg == "" {
			msg = "MaybeFail requested failure"
		}
		return nil, errors.New(msg)
	}
	return &gen.Message{Text: input.GetText()}, nil
}
