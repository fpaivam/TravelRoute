package route

import "testing"

func TestRouteInsert(t *testing.T) {
	routeDB := NewDB()

	routeDB.InsertRoute(Route{"GRU", "CON", 5.2})
	routes := routeDB.getRoutes()

	if len(routes) != 1 {
		t.Errorf("routeDB.getRoutes expected size %v, got %v", 1, len(routes))
	}

	if routes[0].origin != "GRU" || routes[0].destination != "CON" || routes[0].cost != 5.2 {
		t.Errorf("route expected %v, got %v", Route{"GRU", "CON", 5.2}, routes[0])
	}
}
