package domain

import (
	"TravelRoute/dal"
	"testing"
)

func TestGraphShortestPath(t *testing.T) {
	routeDB := dal.NewDB()

	routeDB.InsertRoute(*dal.New("GRU", "BRC", 10))
	routeDB.InsertRoute(*dal.New("BRC", "SCL", 5))
	routeDB.InsertRoute(*dal.New("GRU", "CDG", 75))

	var tests = []struct {
		origin        string
		destination   string
		expectedRoute []string
		expectedCost  float32
	}{
		{"GRU", "CDG", []string{"GRU", "CDG"}, float32(75)},
		{"GRU", "BRC", []string{"GRU", "BRC"}, float32(10)},
		{"SCL", "BRC", []string{}, float32(0)},
	}

	for _, test := range tests {
		route, cost := FindCheapestRoute(routeDB.GetRoutes(), test.origin, test.destination)
		for i := 0; i < len(test.expectedRoute); i++ {
			if route[i] != test.expectedRoute[i] {
				t.Errorf("FindCheapestRoute expected route %v, got %v", test.expectedRoute, route)
			}
		}

		if cost != test.expectedCost {
			t.Errorf("FindCheapestRoute expected cost %v, got %v", test.expectedCost, cost)
		}
	}

}
