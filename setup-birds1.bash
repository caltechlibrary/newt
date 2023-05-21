#!/bin/bash

#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
mkdir -p birds1/htdocs

# Generate the empty files we'll use in the demo.
touch birds1/birds.csv
touch birds1/build.bash
touch birds1/htdocs/index.html

# Generate a README
cat <<EOT>birds1/README.md

# Birds 1, a demo with CSV file and Pandoc 3

This directory holds our demo.

EOT

# Generate some test data to load into our models
cat <<EOT>birds1/birds.csv
bird_name,place,sighted
robin, seen in my backyard,2023-04-16
humming bird, seen in my backyard, 2023-02-28
blue jay, seen on my back porch, 2023-01-12
EOT

# Generate page.tmpl
cat <<EOT>birds1/page.tmpl
<DOCTYPE html lang="en">
<html>
  <head><title>\$title\$</title></head>
  <body>
    <header>\$title\$</header>
	<p>
    <h1>Welcome to the bird list!</h1>
    <p>
\$body\$
	<p>
	<footer>Updated: \$date\$</footer>
  </body>
</html>
EOT

# Generate a simple build.sh script to run Pandoc
D=$(date)
cat <<EOT >birds1/build.sh
#!/bin/sh
echo 'Running Pandoc, converting CSV to Markdown'
pandoc --from csv --to html \
       --metadata title="Birds 1 Demo" \
	   --metadata date="$D" \
	   --template page.tmpl \
	   birds.csv >htdocs/index.html
EOT
chmod 775 birds1/build.sh
START=$(pwd)
cd birds1 && ./build.sh 
cd "${START}"
