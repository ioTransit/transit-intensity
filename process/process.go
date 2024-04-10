package process

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/geops/gtfsparser"
)

type StopProperty struct {
	StopId string `json:"stop_id"`
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

func Digest() []StopFeature {
	fmt.Println("another one")
	downloadFile("https://www.metrostlouis.org/Transit/google_transit.zip", "gtfs.zip")

	feed := gtfsparser.NewFeed()
	feed.Parse("gtfs.zip")

	fmt.Printf("Done, parsed %d agencies, %d stops, %d routes, %d trips, %d fare attributes\n\n",
		len(feed.Agencies), len(feed.Stops), len(feed.Routes), len(feed.Trips), len(feed.FareAttributes))

	var stops_features []StopFeature

	for _, v := range feed.Stops {
		var array []float32
		array = append(array, v.Lat)
		array = append(array, v.Lon)

		feature := StopFeature{
			"feature",
			StopProperty{v.Id},
			geom{"point", array},
		}
		stops_features = append(stops_features, feature)
	}
	return stops_features
}
