package main

import (
	"encoding/json"
	"fmt"
	"github.com/deckarep/gosx-notifier"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const endpoint = "http://api.pugetsound.onebusaway.org"

const routeId316 = "1_100190"
const routeId76 = "1_100270"

const stopId4thAndUniv = "1_682"
const stopId6thAndPike = "1_1190"

type ArrivalDepartures struct {
	ArrivalEnabled             bool          `json:"arrivalEnabled"`
	BlockTripSequence          int           `json:"blockTripSequence"`
	DepartureEnabled           bool          `json:"departureEnabled"`
	DistanceFromStop           float64       `json:"distanceFromStop"`
	Frequency                  interface{}   `json:"frequency"`
	LastUpdateTime             int64         `json:"lastUpdateTime"`
	NumberOfStopsAway          int           `json:"numberOfStopsAway"`
	Predicted                  bool          `json:"predicted"`
	PredictedArrivalInterval   interface{}   `json:"predictedArrivalInterval"`
	PredictedArrivalTime       int64         `json:"predictedArrivalTime"`
	PredictedDepartureInterval interface{}   `json:"predictedDepartureInterval"`
	PredictedDepartureTime     int64         `json:"predictedDepartureTime"`
	RouteID                    string        `json:"routeId"`
	RouteLongName              string        `json:"routeLongName"`
	RouteShortName             string        `json:"routeShortName"`
	ScheduledArrivalInterval   interface{}   `json:"scheduledArrivalInterval"`
	ScheduledArrivalTime       int64         `json:"scheduledArrivalTime"`
	ScheduledDepartureInterval interface{}   `json:"scheduledDepartureInterval"`
	ScheduledDepartureTime     int64         `json:"scheduledDepartureTime"`
	ServiceDate                int64         `json:"serviceDate"`
	SituationIds               []interface{} `json:"situationIds"`
	Status                     string        `json:"status"`
	StopID                     string        `json:"stopId"`
	StopSequence               int           `json:"stopSequence"`
	TotalStopsInTrip           int           `json:"totalStopsInTrip"`
	TripHeadsign               string        `json:"tripHeadsign"`
	TripID                     string        `json:"tripId"`
}

type ArrivalDepsResponse struct {
	Data struct {
		Entry struct {
			ArrivalsAndDepartures []ArrivalDepartures `json:"arrivalsAndDepartures"`
		} `json:"entry"`
	} `json:"data"`
}

// For more information
// http://developer.onebusaway.org/modules/onebusaway-application-modules/current/index.html
func httpCall(url string) (data []byte, err error) {
	log.Println("HTTP GET to ", url)
	var res *http.Response
	res, err = http.Get(url)
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func buildUrl(key, operation, data string) string {
	return fmt.Sprintf("%s/api/where/%s/%s.json?key=%s", endpoint, operation, data, key)
}

func getTimesForRouteAtStop(apiKey, stopId, routeId string) (res []ArrivalDepartures, e error) {
	url := buildUrl(apiKey, "arrivals-and-departures-for-stop", stopId)
	var r []byte
	if r, e = httpCall(url); e != nil {
		return
	}
	var jsonRes ArrivalDepsResponse
	if e = json.Unmarshal(r, &jsonRes); e != nil {
		return
	}

	res = make([]ArrivalDepartures, 0)
	for _, arrDeps := range jsonRes.Data.Entry.ArrivalsAndDepartures {
		if arrDeps.RouteID == routeId {
			res = append(res, arrDeps)
		}
	}
	return
}

func (a *ArrivalDepartures) String() string {
	// check PredictedArrivalTime first <-- this one has real-time information
	//transform scheduledArrivalTime to date and time (ms since epoch)
	// print print if it is predicted
	//t := time.Unix(a.ScheduledArrivalTime/1000, 0)
	t := time.Unix(a.ScheduledArrivalTime/1000, 0)
	return fmt.Sprintf("For stop %v, bus %v  is comming in %0.0f mins, predicted: %v", a.StopID,
		a.RouteShortName, time.Until(t).Minutes(), a.Predicted)
}

func notify(a ArrivalDepartures) {
	fmt.Println(a.String())
	msg := gosxnotifier.NewNotification(a.String())
	msg.Title = "Bus " + a.RouteShortName
	msg.Group = "com.eginez.go.bus.notifier"
	msg.Push()
}

func searchAndNotify(apiKey, stopId, routeId string) {
	arr, _ := getTimesForRouteAtStop(apiKey, stopId, routeId)
	for _, a := range arr {
		notify(a)
	}
}

//Calls arrivals and departures parses the data
//find the route in it and then printout the time
func main() {
	apiKey := os.Getenv("SOUND_TRANSIT_KEY")
	for {
		go searchAndNotify(apiKey, stopId4thAndUniv, routeId76)
		go searchAndNotify(apiKey, stopId4thAndUniv, "1_100069")
		go searchAndNotify(apiKey, stopId4thAndUniv, routeId316)
		time.Sleep(5 * time.Minute)
	}
}
