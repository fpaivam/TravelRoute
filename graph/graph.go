package graph

import (
	"TravelRoute/route"
	"container/list"
)

// connection represents a weighted oriented conenection
type connection struct {
	destination *Node
	weight      float32
}

func newConnection(destination *Node, weight float32) *connection {
	return &connection{destination: destination, weight: weight}
}

// Node represents a graph elemment that has weighted oriented conenections
type Node struct {
	label       string
	connections map[string]*connection
}

func newNode(label string) *Node {
	return &Node{label: label, connections: make(map[string]*connection)}
}

func (n *Node) connect(destination *Node, weigth float32) {
	connection, found := n.connections[destination.label]
	if !found {
		connection = newConnection(destination, weigth)
		n.connections[destination.label] = connection
	} else {
		connection.weight = weigth
	}
}

// Graph represents a graph
type Graph struct {
	nodes map[string]*Node
}

// NewGraph constructs a new Graph
func NewGraph() *Graph {
	return &Graph{nodes: make(map[string]*Node)}
}

// FindCheapestRoute Constructs a graph and finds the shortest (cheapest) route
// between origin and destination
func FindCheapestRoute(routes []route.Route, origin string, destination string) ([]string, float32) {
	routeGraph := NewGraph()
	for _, r := range routes {
		routeGraph.Connect(r.Origin, r.Destination, r.Cost)
	}
	return routeGraph.ShortestPath(origin, destination)
}

// Connect makes a connection between origin and destination with the weigth
func (g *Graph) Connect(origin string, destination string, weigth float32) {
	originNode, found := g.nodes[origin]
	if !found {
		originNode = newNode(origin)
		g.nodes[origin] = originNode
	}

	destinationNode, found := g.nodes[destination]
	if !found {
		destinationNode = newNode(destination)
		g.nodes[destination] = destinationNode
	}

	originNode.connect(destinationNode, weigth)
}

// ShortestPath finds the shortest Path from origin to destination
func (g *Graph) ShortestPath(origin string, destination string) ([]string, float32) {
	// Invalid input
	originNode, found := g.nodes[origin]
	if !found {
		return make([]string, 0), 0
	}
	_, found = g.nodes[destination]
	if !found {
		return make([]string, 0), 0
	}

	// Initializes control tables
	nodeCost := make(map[string]float32)
	nodeBestOrig := make(map[string]string)
	toVisit := list.New()

	// Fisrt pass over connections
	nodeCost[origin] = 0
	for label, connection := range originNode.connections {
		toVisit.PushBack(connection.destination)
		nodeCost[label] = connection.weight
		nodeBestOrig[label] = originNode.label
	}

	// Main loop
	for n := toVisit.Front(); n != nil; {
		visitCost, _ := nodeCost[n.Value.(*Node).label]
		visitLabel := n.Value.(*Node).label
		for label, connection := range n.Value.(*Node).connections {
			currCost, found := nodeCost[label]
			// New or better connection
			if !found || (visitCost+connection.weight) < currCost {
				// Needs a first or extra visit
				toVisit.PushBack(connection.destination)
				// Update cost and best route
				nodeCost[label] = visitCost + connection.weight
				nodeBestOrig[label] = visitLabel
			}
		}
		// Advance and remove visited item from list
		oldN := n
		n = n.Next()
		toVisit.Remove(oldN)
	}

	// Reverse best route from destination
	BestOrigin, found := nodeBestOrig[destination]
	if !found {
		// No route to destination
		return make([]string, 0), 0
	}

	route := []string{destination}
	for BestOrigin != origin {
		route = append([]string{BestOrigin}, route...)
		BestOrigin, _ = nodeBestOrig[BestOrigin]
	}
	route = append([]string{BestOrigin}, route...)

	return route, nodeCost[destination]
}
