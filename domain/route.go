package domain

import (
	"TravelRoute/algorithm"
	"TravelRoute/dal"
)

// FindCheapestRoute Finds the shortest (cheapest) route between origin and destination in routes
// Returns the list of node labels and the total cost
// Return an empty slice and 0 in case there is no route
func FindCheapestRoute(routes []dal.Route, origin string, destination string) ([]string, float32) {
	routeGraph := algorithm.NewGraph()
	for _, r := range routes {
		routeGraph.Connect(r.Origin, r.Destination, r.Cost)
	}
	return routeGraph.ShortestPath(origin, destination)
}
