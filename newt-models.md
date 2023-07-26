
# Modeling structured data in Newt

The Newt YAML file includes support to model simple data structures for using with Postgres+PostgREST.  

... FIXME: describe how models work and are declared

## Picking the right variable names

One challenge in modeling data in the Newt YAML file is picking good variable names. If you use common single words "begin", "end", "group", "update", "delete", "select" they may collide with the SQL reserved worlds in Postgres. One way to avoid this is to compose variable names from words separated by an underscore. If you used a short prefix for a project's name this solve the problem. If my project code prefix was "x", the variable names "begin", "end", "group", "update", "delete" could be expressed as "x_begin", "x_end", "x_group", "x_update", "x_delete".
