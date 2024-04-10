package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"transit-chat/transit-intensity/process"
)

func main() {
	file, err := os.Open("ca-agencies.csv")
	if err != nil {
		fmt.Println("Issue with opening file")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Issue with records")
		return
	}

	var AllStops []process.StopFeature
	for _, row := range records[1:] {
		fmt.Println(row[0])
		AllStops = append(AllStops, process.Digest(row[0])...)
	}
	process.SaveStopsGeojson("stops", AllStops)
}
