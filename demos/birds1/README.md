
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

