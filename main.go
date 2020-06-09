package main

import (
	"TravelRoute/graph"
	"TravelRoute/route"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalln("Usage: TravelRoute FILE.csv")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("could not open file: %v", err)
		return
	}

	routesDB := route.NewDB()
	parser := route.NewCSVParser(routesDB)
	parser.ParseStream(file)

	fmt.Println("Routes added:")
	for _, route := range routesDB.GetRoutes() {
		fmt.Println(route)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Please enter the route origin:")
		scanner.Scan()
		origin := scanner.Text()

		fmt.Println("Please enter the route destination:")
		scanner.Scan()
		destination := scanner.Text()

		fmt.Println("Calculating best route...")
		routeGraph := graph.NewGraph()
		for _, r := range routesDB.GetRoutes() {
			routeGraph.Connect(r.Origin, r.Destination, r.Cost)
		}
		bestRoute, cost := routeGraph.ShortestPath(origin, destination)
		if len(bestRoute) != 0 {
			fmt.Printf("Best route: %v > $%v\n", strings.Join(bestRoute, " - "), cost)
		} else {
			fmt.Println("No route found!")
		}
	}
}
