package web

import (
	"TravelRoute/graph"
	"TravelRoute/route"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestStartStopServer(t *testing.T) {
	srv := Start(route.NewDB(), 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	Stop(srv)
}

func getRoutes(t *testing.T) string {
	resp, err := http.Get("http://localhost:8080/route")
	if err != nil {
		t.Fatalf("http.Get error: %v\n", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error: %v\n", err.Error())
	}

	return string(body)
}

func TestGetRoutes(t *testing.T) {
	routeDB := route.NewDB()

	routeDB.InsertRoute(*route.New("GRU", "BRC", 10))
	routeDB.InsertRoute(*route.New("BRC", "SCL", 5))
	routeDB.InsertRoute(*route.New("GRU", "CDG", 75))

	srv := Start(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	expect := `[{"Origin":"GRU","Destination":"BRC","Cost":10},{"Origin":"BRC","Destination":"SCL","Cost":5},{"Origin":"GRU","Destination":"CDG","Cost":75}]`
	ret := getRoutes(t)
	if ret != expect {
		t.Errorf("Get expected %v, got %v", expect, ret)
	}

	Stop(srv)
}

func TestGetEmptyRoutes(t *testing.T) {
	routeDB := route.NewDB()

	srv := Start(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	expect := `[]`
	ret := getRoutes(t)
	if ret != expect {
		t.Errorf("Get expected %v, got %v", expect, ret)
	}

	Stop(srv)
}

func addRoute(t *testing.T, r route.Route) {
	js, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json.Marshal error: %v\n", err.Error())
	}

	resp, err := http.Post("http://localhost:8080/route", "application/json", bytes.NewBuffer(js))
	if err != nil {
		t.Fatalf("http.Post error: %v\n", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error: %v\n", err.Error())
	}

	ret := string(body)
	if ret != "OK" {
		t.Fatalf("/route response expect %v, got %v\n", "OK", ret)
	}
}

func TestAddRoutes(t *testing.T) {
	routeDB := route.NewDB()

	srv := Start(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	addRoute(t, *route.New("GRU", "BRC", 10))
	addRoute(t, *route.New("BRC", "SCL", 5))
	addRoute(t, *route.New("GRU", "CDG", 75))

	expect := `[{"Origin":"GRU","Destination":"BRC","Cost":10},{"Origin":"BRC","Destination":"SCL","Cost":5},{"Origin":"GRU","Destination":"CDG","Cost":75}]`
	ret := getRoutes(t)
	if ret != expect {
		t.Errorf("Get expected %v, got %v", expect, ret)
	}

	Stop(srv)
}

func getBestRoute(t *testing.T, origin string, destination string) string {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/route/best?Origin=%v&Destination=%v",
		url.QueryEscape(origin), url.QueryEscape(destination)))
	if err != nil {
		t.Fatalf("http.Get error: %v\n", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error: %v\n", err.Error())
	}

	return string(body)
}

func TestBestRoute(t *testing.T) {
	routeDB := route.NewDB()

	srv := Start(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	var tests = []struct {
		origin      string
		destination string
	}{
		{"GRU", "CDG"},
		{"GRU", "BRC"},
		{"BRC", "GRU"},
		{"GRU", "GRU"},
		{"asfd", "CDG"},
		{"BRC", "CDG"},
	}

	addRoute(t, *route.New("GRU", "BRC", 10))
	addRoute(t, *route.New("BRC", "SCL", 5))
	addRoute(t, *route.New("GRU", "CDG", 75))
	addRoute(t, *route.New("GRU", "SCL", 20))
	addRoute(t, *route.New("GRU", "ORL", 56))
	addRoute(t, *route.New("ORL", "CDG", 5))
	addRoute(t, *route.New("SCL", "ORL", 20))

	for _, test := range tests {
		expectedBestRoute, expectedCost := graph.FindCheapestRoute(routeDB.GetRoutes(), test.origin, test.destination)
		resp := bestRouteResponse{Route: expectedBestRoute, Cost: expectedCost}
		expectJS, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("json.Marshal error: %v\n", err.Error())
		}

		bestRoute := getBestRoute(t, test.origin, test.destination)
		if string(expectJS) != bestRoute {
			t.Errorf("BestRoute expected %v, got %v", string(expectJS), bestRoute)
		}
	}

	Stop(srv)
}
