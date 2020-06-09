package graph

import "testing"

func TestNodeConnect(t *testing.T) {
	nodeOrig := newNode("GRU")
	nodeDest := newNode("BRC")

	nodeOrig.connect(nodeDest, 10)
	if len(nodeOrig.connections) != 1 {
		t.Errorf("node.connectios expected size %v, got %v", 1, len(nodeOrig.connections))
	}

	connection, found := nodeOrig.connections[nodeDest.label]
	if !found {
		t.Errorf("node.connectios expected size %v, got %v", 1, len(nodeOrig.connections))
	}
	if connection.destination.label != "BRC" {
		t.Errorf("connection.destination.label expected %v, got %v", "BRC", connection.destination.label)
	}
}

func TestNodeMultipleConnections(t *testing.T) {
	nodes := []*Node{newNode("GRU"), newNode("BRC"), newNode("SCL"), newNode("CDG")}

	nodes[0].connect(nodes[1], 10)
	nodes[1].connect(nodes[2], 5)
	nodes[0].connect(nodes[3], 75)
	nodes[0].connect(nodes[2], 20)

	if len(nodes[0].connections) != 3 {
		t.Errorf("nodes[0].connections expected size %v, got %v", 3, len(nodes[0].connections))
	}

	if len(nodes[1].connections) != 1 {
		t.Errorf("nodes[1].connections expected size %v, got %v", 3, len(nodes[1].connections))
	}

	if len(nodes[2].connections) != 0 {
		t.Errorf("nodes[2].connections expected size %v, got %v", 0, len(nodes[2].connections))
	}

	if len(nodes[3].connections) != 0 {
		t.Errorf("nodes[3].connections expected size %v, got %v", 0, len(nodes[3].connections))
	}
}

func TestGraphConnections(t *testing.T) {
	var expected = []struct {
		label               string
		expectedConnections int
	}{
		{"GRU", 4},
		{"BRC", 1},
		{"SCL", 1},
		{"CDG", 0},
		{"ORL", 1},
	}

	graph := NewGraph()
	graph.Connect("GRU", "BRC", 10)
	graph.Connect("BRC", "SCL", 5)
	graph.Connect("GRU", "CDG", 75)
	graph.Connect("GRU", "SCL", 20)
	graph.Connect("GRU", "ORL", 56)
	graph.Connect("ORL", "CDG", 5)
	graph.Connect("SCL", "ORL", 20)

	for _, e := range expected {
		node, found := graph.nodes[e.label]
		if !found {
			t.Errorf("graph.nodes[%v] not found", e.label)
		}
		if len(node.connections) != e.expectedConnections {
			t.Errorf("graph.nodes[%v].connections expected size %v, got %v", e.label, e.expectedConnections, len(node.connections))
		}
	}
}

func TestGraphShortestPath(t *testing.T) {
	expectedRoute := []string{
		"GRU", "BRC", "SCL", "ORL", "CDG",
	}
	expectedCost := float32(40)

	graph := NewGraph()
	graph.Connect("GRU", "BRC", 10)
	graph.Connect("BRC", "SCL", 5)
	graph.Connect("GRU", "CDG", 75)
	graph.Connect("GRU", "SCL", 20)
	graph.Connect("GRU", "ORL", 56)
	graph.Connect("ORL", "CDG", 5)
	graph.Connect("SCL", "ORL", 20)

	route, cost := graph.ShortestPath("GRU", "CDG")
	for i := 0; i < len(expectedRoute); i++ {
		if route[i] != expectedRoute[i] {
			t.Errorf("graph.ShortestPath expected route %v, got %v", expectedRoute, route)
		}
	}

	if cost != expectedCost {
		t.Errorf("graph.ShortestPath expected cost %v, got %v", expectedCost, cost)
	}
}