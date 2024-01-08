#!/bin/bash
#

# Launch our Mustache engine
../bin/newtmustache -verbose -port 3032 &
PID=$!
sleep 5
if ! curl --verbose \
		-H 'Content-Type: application/json' \
		-d @mustache_data.json \
		http://localhost:3032; then
	echo "Failed to connect with curl"
	kill "$PID"
	exit 1
fi
kill "$PID"
echo "Sucess!"
