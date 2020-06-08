package route

import (
	"io"
	"strconv"
	"strings"
)

type CSVParser struct {
	routeDB *DB
}

func NewCSVParser(routeDB *DB) *CSVParser {
	return &CSVParser{routeDB}
}

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
