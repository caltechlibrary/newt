
# A short history of web development and databases

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. By 1993 the web as an application platform was born[^13] and with it a good platform for providing useful organizational and institutional software.

By the mid 1990s the Open Source databases like MySQL and Postgres were popular choices for building web applications. It is important to note neither MySQL or Postgres spoke HTTP[^14]. To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular Apache web server. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems. By the late 1990s and early 2000s the practice of "mashing up" sites (i.e. content reuse) was the rage. Bespoke systems took advantage of content reuse too. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture"[^15]. Yahoo Pipes inspired Newt's data pipelines.  Eventual the bespoke systems gave way to common use cases[^16]. A good example of a common use case is Apache's [Solr](https://solr.apache.org) search engine. Another example is how bespoke content management systems gave way to [Plone](https://plone.org), [Drupal](https://drupal.org) and [WordPress](https://wordpress.org).

[^13]: Web applications proceeded to eat all the venerable green screen systems they could find. Today's web is mired in surveillance tech and complex solutions. It has drifted far from Sir. Tim's vision of sharing science documents. We need to refocus on the "good ideas" and jettison the complexity that came with the surveillance economy. Newt can be part of that solution. Develop your Newt application with consideration for others.

[^14]: HTTP being the protocol the communicates with. Essentially at the time RDBMS spoke a dialect of SQL as the unifying language. The web of the time understood HTML and to a certain degree XML. By 2000 people were looking for something simpler than XML to move structured data about. [JSON](https://en.wikipedia.org/wiki/JSON) quickly became the answer.

[^15]: The basic concept was to make it easy to work with "data feeds" and combined them into a useful human friendly web pages. It even included a visual programming language to make it friendly to the non-programmer crowd.

[^16]: If a use case is solved reliably enough it becomes "off the shelf" software.

