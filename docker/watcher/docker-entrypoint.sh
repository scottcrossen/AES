#!/bin/bash

SLEEPAMNT=2

function finish {
	if [ $PID1 -ge 0 ]; then
		kill $PID1
	fi
	kill $PID2
}

trap 'finish' SIGTERM

function test {
	echo -e "\nTesting Started."
	go test ./...
	echo -e "\nTesting Finished."
}

PID1=-1
while true; do
	if [ $PID1 -ge 0 ]; then
		echo -e "\n[INFO] Restarting test cases in ${SLEEPAMNT} seconds\n"
		kill $PID1
		PID1=-1
	fi
	sleep ${SLEEPAMNT} && glide up && test &
	PID1=$!
	echo -e "\nBuilding release."
	go install main
	mkdir ../rel > /dev/null 2>&1
	cp ../bin/main ../rel/main
	echo -e "\nBuilding Finished."
	inotifywait -e modify -e delete -e create -e attrib ./* --exclude vendor
done &
PID2=$!

wait
