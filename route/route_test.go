package route

import "testing"

func TestRouteInsert(t *testing.T) {
	routeDB := NewDB()

	routeDB.InsertRoute(Route{"GRU", "CON", 5.2})
	routes := routeDB.GetRoutes()

	if len(routes) != 1 {
		t.Errorf("routeDB.getRoutes expected size %v, got %v", 1, len(routes))
	}

	if routes[0].Origin != "GRU" || routes[0].Destination != "CON" || routes[0].Cost != 5.2 {
		t.Errorf("route expected %v, got %v", Route{"GRU", "CON", 5.2}, routes[0])
	}
}

func TestEmptyRoute(t *testing.T) {
	routeDB := NewDB()

	routes := routeDB.GetRoutes()

	if len(routes) != 0 {
		t.Errorf("routeDB.getRoutes expected size %v, got %v", 0, len(routes))
	}

	if routes == nil {
		t.Errorf("routes expected %v, got %v", make([]Route, 0), routes)
	}
}
