package controller

import (
	"TravelRoute/dal"
	"TravelRoute/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// TravelServer defines the server
type TravelServer struct {
	srv *http.Server
	wg  *sync.WaitGroup
}

// StartWebServer starts the webserver at the provided port with the provided databse
func StartWebServer(routeDB *dal.DB, port int) *TravelServer {
	srv := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: newWebServer(routeDB)}

	// Used to syncronize Stop call
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// let Stop know we are done
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: " + err.Error())
		}
	}()

	return &TravelServer{srv, wg}
}

// StopWebServer stops the webserver at the provided port
func StopWebServer(ts *TravelServer) {
	if err := ts.srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	// wait for goroutine started in Start() to finish
	ts.wg.Wait()
}

// webServer defines a route's webserver
type webServer struct {
	mux     *http.ServeMux
	routeDB *dal.DB
}

// ServeHTTP uses the default ServerHTTP from http
func (ws *webServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws.mux.ServeHTTP(w, r)
}

func (ws *webServer) routeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		routes := ws.routeDB.GetRoutes()
		js, err := json.Marshal(routes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case http.MethodPost, http.MethodPut:
		var route dal.Route
		err := json.NewDecoder(r.Body).Decode(&route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ws.routeDB.InsertRoute(route)
		fmt.Printf("Route added: %v\n", route)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	default:
		http.Error(w, fmt.Sprintf("%v: Method not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (ws *webServer) bestRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		origin := r.FormValue("Origin")
		if origin == "" {
			http.Error(w, "Missing 'origin' param", http.StatusBadRequest)
			return
		}

		destination := r.FormValue("Destination")
		if destination == "" {
			http.Error(w, "Missing 'Destination' param", http.StatusBadRequest)
			return
		}

		expectedBestRoute, expectedCost := domain.FindCheapestRoute(ws.routeDB.GetRoutes(), origin, destination)
		resp := bestRouteResponse{Route: expectedBestRoute, Cost: expectedCost}
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	default:
		http.Error(w, fmt.Sprintf("%v: Method not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

type bestRouteResponse struct {
	Route []string
	Cost  float32
}

// newWebServer constructs a new Webserver, if no port is provided defaults to 8080
func newWebServer(routeDB *dal.DB) *webServer {
	mux := http.NewServeMux()
	ws := &webServer{mux, routeDB}
	mux.HandleFunc("/route", ws.routeHandler)
	mux.HandleFunc("/route/best", ws.bestRouteHandler)
	return ws
}
