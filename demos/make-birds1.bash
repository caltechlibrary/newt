#!/bin/bash

#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
if [ -d birds1 ]; then
	rm -fR birds1
fi
mkdir -p birds1/htdocs

# Generate a README
cat <<EOT>birds1/README.md

# Birds 1 demo

This demo presents a list of birds. It is implement
as a static site rendered with Pandoc.

Pandoc takes the CSV file, "bids.csv" and converts the
content to an HTML table then drops that in the "body"
of the template called "page.tmpl". A simple shell script
called "build.sh" is provided to save some typing. The
site is rendered to the "htdocs" directory. This can be copied
to your static site host or served out with your web server.

To update the website you edit the CSV file then re-run the
"build.sh" script.

~~~
edit birds.csv
./build.sh
~~~

EOT

# Generate some test data to load into our models
cat <<EOT>birds1/birds.csv
"bird","place","sighted"
"robin","seen in my backyard","2023-04-16"
"humming bird","seen in my backyard","2023-02-28"
"blue jay","seen on my back porch","2023-01-12"
EOT

# Generate page.tmpl
cat <<EOT>birds1/page.tmpl
<DOCTYPE html lang="en">
<html>
  <body>
    <h1>Welcome to the bird list!</h1>
\${body}
  </body>
</html>
EOT

# Generate a simple build.sh script to run Pandoc
cat <<EOT >birds1/build.sh
#!/bin/sh
echo 'Running Pandoc, converting CSV to Markdown'
pandoc --from csv --to html \\
       --metadata title="Birds 1 Demo" \\
	   --template page.tmpl \\
	   birds.csv >htdocs/index.html
EOT
chmod 775 birds1/build.sh
START=$(pwd)
cd birds1 && ./build.sh 
cd "${START}" || exit 1
tree birds1
wc -l birds1/*.* birds1/htdocs/*.*
