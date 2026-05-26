package nodes_test

import (
	"context"
	"encoding/json"
	"testing"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
	"axiom-official/axiom-durable-test/nodes"
)

// ADR-050 (2026-05-26): unit test that ReflectionProbe produces the
// expected JSON snapshot from a populated FlowReflection. The fixture
// mirrors the 2-node Identity → ReflectionProbe graph the e2e-gate spec
// drives end-to-end through the worker; this test pins the JSON shape
// without running the worker.

func TestReflectionProbe_Snapshot(t *testing.T) {
	ctx := context.Background()
	tc := newTestContext(t)
	tc.reflection = testReflection{
		flow: fixtureFlowReflection{
			graphID: "01HXTESTGRAPH00000000000000",
			nodes: []axiom.ReflectionNode{
				{
					InstanceID:   0,
					Name:         "Identity",
					NodeType:     "node",
					CanvasNodeID: "ident-0",
				},
				{
					InstanceID:   1,
					Name:         "ReflectionProbe",
					NodeType:     "node",
					CanvasNodeID: "probe-1",
				},
			},
			edges: []axiom.ReflectionEdge{
				{SrcInstance: 0, DstInstance: 1, CanvasEdgeID: "e0", HasCondition: false},
			},
			position: axiom.FlowPosition{
				CurrentInstance: 1, // we are the probe node
				Depth:           0,
			},
		},
	}

	in := &gen.Message{Text: "ignored", Value: 7, Note: "n"}
	got, err := nodes.ReflectionProbe(ctx, tc, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetValue() != 7 || got.GetNote() != "n" {
		t.Errorf("probe should pass value/note through; got value=%d note=%q",
			got.GetValue(), got.GetNote())
	}

	var snap map[string]any
	if err := json.Unmarshal([]byte(got.GetText()), &snap); err != nil {
		t.Fatalf("probe output is not valid JSON: %v\nraw: %s", err, got.GetText())
	}

	if snap["graph_id"] != "01HXTESTGRAPH00000000000000" {
		t.Errorf("graph_id = %v", snap["graph_id"])
	}
	if nc, _ := snap["node_count"].(float64); int(nc) != 2 {
		t.Errorf("node_count = %v want 2", snap["node_count"])
	}
	if ec, _ := snap["edge_count"].(float64); int(ec) != 1 {
		t.Errorf("edge_count = %v want 1", snap["edge_count"])
	}
	if ci, _ := snap["current_instance"].(float64); int(ci) != 1 {
		t.Errorf("current_instance = %v want 1", snap["current_instance"])
	}
	if d, _ := snap["depth"].(float64); int(d) != 0 {
		t.Errorf("depth = %v want 0", snap["depth"])
	}
	if name, _ := snap["my_node_name"].(string); name != "ReflectionProbe" {
		t.Errorf("my_node_name = %q want %q", name, "ReflectionProbe")
	}
	if myCanvas, _ := snap["my_canvas_node_id"].(string); myCanvas != "probe-1" {
		t.Errorf("my_canvas_node_id = %q want %q", myCanvas, "probe-1")
	}
	downstream, _ := snap["downstream_node_names"].([]any)
	if len(downstream) != 0 {
		t.Errorf("downstream_node_names should be empty for terminal probe; got %v", downstream)
	}
	if cond, _ := snap["any_conditional_edge"].(bool); cond {
		t.Errorf("any_conditional_edge should be false")
	}
}

func TestReflectionProbe_DownstreamAndCondition(t *testing.T) {
	// Variant: probe sits in the middle of the graph, so it has a
	// downstream successor; one edge in the graph has a condition.
	ctx := context.Background()
	tc := newTestContext(t)
	tc.reflection = testReflection{
		flow: fixtureFlowReflection{
			graphID: "g",
			nodes: []axiom.ReflectionNode{
				{InstanceID: 0, Name: "Identity", NodeType: "node"},
				{InstanceID: 1, Name: "ReflectionProbe", NodeType: "node"},
				{InstanceID: 2, Name: "Echo", NodeType: "node"},
			},
			edges: []axiom.ReflectionEdge{
				{SrcInstance: 0, DstInstance: 1, HasCondition: true},
				{SrcInstance: 1, DstInstance: 2},
			},
			position: axiom.FlowPosition{CurrentInstance: 1},
		},
	}
	got, err := nodes.ReflectionProbe(ctx, tc, &gen.Message{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	var snap map[string]any
	if err := json.Unmarshal([]byte(got.GetText()), &snap); err != nil {
		t.Fatalf("snap parse: %v", err)
	}
	downstream, _ := snap["downstream_node_names"].([]any)
	if len(downstream) != 1 || downstream[0] != "Echo" {
		t.Errorf("downstream = %v want [Echo]", downstream)
	}
	if cond, _ := snap["any_conditional_edge"].(bool); !cond {
		t.Errorf("any_conditional_edge should be true")
	}
}
