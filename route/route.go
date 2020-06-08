package route

type Route struct {
	origin      string
	destination string
	cost        float32
}

func New(origin string, destination string, cost float32) *Route {
	return &Route{origin, destination, cost}
}

type DB struct {
	routes []Route
}

func NewDB() *DB {
	return &DB{}
}

func (rDB *DB) InsertRoute(route Route) {
	rDB.routes = append(rDB.routes, route)
}

func (rDB *DB) getRoutes() []Route {
	return rDB.routes
}
