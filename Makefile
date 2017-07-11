
build:
	go build && cp gosoundtransit GoSoundTransit.app/Contents/MacOS/

run: build
	./GoSoundTransit.app/Contents/MacOS/gosoundtransit

cronfile: build
	echo "30 16 * * * /usr/bin/open ${GOPATH}/src/github.comeginez/gosoundtransit/GoSoundTransit.app" > cronfile

install-cron: cronfile
	crontab cronfile && crontab -l

delete-cron:
	crontab -r


