
build:
	go build && cp gosoundtransit GoSoundTransit.app/Contents/MacOS/

run: build
	./GoSoundTransit.app/Contents/MacOS/gosoundtransit

cronfile: build
	echo "30 16 * * * /usr/bin/open ${GOPATH}/src/github.comeginez/gosoundtransit/GoSoundTransit.app" > cronfile

install-cron: cronfile
	crontab cronfile && crontab -l

launchdfile: build
	sed -e s,GOPATH,"${GOPATH}",g cron.plist > cron-gosoundtransit.plist

install-launchdfile: launchdfile
	cp cron-gosoundtransit.plist ~/Library/LaunchAgents
	launchctl load ~/Library/LaunchAgents/cron-gosoundtransit.plist

delete-launchdfile:
	launchctl remove cron_gosoundtransit

delete-cron:
	crontab -r

clean:
	rm cronfile cron-gosoundtransit.plist gosoundtransit
