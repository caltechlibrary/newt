
# Newt a configuration driven router

With the combination of Postgres, PostgREST and Pandoc server we almost have antire web platform. What missing is how to map URL requests to those micro services. This mapping can be formed by Newt.  Newt is designed to work with and like PostgREST and Pandoc server but it's role is treat working with both as a pipe line based on the request received. Newt itself is controlled by Postgres SQL tables that define where to find Pandoc templates and how to map a given URL request to a specific pipe line sequence before returning the results to the front-end web server. Newt only does routing and page assembly. It doens't to user managment or content crontrol. That is left to either Postgres+PostgREST or to the front-end web server (e.g. via Shibboleth or other single sigh-on mechanism).


