package main

import (
	"TravelRoute/controller"
	"TravelRoute/dal"
	"TravelRoute/domain"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func buildRoutesDB() *dal.DB {
	file, err := os.OpenFile(os.Args[1], os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}

	routesDB := dal.NewDB(file)
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
	srv := controller.StartWebServer(routesDB, 8080)

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
		bestRoute, cost := domain.FindCheapestRoute(routesDB.GetRoutes(), origin, destination)
		if len(bestRoute) != 0 {
			fmt.Printf("Best route: %v > $%v\n", strings.Join(bestRoute, " - "), cost)
		} else {
			fmt.Println("No route found!")
		}
	}

	controller.StopWebServer(srv)
}
