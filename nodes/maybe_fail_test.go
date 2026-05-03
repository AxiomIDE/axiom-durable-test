package nodes_test

import (
	"context"
	"testing"

	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

func TestMaybeFail_Success(t *testing.T) {
	got, err := nodes.MaybeFail(context.Background(), newTestContext(t),
		&gen.MaybeFailRequest{Text: "ok", Fail: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetText() != "ok" {
		t.Errorf("expected text=ok, got %q", got.GetText())
	}
}

func TestMaybeFail_FailureWithCustomMessage(t *testing.T) {
	_, err := nodes.MaybeFail(context.Background(), newTestContext(t),
		&gen.MaybeFailRequest{Fail: true, ErrorMessage: "boom"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "boom" {
		t.Errorf("expected error=boom, got %q", err.Error())
	}
}

func TestMaybeFail_FailureDefaultMessage(t *testing.T) {
	_, err := nodes.MaybeFail(context.Background(), newTestContext(t),
		&gen.MaybeFailRequest{Fail: true})
	if err == nil || err.Error() == "" {
		t.Errorf("expected non-empty default error, got %v", err)
	}
}
