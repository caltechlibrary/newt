#!/bin/bash

for I in 1 2 3; do
    pandoc -s --from markdown --to html5 \
	    --metadata title="Birds ${I} Demo" \
        "birds${I}/README.md" -o "birds${I}/README.html"
done
