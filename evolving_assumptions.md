
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

Web development seen in this abbreviated timeline shows growth complexity of implement and expectations. In 2024 I think many people and organizations operating on the maxim, "the web is complex so my tools need to be complex and what I create will be necessarily complex". The maxim is not sustainable. 

The seeds of simplification are out there. They've cropped up at almost every turning point in the web evolution. Sometimes the seeds of simplification can actually result in more complexity, e.g. my experience with early NodeJS work was liberating. I went from projects where I dealt with give or six programming and configuration languages to four language (e.g. JavaScript, HTML, CSS and SQL). Yet the NodeJS/NPM combo also ushered in the assumption of a complex ecosystem to build your "web application" from.  Static website deployments came back into vogue (in part due to cost advantages of using S3 buckets) but we also got increasing numbers of web services, API, and "single page apps" which are just a web site with navigation taken out of the hands of the URL. 

Importantly I think there is beginning to be a sense of a course change. Some of this is technical, some political, some social and importantly economic. Today the web and the "cloud" is a landlord's market. Just as many countries are facing housing crisis for various reasons (e.g. war, wealth gaps, climate change) the web experience is facing one now. In North America organizations are faced with renting resources in data centers just as they faced renting the communications lines in the past. The big change is that you no longer have the option of ownership in the sense that your software isn't under your control. It's under the control of the landlord and their interest and how they want to rent you out.

We need to focus on simplification at many levels. Simplification of our starts (or at least how we conceptualize them), simplification in the code we write, and simplification of our expectations about web we want from the web.

[^2]: E.g. Prototype 2005, Mootools, YUI and jQuery 2006.

[^3]: Ryan Dahl releases introduces NodeJS. This quickly replaces Rhino, Narwhal and Jaxer JavaScript server side efforts

