#!/bin/sh
echo 'Running Pandoc, converting CSV to Markdown'
pandoc --from csv --to html \
       --metadata title="Birds 1 Demo" \
	   --template page.tmpl \
	   birds.csv >htdocs/index.html
