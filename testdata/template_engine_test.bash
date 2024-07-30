#!/bin/bash
#

# Launch our Template Engine
../bin/nte -verbose test_nte.yaml &
PID=$!
sleep 5
if ! curl --verbose \
		-H 'Content-Type: application/json' \
		-d @template_data.json \
		http://localhost:3032; then
	echo "Failed to connect with curl"
	kill "$PID"
	exit 1
fi
kill "$PID"
echo "Sucess!"
