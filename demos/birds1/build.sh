#!/bin/sh
echo 'Running Pandoc, converting CSV to Markdown'
pandoc --from csv --to html        --metadata title="Birds 1 Demo" 	   --metadata date="Thu Jun 22 09:26:11 PDT 2023" 	   --template page.tmpl 	   birds.csv >htdocs/index.html
