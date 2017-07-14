# GoSoundTransit
A simple go program that looks that monitors buses at stops. It uses notifications when it finds buses that match the criteria. It monitors for 1 hour starting at a given time of the day.

## How to install it

```bash 
go build && mv gosoundtransit GoSoundTransit.app/Contents/MacOS/
```

## How to run it
```bash
make install-launchdfile
```
## Example of a configuration file
```json
{
    "pudgetSoundApiKey":"",
    "stopsToMonitor":[
        {
            "stopId":"1_682",
            "name": "4th and University",
            "routes":["1_100190", "1_100270"]
        },
        {
            "stopId":"1_1190",
            "name": "6th and Pike",
            "routes":["1_100190", "1_100270"]
        }

    ],
    "frequencyToMonitor": 5,
    "monitorDuration": 60,
    "startMonitoringHour": 21,
    "startMonitoringMinute": 0
}```




