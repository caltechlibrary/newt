#!/bin/bash

for I in 1 2 3; do
    pandoc -s --template=../page.tmpl --from markdown --to html5 \
	    --metadata title="Birds ${I} Demo" \
        birds${I}/README.md -o README-birds${I}.html
done
