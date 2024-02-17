
# Preface, the Newt Experiment, February 2024

The Newt Project was started in 2023. The first prototype was implemented and culminated with a [presentation](https://caltechlibrary.github.io/newt/presentation/) to the SoCal Code4Lib group. The lessons learned from the presentation included

- Postgres+PostgREST was exciting (expected)
- SQL was a problem for many people (expected) 
- Pandoc didn't ring a bell with the web developers (a surprise)
- A service that combined only Postgres+PostgREST with Pandoc server wasn't compelling

What was nice about presenting was the excellent and supportive feed back. It is all positive because it set me on a better path. It helped me understand why my colleagues were politely ambivolent to my Newt. I wasn't about the amphibian.

## Summer, Fall 2023

By the time I had wrapped up coding the first prototype I was successfully generating Postgres SQL and PostgREST configuration files. I was troubled by my horrible choices in the YAML syntax I had evolved. You had to know too much about SQL and specifically Postgres. My YAML syntax was a complete mess.

I continued to experiment with Pandoc server. While it is capable and cleverly simple it wasn't particularly fun to work with. If you send invalid JSON you get an error message. This is very helpful. When you sent it a valid JSON POST and Pandoc returns an empty request body it is really hard to see where you went wrong. I found that a completely frustrating experience. I could not ask others to live with that.

My YAML is horrible and Pandoc is out, what next?

## Fall 2023 thru January 2024

I kept purcolating on Newt. Was I building a solution to a problem only I had? I dreaded writing yet another metadata curation app. I had three scheduled to be written and that waa three too many. I was left pondering. Is it me or do web applications really need to be complex?  Surely there is a simpler way. Or di that vanish with the late Prof. Niklaus Wirth?

I needed a clean slate, a new prototype. What should that look like?

## February 2024, a second prototype

My motivating insight came from a cool demo by my colleague [Tommy Keswick](https://library.caltech.edu/about/directory) gave for a project he did with Caltech Archives. He showed the GitHub YAML issue syntax[^100] needed to use GitHub issues to trigger GitHub workflows.  When I saw how the syntax setup the data needed for the actions it struck a chord. That is what I needed for describing my data models. I don't know who designed the GitHub YAML issue template syntax but they did a really good job. It feels like you're talking about the elements in a web form. That means as a web developer I'm building on elements and concepts that are second nature.  I reread the documenetation at [MDN](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types) for HTML5 form element types. The light buld turned on. I could infer directly using the GitHub YAML issue template syntax (abbr: GHYTS) the SQL I needed in Postgres or SQLite3 to hold the data. GHYTS can indeed be used as a modeling language targeting an SQL database.  Of course it is designed to model an HTML form so that problem gets solved too.

My sage colleage [Mike Hucka](https://library.caltech.edu/about/directory) pointed to me the problems of inventing new languages, even new domain spefic languages. When you invent a new language or new DSL you are increasing cognative load on your the developer.  By using GHYTS I side step that problem. Our developers already know how to create the issue temaplates on GitHub. I am repurposing that knowledge for code generation to build a custom application. Time to start thinking through the specifics of the new prototype.

## The second prototype

Newt tools need configuration. YAML is the language used to express that configuration. Newt's second prototype uses CITATION.cff to describe the application mentadata for the Newt based project. Actually all it needs is to point to the CITATION.cff[^101] and the metadata about the application you're creating is complete.  The horrible part of the first YAML for the first prototype was the modeling syntax. GHYTS solves that. What's left to learn is much easier to describe. It's how to route the data and how to provide templates support. But I am getting ahead of my story so I'll leave that for the next section.


[^100]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>, as viewed Feb 2024

[^101]: See <https://citation-file-format.github.io/>

