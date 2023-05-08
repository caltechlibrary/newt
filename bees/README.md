
# Bees

In this demo we are building an application for Bee sightings. This is similar to our previous demo [birds](../birds/) except that we're going to use Pandoc server to provide a template engine for rendering HTML we include in our web browser managed UI.  Note we still require JavaScript in the web brower to perform the requests to our micro services (i.e. PostgREST and Pandoc srever). This approach is not efficient. E.g. the templates are read in as a static file and that along with our JSON API results get fed back into Pandoc server before taking the resulting HTML and inserting it into our web page. But it does demonstrate adding a well documented template engine as a micro service and serves a stepping stone to our next demo, [flowers](../flowers/).
