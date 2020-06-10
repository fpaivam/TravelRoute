package dal

import (
	"bytes"
	"testing"
)

func TestProcessLine(t *testing.T) {
	var tests = []struct {
		input    string
		err      bool
		expected Route
	}{
		{"GRU,BRC,10", false, Route{"GRU", "BRC", 10}},
		{"BRC,SCL,5", false, Route{"BRC", "SCL", 5}},
		{"GRU,CDG,75", false, Route{"GRU", "CDG", 75}},
		{"GRU,SCL,20", false, Route{"GRU", "SCL", 20}},
		{"GRU,ORL,56", false, Route{"GRU", "ORL", 56}},
		{"ORL,CDG,5", false, Route{"ORL", "CDG", 5}},
		{"SCL,ORL,20", false, Route{"SCL", "ORL", 20}},
		{"SCL,ORL,20,asdjfh", true, Route{}},
		{"SCL,ORL,", true, Route{}},
		{"sdkfjasdfsdfj", true, Route{}},
	}

	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			route, err := processLine(tt.input)

			if err != tt.err {
				t.Fatalf("route.processLine expected %v, got %v", tt.err, err)
			}

			if err {
				return
			}

			if tt.expected != *route {
				t.Errorf("route expected %v, got %v", tt.expected, route)
			}
		})
	}
}

func TestToLine(t *testing.T) {
	var tests = []struct {
		expected string
		input    *Route
	}{
		{"GRU,BRC,10.00\n", &Route{"GRU", "BRC", 10}},
		{"BRC,SCL,5.00\n", &Route{"BRC", "SCL", 5}},
		{"GRU,CDG,75.00\n", &Route{"GRU", "CDG", 75}},
		{"GRU,SCL,20.00\n", &Route{"GRU", "SCL", 20}},
		{"GRU,ORL,56.00\n", &Route{"GRU", "ORL", 56}},
		{"ORL,CDG,5.00\n", &Route{"ORL", "CDG", 5}},
		{"SCL,ORL,20.00\n", &Route{"SCL", "ORL", 20}},
		{"", nil},
	}

	for _, tt := range tests {
		testname := tt.expected
		t.Run(testname, func(t *testing.T) {
			value := toLine(tt.input)

			if tt.expected != value {
				t.Errorf("value expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func TestSplitLines(t *testing.T) {
	var tests = []struct {
		name          string
		input         string
		bytesConsumed int
		expectedLines []string
	}{
		{"ProvidedInput",
			`GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
`,
			75,
			[]string{
				"GRU,BRC,10",
				"BRC,SCL,5",
				"GRU,CDG,75",
				"GRU,SCL,20",
				"GRU,ORL,56",
				"ORL,CDG,5",
				"SCL,ORL,20"}},
		{"HalfLine",
			`GRU,BRC,10
BRC,SCL,5
GRU,CD`,
			21,
			[]string{
				"GRU,BRC,10",
				"BRC,SCL,5"}},
		{"Empty", "",
			0,
			[]string{}},
		{"MixedLineTerminators", "GRU,BRC,10\r\nBRC,SCL,5\nGRU,CD",
			22,
			[]string{"GRU,BRC,10",
				"BRC,SCL,5"}},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			lines, bytesConsumed := splitLines(tt.input)

			if tt.bytesConsumed != bytesConsumed {
				t.Fatalf("route.splitLines expected %v bytes consumed, got %v", tt.bytesConsumed, bytesConsumed)
			}

			if len(tt.expectedLines) != len(lines) {
				t.Fatalf("route.splitLines expected %v lines, got %v", len(tt.expectedLines), len(lines))
			}

			for i := range lines {
				if tt.expectedLines[i] != lines[i] {
					t.Errorf("route.splitLines expected %v, got %v", tt.expectedLines[i], lines[i])
				}
			}
		})
	}
}

func TestProccessLines(t *testing.T) {
	var tests = []struct {
		name           string
		input          string
		bytesConsumed  int
		expectedRoutes []Route
	}{
		{"ProvidedInput",
			`GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
`,
			75,
			[]Route{
				{"GRU", "BRC", 10},
				{"BRC", "SCL", 5},
				{"GRU", "CDG", 75},
				{"GRU", "SCL", 20},
				{"GRU", "ORL", 56},
				{"ORL", "CDG", 5},
				{"SCL", "ORL", 20}}},
		{"InvalidRoute",
			`GRU,BRC,10
BRC,SCL,5,asjdfa
GRU,CDG,75
`,
			39,
			[]Route{
				{"GRU", "BRC", 10},
				{"GRU", "CDG", 75}}},
		{"Empty", "",
			0,
			[]Route{}},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			routes, bytesConsumed := processLines(tt.input)

			if tt.bytesConsumed != bytesConsumed {
				t.Fatalf("route.processLines expected %v bytes consumed, got %v", tt.bytesConsumed, bytesConsumed)
			}

			if len(tt.expectedRoutes) != len(routes) {
				t.Fatalf("route.processLines expected %v routes, got %v", len(tt.expectedRoutes), len(routes))
			}

			for i := range routes {
				if tt.expectedRoutes[i] != routes[i] {
					t.Errorf("route.processLines expected %v, got %v", tt.expectedRoutes[i], routes[i])
				}
			}
		})
	}
}

func TestParseStream(t *testing.T) {
	var tests = []struct {
		name           string
		input          string
		expectedRoutes []Route
	}{
		{"ProvidedInput",
			`GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20`,
			[]Route{
				{"GRU", "BRC", 10},
				{"BRC", "SCL", 5},
				{"GRU", "CDG", 75},
				{"GRU", "SCL", 20},
				{"GRU", "ORL", 56},
				{"ORL", "CDG", 5},
				{"SCL", "ORL", 20}}},
		{"InvalidRoute",
			`GRU,BRC,10
BRC,SCL,5,asjdfa
GRU,CDG,75`,
			[]Route{
				{"GRU", "BRC", 10},
				{"GRU", "CDG", 75}}},
		{"Empty", "",
			[]Route{}},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			routeDB := NewDB(bytes.NewBufferString(tt.input))
			routes := routeDB.GetRoutes()

			if len(tt.expectedRoutes) != len(routes) {
				t.Fatalf("route.processLines expected %v routes, got %v", len(tt.expectedRoutes), len(routes))
			}

			for i := range routes {
				if tt.expectedRoutes[i] != routes[i] {
					t.Errorf("route.processLines expected %v, got %v", tt.expectedRoutes[i], routes[i])
				}
			}
		})
	}
}

func TestWriteStream(t *testing.T) {
	var buf bytes.Buffer
	routeDB := NewDB(&buf)
	routeDB.InsertRoute(*NewRoute("GRU", "BRC", 10))
	routeDB.InsertRoute(*NewRoute("BRC", "SCL", 5))
	routeDB.InsertRoute(*NewRoute("GRU", "CDG", 75))

	expected := "GRU,BRC,10.00\nBRC,SCL,5.00\nGRU,CDG,75.00\n"
	result := buf.String()
	if expected != result {
		t.Errorf("value expected %v, got %v", expected, result)
	}
}
