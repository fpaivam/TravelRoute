package domain

import (
	"TravelRoute/algorithm"
	"TravelRoute/dal"
)

// FindCheapestRoute Constructs a graph and finds the shortest (cheapest) route
// between origin and destination
func FindCheapestRoute(routes []dal.Route, origin string, destination string) ([]string, float32) {
	routeGraph := algorithm.NewGraph()
	for _, r := range routes {
		routeGraph.Connect(r.Origin, r.Destination, r.Cost)
	}
	return routeGraph.ShortestPath(origin, destination)
}
