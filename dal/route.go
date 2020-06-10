// Package dal implements simple functions to store and retrieve routes.
package dal

import "io"

// Route defines a weighted oriented connection between 2 airports
type Route struct {
	Origin      string
	Destination string
	Cost        float32
}

// NewRoute Constructs a route given an origin destination and cost
func NewRoute(origin string, destination string, cost float32) *Route {
	return &Route{origin, destination, cost}
}

// DB Defines an memory DataBase to store our routes
type DB struct {
	routes []Route
	stream *io.ReadWriter
}

// NewDB constructs a new Route Database
func NewDB(stream io.ReadWriter) *DB {
	db := DB{make([]Route, 0), &stream}
	newCSVParser(&db).parseStream(&stream)
	return &db
}

// InsertRoute inserts a route in the database
func (rDB *DB) InsertRoute(route Route) {
	rDB.routes = append(rDB.routes, route)
	newCSVParser(rDB).writeLastRouteToStream(rDB.stream)
}

// GetRoutes retrieves all routes stored in the Databse
func (rDB *DB) GetRoutes() []Route {
	return rDB.routes
}
