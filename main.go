package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/shurcooL/trayhost"
)

const endpoint = "http://api.pugetsound.onebusaway.org"

const routeId316 = "1_100190"
const routeId76 = "1_100270"

const starMonitoringHour = 17
const starMonitoringMinute = 0

const stopId4thAndUniv = "1_682"
const stopId6thAndPike = "1_1190"
const monitoringDuration = time.Minute * 60
const frequencyToMonitor = time.Minute * 5

var StopsToMonitor = map[string][]string{
	stopId4thAndUniv: {routeId76, routeId316},
	stopId6thAndPike: {routeId76, routeId316},
}

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
	if err != nil {
		return
	}
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
	notification := trayhost.Notification{
		Title:   "Bus " + a.RouteShortName,
		Body:    a.String(),
		Handler: nil,
	}

	notification.Display()
}

func searchAndNotify(apiKey, stopId, routeId string) {
	arr, _ := getTimesForRouteAtStop(apiKey, stopId, routeId)
	for _, a := range arr {
		notify(a)
	}
}

func makeMenu() (menus []trayhost.MenuItem) {
	menus = []trayhost.MenuItem{
		trayhost.MenuItem{
			Title:   "Quit",
			Enabled: nil,
			Handler: func() { trayhost.Exit() },
		},
	}
	return
}

func startMonitoring(apiKey string) {
	startTime := time.Now()
	for {
		for k, v := range StopsToMonitor {
			for _, r := range v {
				go searchAndNotify(apiKey, k, r)
			}
		}

		if time.Since(startTime) > time.Duration(monitoringDuration) {
			log.Println("Exiting after monitoring for ", monitoringDuration, " minutes")
			return
		} else {
			time.Sleep(frequencyToMonitor)
		}
	}

}

func initApp() {
	var ep string
	var err error
	var imgData []byte

	ep, err = os.Executable()
	if err != nil {
		panic(err)
	}

	imgData, err = ioutil.ReadFile(path.Join(filepath.Dir(ep), "..", "Resources", "tray_icon.png"))
	if err != nil {
		panic(err)
	}
	trayhost.Initialize("GoBus", imgData, makeMenu())

}

// Calculates the time until the first monitoring, if the desired time is within an our
// of the runtime, it will be schedule to run ASAP, otherwise will schedule for the next day
func CalculateDurationUntilFirstTrigger(hourOfDay int, runtime time.Time) time.Duration {
	startTime := time.Date(runtime.Year(), runtime.Month(), runtime.Day(), hourOfDay, 0, 0, 0, runtime.Location())
	// If the happen to be running an hour past our start time
	if startTime.Before(runtime.Add(-1 * time.Hour)) {
		// Set it for the next day
		startTime = startTime.Add(24 * time.Hour)
		return startTime.Sub(runtime)
	}

	if startTime.Before(runtime) {
		return 0 * time.Minute
	}
	return startTime.Sub(runtime)
}

func monitoringLoop(apiKey string, hourOfDay, minuteOfHour int) {
	for {
		nextStart := CalculateDurationUntilFirstTrigger(hourOfDay, time.Now())
		notification := trayhost.Notification{
			Title:   "Start Bus monitoring",
			Body:    "Will start bus monitoring in " + nextStart.String(),
			Timeout: time.Second * 30,
			Handler: nil,
		}
		notification.Display()
		<-time.After(nextStart)
		startMonitoring(apiKey)
	}
}

//Calls arrivals and departures parses the data
//find the route in it and then printout the time
func main() {
	apiKey := os.Getenv("SOUND_TRANSIT_KEY")
	initApp()
	go monitoringLoop(apiKey, starMonitoringHour, starMonitoringMinute)
	trayhost.EnterLoop()
}
