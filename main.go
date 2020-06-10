package main

import (
	"TravelRoute/graph"
	"TravelRoute/route"
	"TravelRoute/web"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func buildRoutesDB() *route.DB {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}

	routesDB := route.NewDB()
	parser := route.NewCSVParser(routesDB)
	parser.ParseStream(file)

	fmt.Println("Routes added:")
	for _, route := range routesDB.GetRoutes() {
		fmt.Println(route)
	}
	return routesDB
}

func readInput(scanner *bufio.Scanner) (string, bool) {
	scanner.Scan()
	if scanner.Text() == "q" {
		return "", true
	}
	return scanner.Text(), false
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: TravelRoute FILE.csv\n\tPress 'q' to exit")
		os.Exit(1)
		return
	}

	routesDB := buildRoutesDB()
	srv := web.Start(routesDB, 8080)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Please enter the route origin:")
		origin, exit := readInput(scanner)
		if exit {
			break
		}

		fmt.Println("Please enter the route destination:")
		destination, exit := readInput(scanner)
		if exit {
			break
		}

		fmt.Println("Calculating best route...")
		bestRoute, cost := graph.FindCheapestRoute(routesDB.GetRoutes(), origin, destination)
		if len(bestRoute) != 0 {
			fmt.Printf("Best route: %v > $%v\n", strings.Join(bestRoute, " - "), cost)
		} else {
			fmt.Println("No route found!")
		}
	}

	web.Stop(srv)
}
