// Package route implements simple functions to store and retrieve routes.
package route

// Route defines a weighted oriented connection between 2 airports
type Route struct {
	Origin      string
	Destination string
	Cost        float32
}

// New Constructs a route given an origin destination and cost
func New(origin string, destination string, cost float32) *Route {
	return &Route{origin, destination, cost}
}

// DB Defines an memory DataBase to store our routes
type DB struct {
	routes []Route
}

// NewDB constructs a new Route Database
func NewDB() *DB {
	return &DB{make([]Route, 0)}
}

// InsertRoute inserts an route in the database
func (rDB *DB) InsertRoute(route Route) {
	rDB.routes = append(rDB.routes, route)
}

// GetRoutes retrieves all routes stored in the Databse
func (rDB *DB) GetRoutes() []Route {
	return rDB.routes
}
