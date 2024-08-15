
# Development Notes and Ideas

## Ideas

### release action

At somepoint in development your Newt app you will want to prepare it for deployment. This could be done with a "release" action where
the Bash/Powershell scripts get generated for deployment as well as additional service files such as those needed for systemd, launchd
or whatever Windows uses to run services.

### Validation

One idea was to push validation back into Postgres+PostsgRESTS, this is tricky as debugging SQL takes a long time.

Idea two was to generate a validation service for the models associated with a new project.  The question is what langauge? TypeScript seems reasonable as the validation code could be shared browser side via TypeScript transpilation to JavaScript. Deno seems the right platform for this.

An older idea I considered was to generate Lua based functions that run in the router.  A route could be flagged to use the validation for a specific model so input would be validated before entering the pipeline.  This has the advantage that Lua is designed to be embedded and I would not require Deno to be available. This in effect would be a pseudo service in the data pipeline. Downside could be performance of the Newt data router depending on implementation details. Another downside is that generating browser side validation, beyond that which HTML 5 provides, is not achieved.

### Pseudo services

There are two use cases where having a pseudo service available before or after the pipeline makes sense.  Validation before starting the data pipeline makes sense. Another case is when creating, updating or deleting objects.  The POST, PUT and DELETE are easily handled by the pipeline but right now I wind up with two templates for each action because I need to handle the care of deliverying the web form and of course handling the response.  It would be nice to simply redirect a successful create or update to the read end point. In the case of create before record creation the identifier isn't known but when the process completes it could be used to redirect the browser to the read page for the newly created object.  Update also could use the same mechanism for redirection even though the identifier is available via the URL path segment.  Finally delete on success could redirect to either the read page with object status (e.g. not found or deleted), the list of objects page or an error page that idendicates the object is no longer available.

## Supporting multiple "back ends"

The initial versions of Newt prototypes assume the back end is Postgres+PostgREST. While robust Postgres and PostgREST are non-trivial applications which places a burden on the novice end user using Newt to generate their application.  Newt should support alternatives. What might those be?

1. Dataset+datasetd provides a nice object collection platform and that itself supports Postgres and SQLite3. This platform doens't currently support object validation but using the YAML notation from Newt it would be easy to implement.  This would encourage a short pipeline with allot of flexibility to add queries and such that datasetd supports.
2. Deno+TypeScript+SQLite3 would be a nice back end. The TypeScript API service could be generated using a very similar API to Dataset+datasetd. This would be easy to evolve over time when the project owner outgrew Newt. Deno+TypeScript can be compiled down to standalone binaries including limited cross compilation. This would open the door to creating a non-Newt hosted application.

If I provide back end choices I should default to one to keep the "model", "generate", "run" sequence simple but allow the other back ends to be configured if needed.

## Code organization

The current code organization needs to be refactor to allow for better integration of future features. One challenging point is managing the routes and template maps based on actions. If the back end changes from Postgres+PostgREST to something else the routes and their pipelines will likely change, additionally it may impact the templates or their mapping.  Right now generating the YAML for routes and template maps is conflated. It would be better to organize that by action which then could use the AST applications proprerty to call the right generation code.

## Questions

Should routes and template maps be generated at model change or at the generation stage?



