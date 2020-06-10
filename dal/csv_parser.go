package dal

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// csvParser defines a Routes CSV Parser
type csvParser struct {
	routeDB *DB
}

// newCSVParser constructs a new Routes CSV Parser given a route Database
func newCSVParser(routeDB *DB) *csvParser {
	return &csvParser{routeDB}
}

// parseStream parses CSV stream and fills the Route Database
func (csv *csvParser) parseStream(reader *io.ReadWriter) {
	internalBuffer := make([]byte, 0)

	for {
		temporaryBuffer := make([]byte, 1024)
		bytesRead, err := (*reader).Read(temporaryBuffer)
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}

		newInternalBuffer := make([]byte, len(internalBuffer)+bytesRead)
		copy(newInternalBuffer, internalBuffer)
		copy(newInternalBuffer[len(internalBuffer):], temporaryBuffer[:bytesRead])

		// Adds a end of line caracter to diferenciate from an incomplete line
		if err == io.EOF {
			newInternalBuffer = append(newInternalBuffer, '\n')
		}

		routes, bytesConsumed := processLines(string(newInternalBuffer))
		internalBuffer = newInternalBuffer[bytesConsumed:]

		for _, route := range routes {
			csv.routeDB.routes = append(csv.routeDB.routes, route)
		}

		if err == io.EOF {
			break
		}
	}
}

// writeLastRouteToStream writes in CSV format the last added route to the stream
func (csv *csvParser) writeLastRouteToStream(writer *io.ReadWriter) {
	if len(csv.routeDB.routes) == 0 {
		return
	}

	route := csv.routeDB.routes[len(csv.routeDB.routes)-1]
	_, err := io.WriteString(*writer, toLine(&route))
	if err != nil {
		log.Fatal(err)
	}
}

// splitLines splits the data input into lines.
// Returns an array of lines and the amount of data consumed
func splitLines(data string) ([]string, int) {
	lines := make([]string, 0)
	bytesConsumed := 0

	index := strings.IndexAny(data, "\r\n")
	// No new Line
	for index != -1 {
		line := data[:index]
		if line != "" {
			lines = append(lines, line)
		}
		data = data[index+1:]
		bytesConsumed += index + 1
		index = strings.IndexAny(data, "\r\n")
	}

	return lines, bytesConsumed
}

// processLines splits the input in lines and decode them into Route structs
// Returns an array of Routes and the amount of data consumed
func processLines(data string) ([]Route, int) {
	routes := make([]Route, 0)
	lines, bytesConsumed := splitLines(data)

	for _, line := range lines {
		route, err := processLine(line)
		if err {
			continue
		}
		routes = append(routes, *route)
	}

	return routes, bytesConsumed
}

// toLine transforms the Route Object into a comma separated line
func toLine(route *Route) string {
	if route == nil {
		return ""
	}

	values := make([]string, 3)
	values[0] = route.Origin
	values[1] = route.Destination
	values[2] = fmt.Sprintf("%.2f", route.Cost)

	return strings.Join(values, ",") + "\n"
}

// processLine splits comma separated input and decode it into a Route struct
// Returns a Route pointer and an error flag.
// It will either return nil, true or *Route, false
func processLine(line string) (*Route, bool) {
	values := strings.Split(line, ",")
	if len(values) != 3 {
		return nil, true
	}

	origin := values[0]
	destination := values[1]
	cost, err := strconv.ParseFloat(values[2], 32)
	if err != nil {
		return nil, true
	}

	return NewRoute(origin, destination, float32(cost)), false
}
