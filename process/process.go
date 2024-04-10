package process

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/geops/gtfsparser"
)

type feature_collection struct {
	Type     string        `json:"type"`
	Features []StopFeature `json:"features"`
}

type StopProperty struct {
	StopId string `json:"stop_id"`
	Count  int    `json:"stop_count"`
}

type geom struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

type StopFeature struct {
	Type       string       `json:"type"`
	Properties StopProperty `json:"properties"`
	Geometry   geom         `json:"geometry"`
}

type StopCounts struct {
	CountMap map[string]int
}

// IncrementCount increments the count for a given stop_id
func (sc *StopCounts) IncrementCount(stopID string) {
	if _, ok := sc.CountMap[stopID]; ok {
		// Increment count if stopID already exists in the map
		sc.CountMap[stopID]++
	} else {
		// Initialize count to 1 if stopID does not exist in the map
		sc.CountMap[stopID] = 1
	}
}

func downloadFile(url string, filePath string) error {
	// Create the file
	output, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

	// Get the data from the URL
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error getting response: %v", err)
	}
	defer response.Body.Close()

	// Check if response is successful
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", response.Status)
	}

	// Write the body to file
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func NewStopCounts() *StopCounts {
	return &StopCounts{
		CountMap: make(map[string]int),
	}
}

func SaveStopsGeojson(fileName string, stops_features []StopFeature) {
	json_file, err := os.Create(fileName + ".json")
	if err != nil {
		fmt.Println("Error creating file")
		return
	}
	encoder := json.NewEncoder(json_file)

	err = encoder.Encode(feature_collection{"FeatureCollection", stops_features})
	if err != nil {
		fmt.Println("issue with encoder")
	}
	fmt.Println("successfully written to file")
}

func Digest() {
	fmt.Println("another one")
	downloadFile("https://www.metrostlouis.org/Transit/google_transit.zip", "gtfs.zip")

	feed := gtfsparser.NewFeed()
	feed.Parse("gtfs.zip")

	fmt.Printf("Done, parsed %d agencies, %d stops, %d routes, %d trips, %d fare attributes\n\n",
		len(feed.Agencies), len(feed.Stops), len(feed.Routes), len(feed.Trips), len(feed.FareAttributes))

	var stops_features []StopFeature
	stopCounts := NewStopCounts()

	now := time.Now()
	for _, trip := range feed.Trips {
		startBeforeNow := now.Before(trip.Service.Start_date.GetTime())
		endAfterNow := now.After(trip.Service.End_date.GetTime())
		if startBeforeNow || endAfterNow {
			continue
		}
		for _, stop := range trip.StopTimes {
			stopCounts.IncrementCount(stop.Stop.Id)
		}
	}

	for _, v := range feed.Stops {
		var array []float32
		array = append(array, v.Lon)
		array = append(array, v.Lat)

		count := stopCounts.CountMap[v.Id]

		feature := StopFeature{
			"Feature",
			StopProperty{v.Id, count},
			geom{"Point", array},
		}
		stops_features = append(stops_features, feature)
	}

	SaveStopsGeojson("stlouis", stops_features)
}
