package route

type Route struct {
	origin      string
	destination string
	cost        float32
}

type RouteDB struct {
	routes []Route
}

func New() RouteDB {
	rDB := RouteDB{}
	return rDB
}

func (rDB *RouteDB) InsertRoute(route Route) {
	rDB.routes = append(rDB.routes, route)
}

func (rDB *RouteDB) getRoutes() []Route {
	return rDB.routes
}
