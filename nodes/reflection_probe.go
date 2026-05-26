package nodes

import (
	"context"
	"encoding/json"

	"axiom-official/axiom-durable-test/axiom"
	gen "axiom-official/axiom-durable-test/gen"
)

// ADR-050 (2026-05-26): ReflectionProbe is the demo / e2e-gate witness node
// for `ax.reflection.flow.*`. It reads the reflection surface that the
// worker attaches to every SidecarRequest and the sidecar forwards on every
// NodeRequest, then JSON-serializes a snapshot of what it saw into
// Message.text so the e2e-gate spec can assert on it.
//
// Why JSON-in-text instead of a typed snapshot message: keeps the package
// proto schema unchanged so the worker / SPA / registry don't need to know
// about a new message type. The cost is the e2e spec has to JSON.parse the
// terminal output, which Playwright handles natively.
//
// Output schema (kept stable for e2e assertions):
//
//	{
//	  "graph_id":              "01HX...",
//	  "node_count":            <int>,
//	  "edge_count":            <int>,
//	  "loop_edge_count":       <int>,
//	  "current_instance":      <uint>,
//	  "depth":                 <uint>,
//	  "subflow_stack":         [<graph_id>, ...],
//	  "loop_iterations":       {<dst_instance>: <count>, ...},
//	  "my_node_name":          "<registry name of the instance currently executing>",
//	  "my_node_type":          "node" | "subflow" | "pipeline",
//	  "my_canvas_node_id":     "<canvas id>",
//	  "downstream_node_names": ["<successor1>", "<successor2>", ...],
//	  "any_conditional_edge":  <bool — true if any edge has a condition>
//	}
//
// Input.value and input.note are passed through unchanged so a downstream
// node in the e2e graph can still see them — the probe is non-destructive.
func ReflectionProbe(ctx context.Context, ax axiom.Context, input *gen.Message) (*gen.Message, error) {
	_ = ctx
	flow := ax.Reflection().Flow()
	pos := flow.Position()

	// Resolve "my own" node from the reflection node list keyed by instance id.
	var myName, myType, myCanvasID string
	for _, n := range flow.Nodes() {
		if n.InstanceID == pos.CurrentInstance {
			myName = n.Name
			myType = n.NodeType
			myCanvasID = n.CanvasNodeID
			break
		}
	}

	// Find downstream node names by walking edges from the current instance.
	var downstream []string
	anyConditional := false
	nodeNameByInstance := make(map[uint32]string, len(flow.Nodes()))
	for _, n := range flow.Nodes() {
		nodeNameByInstance[n.InstanceID] = n.Name
	}
	for _, e := range flow.Edges() {
		if e.HasCondition {
			anyConditional = true
		}
		if e.SrcInstance == pos.CurrentInstance {
			if name, ok := nodeNameByInstance[e.DstInstance]; ok {
				downstream = append(downstream, name)
			}
		}
	}
	if downstream == nil {
		downstream = []string{}
	}

	loopIterations := map[uint32]uint32{}
	for k, v := range pos.LoopIterations {
		loopIterations[k] = v
	}

	subStack := pos.SubflowStackGraphIDs
	if subStack == nil {
		subStack = []string{}
	}

	snapshot := map[string]any{
		"graph_id":              flow.GraphID(),
		"node_count":            len(flow.Nodes()),
		"edge_count":            len(flow.Edges()),
		"loop_edge_count":       len(flow.LoopEdges()),
		"current_instance":      pos.CurrentInstance,
		"depth":                 pos.Depth,
		"subflow_stack":         subStack,
		"loop_iterations":       loopIterations,
		"my_node_name":          myName,
		"my_node_type":          myType,
		"my_canvas_node_id":     myCanvasID,
		"downstream_node_names": downstream,
		"any_conditional_edge":  anyConditional,
	}
	out, err := json.Marshal(snapshot)
	if err != nil {
		return nil, err
	}

	ax.Log().Info("ReflectionProbe",
		"graph_id", flow.GraphID(),
		"node_count", len(flow.Nodes()),
		"current_instance", pos.CurrentInstance,
		"my_node_name", myName,
	)

	return &gen.Message{
		Text:  string(out),
		Value: input.GetValue(),
		Note:  input.GetNote(),
	}, nil
}
