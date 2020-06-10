package controller

import (
	"TravelRoute/dal"
	"TravelRoute/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestStartStopServer(t *testing.T) {
	srv := StartWebServer(dal.NewDB(), 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	StopWebServer(srv)
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
	routeDB := dal.NewDB()

	routeDB.InsertRoute(*dal.New("GRU", "BRC", 10))
	routeDB.InsertRoute(*dal.New("BRC", "SCL", 5))
	routeDB.InsertRoute(*dal.New("GRU", "CDG", 75))

	srv := StartWebServer(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	expect := `[{"Origin":"GRU","Destination":"BRC","Cost":10},{"Origin":"BRC","Destination":"SCL","Cost":5},{"Origin":"GRU","Destination":"CDG","Cost":75}]`
	ret := getRoutes(t)
	if ret != expect {
		t.Errorf("Get expected %v, got %v", expect, ret)
	}

	StopWebServer(srv)
}

func TestGetEmptyRoutes(t *testing.T) {
	routeDB := dal.NewDB()

	srv := StartWebServer(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	expect := `[]`
	ret := getRoutes(t)
	if ret != expect {
		t.Errorf("Get expected %v, got %v", expect, ret)
	}

	StopWebServer(srv)
}

func addRoute(t *testing.T, r dal.Route) {
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
	routeDB := dal.NewDB()

	srv := StartWebServer(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	addRoute(t, *dal.New("GRU", "BRC", 10))
	addRoute(t, *dal.New("BRC", "SCL", 5))
	addRoute(t, *dal.New("GRU", "CDG", 75))

	expect := `[{"Origin":"GRU","Destination":"BRC","Cost":10},{"Origin":"BRC","Destination":"SCL","Cost":5},{"Origin":"GRU","Destination":"CDG","Cost":75}]`
	ret := getRoutes(t)
	if ret != expect {
		t.Errorf("Get expected %v, got %v", expect, ret)
	}

	StopWebServer(srv)
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
	routeDB := dal.NewDB()

	srv := StartWebServer(routeDB, 8080)
	if srv == nil {
		t.Errorf("TravelServer expected not nil, got nil")
	}

	var tests = []struct {
		origin      string
		destination string
	}{
		{"GRU", "CDG"},
	}

	addRoute(t, *dal.New("GRU", "BRC", 10))
	addRoute(t, *dal.New("BRC", "SCL", 5))
	addRoute(t, *dal.New("GRU", "CDG", 75))
	addRoute(t, *dal.New("GRU", "SCL", 20))
	addRoute(t, *dal.New("GRU", "ORL", 56))
	addRoute(t, *dal.New("ORL", "CDG", 5))
	addRoute(t, *dal.New("SCL", "ORL", 20))

	for _, test := range tests {
		expectedBestRoute, expectedCost := domain.FindCheapestRoute(routeDB.GetRoutes(), test.origin, test.destination)
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

	StopWebServer(srv)
}
