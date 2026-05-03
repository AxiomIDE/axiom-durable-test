package nodes_test

import (
	"context"
	"testing"

	"axiom-official/axiom-durable-test/axiom"
)

// testContext is a testing.T-backed axiom.Context shared by every *_test.go in
// this package. Each test file's CLI scaffold defined its own copy; that
// produced duplicate declarations once the package had 2+ nodes, so we keep a
// single canonical definition here.
type testContext struct {
	t          *testing.T
	secretsMap map[string]string
}

func newTestContext(t *testing.T) *testContext {
	return &testContext{t: t, secretsMap: map[string]string{}}
}

type testLogger struct{ t *testing.T }

func (l *testLogger) Debug(msg string, args ...any) { l.t.Logf("DEBUG  %s %v", msg, args) }
func (l *testLogger) Info(msg string, args ...any)  { l.t.Logf("INFO   %s %v", msg, args) }
func (l *testLogger) Warn(msg string, args ...any)  { l.t.Logf("WARN   %s %v", msg, args) }
func (l *testLogger) Error(msg string, args ...any) { l.t.Logf("ERROR  %s %v", msg, args) }

type testSecrets struct{ m map[string]string }

func (s testSecrets) Get(name string) (string, bool) { v, ok := s.m[name]; return v, ok }

type testAgent struct{}

func (testAgent) Memory() axiom.AgentMemory { return testAgentMemory{} }

type testAgentMemory struct{}

func (testAgentMemory) Session(_ string) axiom.SessionMemory { return testSessionMemory{} }
func (testAgentMemory) Search(_ context.Context, _ string, _ int) ([]axiom.MemoryEntry, error) {
	return nil, nil
}
func (testAgentMemory) Write(_ context.Context, _ string, _ float32) (string, error) {
	return "", nil
}

type testSessionMemory struct{}

func (testSessionMemory) Search(_ context.Context, _ string, _ int) ([]axiom.MemoryEntry, error) {
	return nil, nil
}
func (testSessionMemory) Write(_ context.Context, _ string, _ float32) (string, error) {
	return "", nil
}
func (testSessionMemory) History() axiom.SessionHistory { return testSessionHistory{} }
func (testSessionMemory) End(_ context.Context) error   { return nil }

type testSessionHistory struct{}

func (testSessionHistory) Last(_ context.Context, _ int) ([]axiom.ConversationTurn, error) {
	return nil, nil
}
func (testSessionHistory) Append(_ context.Context, _, _ string) error { return nil }

func (c *testContext) Log() axiom.Logger      { return &testLogger{c.t} }
func (c *testContext) Secrets() axiom.Secrets { return testSecrets{c.secretsMap} }
func (c *testContext) Agent() axiom.Agent     { return testAgent{} }
func (c *testContext) ExecutionID() string    { return "test-execution-id" }
func (c *testContext) FlowID() string         { return "test-flow-id" }
func (c *testContext) TenantID() string       { return "test-tenant-id" }

var _ axiom.Context = (*testContext)(nil)
