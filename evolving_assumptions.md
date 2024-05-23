
# A short history of web development and databases

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. By 1993 the web as an application platform was born[^13] and with it a good platform for providing useful organizational and institutional software.

By the mid 1990s the Open Source databases like MySQL and Postgres were popular choices for building web applications. It is important to note neither MySQL or Postgres spoke HTTP[^14]. To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular Apache web server. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems. By the late 1990s and early 2000s the practice of "mashing up" sites (i.e. content reuse) was the rage. Bespoke systems took advantage of content reuse too. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture"[^15]. Yahoo Pipes inspired Newt's data pipelines.  Eventual the bespoke systems gave way to common use cases[^16]. A good example of a common use case is Apache's [Solr](https://solr.apache.org) search engine. Another example is how bespoke content management systems gave way to [Plone](https://plone.org), [Drupal](https://drupal.org) and [WordPress](https://wordpress.org).

[^13]: Web applications proceeded to eat all the venerable green screen systems they could find. Today's web is mired in surveillance tech and complex solutions. It has drifted far from Sir. Tim's vision of sharing science documents. We need to refocus on the "good ideas" and jettison the complexity that came with the surveillance economy. Newt can be part of that solution. Develop your Newt application with consideration for others.

[^14]: HTTP being the protocol the communicates with. Essentially at the time RDBMS spoke a dialect of SQL as the unifying language. The web of the time understood HTML and to a certain degree XML. By 2000 people were looking for something simpler than XML to move structured data about. [JSON](https://en.wikipedia.org/wiki/JSON) quickly became the answer.

[^15]: The basic concept was to make it easy to work with "data feeds" and combined them into a useful human friendly web pages. It even included a visual programming language to make it friendly to the non-programmer crowd.

[^16]: If a use case is solved reliably enough it becomes "off the shelf" software.

## Evolving assumptions

- 1992, implementing a list of authors would required multiple round trips to the web server
  - A form submit would be required to add or remove an author from a list
  - The original HTML inputs were simpler than the current HTML5 iteration (e.g. no date, URL, color types)
  - You synthesized an list behavior by leveraging input element names (e.g. `family_name_1`)
- 1993 NCSA httpd web server introduced
  - CGI, common gateway interface, enable an easier path for dynamic web development
- 1995, Apache Web Server JavaScript is introduced
 - non-standard, initially Netscape only, IE adds slightly incompatible support trying to push embedded Visual Basic
- 1996, CSS introduced, CSS is a DSL for controlling the display/layout of HTML elements
 - Initial browser support problematic
 - Slow convergence of supported features
 - Differences leads to the practice of creating a CSS reset style
- 1998, term LAMP coined, becoming a mainstream web platform 
 - Linux, Apache, MySQL and (Perl, PHP or Python)
- 1999, JavaScript can call home
 - Event Handlers allow form validation, menus, etc.
 - XMLHttpRequest and JavaScript can call home or anyplace else of the net XMLHttpRequest
 - **Building a list of authors doesn't require a page refresh**, feels a little more like a "native application"
- 2000, JSON "discovered" by Douglas Crockford, slow rise replacing XML for transfer data
  - You start to see a divergence between "front end" and "back end" web developers
  - A "web designer" and a "web developer" are not necessarily the same thing
  - User Experience is a big topic eventually leading to a specialist role in web developer down the road
- 2005, 2006 JavaScript libraries and frameworks[^2] rise in adoption
  - Designers start needing programming knowledge
  - Frameworks promise to make building rich web forms and site navigation easier
  - CSS is included as part of JavaScript frameworks, e.g. [YUI](https://en.wikipedia.org/wiki/YUI_Library)
    - CSS reset stylesheets appear
- 2006, Mozilla [XULRunner](https://en.wikipedia.org/wiki/XULRunner) allows building stand alone desktop apps using HTML, CSS and JavaScript, SQLite and XML
- 2007, independent CSS frameworks start to appear (e.g. [Blueprint](https://en.wikipedia.org/wiki/Blueprint_(CSS_framework))
- 2007, iPhone introduced, initially 3rd Party apps are expected to be websites
- 2008/2009, [PhoneGap](https://en.wikipedia.org/wiki/Apache_Cordova) brings web development mobile native platforms
  - You can develop a mobile app using HTML, CSS and JavaScript and target both iPhones and Android devices
- 2009, Server side JavaScript goes mainstream
  - Ryan Dahl introduces NodeJS[^3], demos a [C10K](https://en.wikipedia.org/wiki/C10k_problem) solution in JavaScript using efficient event loops
  - NodeJS supplants [Rhino](https://en.wikipedia.org/wiki/Rhino_(JavaScript_engine)), Narwhal, [Jaxer](https://en.wikipedia.org/wiki/Aptana#Aptana_Jaxer)
  - NodeJS inspires embedded web servers in other programming languages (e.g. Python, PHP, Go)
- 2009, A "web developer" can create back end web services, native Desktop and mobile applications
- 2009, static web site builders emerge and become popular
- 2010, npm building on the prior art of other programming language package systems is introduce
  - This starts a shift in "front end" web developers using GUI/WYSIWYG tools towards tools on the command line
  - Unix shell/command line as a resurgence in popularity
  - "Build systems" become mainstream as part of front end development
- 2014, seeds of Simplification
  - [Tilde club](https://medium.com/message/tilde-club-i-had-a-couple-drinks-and-woke-up-with-1-000-nerds-a8904f0a2ebf)
  - A renaissance for public Unix systems and Gopher?
- 2016 (approaching a decade ago), assumptions start to settle in
  - front end developers are expected to work at the command line
  - they are expected to use complex build systems
  - assumption is they will use a framework for CSS and one for JavaScript
  - NPM, the JavaScript package manager, has a common tool to build on for automating front end packaging and assembly (e.g. Grunt, Gulp, Yeoman)
  - A "front end developer" probably isn't a web designer
- 2019, Gemini Project is initiated, something between Gopher and a text centered web wide web experience

The evolution seen in this abbreviated timeline illustrates a growth complexity, implementation and expectations around web development. In 2024 many organizations and people appear to operating with a maxim, "the web is complex, my tools need to be complex, I create complex things". This maxim is not sustainable in software. 

The seeds of simplification are out there. They've cropped up at almost every turning point in the web evolution. Sometimes the seeds of simplification can actually result in more complexity, e.g. my experience with early NodeJS work was liberating. I went from projects where I worked in five or six programming and configuration languages to four language (e.g. JavaScript, HTML, CSS and SQL). That same liberation also laid the foundation for NodeJS+NPM which ushered in the assumption of a complex ecosystem.  Static website deployments came back into vogue (in part due to cost advantages of using S3 buckets) but some of the static site generators, e.g. Jekyll, seemed to have missed the boat on simplifying things.  These types of simplifications came from developers for developers in many cases. One simplification freeing conceptual bandwidth to facilate an increase in complexity.

On the other hand simplifications which had multiple motivators seem to stick around a while. Static sites are common practice for libraries since they provide a robust platform for distributing information and are also low cost to support (rental of an S3 bucket is cheap for small files). 

You also see other forces encouraging a rethink. The race to the "cloud" has lead to a landlord's market place. For commercial software it is difficult to "buy" software but often forcing us to rent[^4].

I think it is high time to focus our simplification at all levels. Getting to simple in part is an engineering problem, in part a human organizational problem and significantly a "market force" problem. As humans we can take steps to change that if we choose to.

[^2]: E.g. Prototype 2005, Mootools, YUI and jQuery 2006.

[^3]: Ryan Dahl releases introduces NodeJS. This quickly replaces Rhino, Narwhal and Jaxer JavaScript server side efforts

[^4]: Example try buying Windows, macOS, or Adobe Photoshop, we rent the software only. 

