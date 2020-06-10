package dal

import (
	"io"
	"strconv"
	"strings"
)

// CSVParser defines a Routes CSV Parser
type CSVParser struct {
	routeDB *DB
}

// NewCSVParser constructs a new Routes CSV Parser given a route Database
func NewCSVParser(routeDB *DB) *CSVParser {
	return &CSVParser{routeDB}
}

// ParseStream Parses CSV stream and fills the Route Database
func (csv *CSVParser) ParseStream(reader io.Reader) {
	internalBuffer := make([]byte, 0)

	for {
		temporaryBuffer := make([]byte, 1024)
		bytesRead, err := reader.Read(temporaryBuffer)

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
			csv.routeDB.InsertRoute(route)
		}

		if err == io.EOF {
			break
		}
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

	return New(origin, destination, float32(cost)), false
}
