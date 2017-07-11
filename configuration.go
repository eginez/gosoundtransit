package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type StopInformation struct {
	Id     string   `json:"stopId"`
	Name   string   `json:"name"`
	Routes []string `json:"routes"`
}
type GoTransitConf struct {
	ApiKey                string            `json:"pudgetSoundApiKey"`
	StopsToMonitor        []StopInformation `json:"stopsToMonitor"`
	MonitorDuration       int               `json:"monitorDuration"`
	FrequencyToMonitor    int               `json:"frequencyToMonitor"`
	StartMonitoringHour   int               `json:"startMonitoringHour"`
	StartMonitoringMinute int               `json:"startMonitoringMinute"`
}

const confFileName = ".gotransit.json"

func ReadConfiguration() (res GoTransitConf, err error) {
	home := os.Getenv("HOME")
	confFile := path.Join(home, confFileName)
	var f []byte
	if f, err = ioutil.ReadFile(confFile); err != nil {
		return
	} else {
		err = json.Unmarshal(f, &res)
	}
	return
}

func StopIdToName(stops []StopInformation) (res map[string]string) {
	res = make(map[string]string)
	for _, s := range stops {
		res[s.Id] = s.Name
	}
	return
}
