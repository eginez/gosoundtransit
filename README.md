# GoSoundTransit
A simple go program that looks that monitors buses at stops. It uses notifications when it finds buses that match the criteria. It monitors for 1 hour starting at a given time of the day.

## How to install it

```bash 
go build && mv gosoundtransit GoSoundTransit.app/Contents/MacOS
```

## How to run it
```bash
SOUND_TRANSIT_KEY=[key] open GoSoundTransit.app
```





