
# Supporting complex form elements, history and brain storm

There are many use cases for for an individual metadata record containing a lists of things when considering applications for galleries, libraries, archives and museums.  Example, a list of authors. A list of authors sounds simple but can be a complex data type. Today it is a complex data type as implemented in repository and catalog systems.

HTML 5 provides a rich set of input types but does not include list element type[^1]. The data model of a web form without enhancement is flat. It is a series of key/value pairs where the data when transmitted is urlencoded and sent as text to the receiving web service. Let's put a pin in that for later.

[^1]: For the this exploration I'm going to table discussing file and select element types.

## Changing base assumptions

In 1992 implementing a list of authors would required multiple round trips to the web server. I submitted the page when you added an author, removed an author or finally submitted the author list. Why so much traffic? The browser only knew the simple native HTML element model. You synthesized an list by imposing a name convention on the input names (e.g. `family_name_1`). HTML only support a few simple element types. In 1995 JavaScript is introduced. In 1999 JavaScript had reach a point where it could be reliably used to smooth things over. It could make adding and removing items in a list on a web page happen without a round trip to the web server.  By 2005 and 2006 JavaScript frameworks start appearing[^2]. These frameworks made building these type of form elements easier. The frameworks competed not just on how easy they were to use but also on richness of the UI elements they implemented. The frameworks also masked over the differences between JavaScript implementations based on the type and version of the web browser being used. Latter we'd called this shimming. In 2009 JavaScript server side and at the command line finally takes hold[^3]. This was significant. It meant that you could use the same programming language on the back end as the front. This simplify form validation by letting your right one set of validation routines rather than two (e.g. once in JavaScript for browser and again in PHP, Python or Perl before sending data to your database). Static web sites started to coming back into fashion too. This I think was influenced by the fact you could create a web server in a few lines of JavaScript, other languages quickly followed that model too.  Fast forward to 2016 and front end developers are expected to work at the command line. They are expected to use complex build systems. The assumption is you're using a framework. The frameworks require you to know the build tools to use them effectively. In 2024 many organizations and individuals still seem to think this is the way it has to be. You can express this approach as the maxim, "the web is complex so the tools to create it must be complex too".

[^2]: E.g. Prototype 2005, Mootools, YUI and jQuery 2006.

[^3]: Ryan Dahl releases introduces NodeJS. This quickly replaces Rhino, Narwhal and Jaxer JavaScript server side efforts

## A radical simplification is due

Today repository and catalog systems rely on these heavy weight frameworks to provide a decent data entry experience. The maxim that evolved from the 2010's is not required. In 2024 we have a modern JavaScript implemented available in all the major web browsers. The browsers are evergreen in the sense they will update themselves as improved versions are released.  This means improvements flow quickly to every network connected computer. This is true even if your "computer" is a tablet, phone or watch.

We have decades of experience with progressive enhancement now. With a couple of heuristics you can avoid jumping through most user interface
development hoops we encountered when the first iPhones were introduced.

1. Design for the smallest screen used to access your web site or application (often that is a phone today)
2. Use CSS to arrange elements on the screen as the screen size expands (e.g. your screen might be a wall sized television)
3. Limit JavaScript for orchestrating behavior (e.g. defining web components) or access interacting with data services
4. Only right the minimum of CSS or JavaScript needed to achieve the desired outcome

Using these four heuristics you can avoid most frameworks be they it JavaScript or even CSS. 

## The challenge poised by web components

The big trouble with web components is JavaScript. JavaScript is used to define them and to run them. If the browser has JavaScript disabled, it is unavailable or it is not supported by the web browser then you're stuck[^4].

[^4]: Lynx, Dillo and NetSurf are browsers that don't support JavaScript. A network disruptions can also prevent JavaScript from reaching your favorite web browser. Filling out a web form in a subway can be a nice place to experience this phenomena.

This is why in spite of evergreen browsers progressive enhancement remains relevant today. In 2005 progressive enhancement meant testing the browser
to see what features were supported and then shimming the browser to support the missing features. A framework like jQuery or YUI was powerful
because it handled all that test and shim work for you. Similarly with CSS it was common to use a "reset" CSS style file to get a CSS baseline
when designing for cross browser compatibility. Today this set of problems is a thing of the past.

Today progressive enhancement more often focuses on site accessibility issues. Can I make my content useful for someone who can't see the screen, hear the recording, track the content or has limited dexterity. The easiest way to get started doing this is to use the native HTML elements. If you need more than that then make sure your web component extends a native element. That way it will inherit those accessibility features. 


My proposal is for simplification in front end development is use our four heuristics in your human interface design. Avoid requiring a complex build system. Build directly with vanilla CSS and vanilla JavaScript writing as little code as you can get away with.  If you need something more than a native HTML element then build a web component but that falls back to a simple web component. The leaves us with the question, which native HTML element to I use as a fallback when JavaScript isn't available?

When I first started using JavaScript to enhance data entry web pages I tended to rely on JSON. First JSON is available in the browser without requiring an extra JavaScript package or include. If we're planning to handle the case where our complex data is being rendering using a simple element then we need to think about how we're representing it. JSON isn't always the best fit.

Here's an example of an person object in JSON.

~~~json
{
  "family_name": "Doiel",
  "lived_name": "Robert",
  "identifers": [
    { 
      "orcid": "0000-0003-0900-6903",
      "clpid": "Doiel-R-S"
    }
  ]
}
~~~

The data is easy enough to read if you're a programmer. People can pick out my name. If they are a researcher or librarian then they might recognize an ORCID id.
The trouble is this notation is visibly busy. The punctuation matters and also make editing it brittle. Dropping a colon, comma or imbalanced
quotes is very common when hand editing JSON. These leave you with invalid JSON.

YAML is another way to represent structured data.  It has become more common in the data science community with the adoption of static site generators where YAML is used to encode "front matter" or document metadata and is also used with rMarkdown and configuration of Jupyter Notebooks. In the YAML spec says JSON is a subset of YAML so the above record actually works as YAML. We can do batter by expressing it more succinctly in YAML.

Here's my person object as YAML.

~~~yaml
fmaily_name: Doiel
lived_name: Robert
idntifiers:
  - orcid: 0000-0003-0900-6903
    clpid: Doiel-R-S
~~~

YAML is expressed in a list oriented manner. Punctuation still counts, e.g. colon and dash mean specific things in YAML. Yet if you need to correct the spelling of my name or ORCID you can do so with less worry because we don't need to worry about quoting. I also think it is easier to read than the JSON. This is especially true if I didn't have the JSON pretty printed.

Let's look at a list of authors since that's the example I am using as a thread through this brain storm.

~~~json
[
  {
    "family_name": "Doiel",
    "lived_name": "Robert",
    "identifers": [
      { "orcid": "0000-0003-0900-6903", "clpid": "Doiel-R-S" }
    ]
  },
  {
    "family_name": "Morrell",
    "lived_name": "Thomas",
    "identifiers": [
       { "orcid": "0000-0001-9266-5146", "clpid": "Morrell-T-E" }
    ]
  }
]
~~~

I think lists of JSON objects definitely too much to typing even if you are a programmer used to working with JSON.

Let's look at the same thing in YAML.

~~~yaml
- family_name: Doiel
  lived_name: Robert
  identifers:
  - clpid: Doiel-R-S
    orcid: 0000-0003-0900-6903
- family_name: Morrell
  lived_name: Thomas
  identifiers:
  - clpid: Morrell-T-E
    orcid: 0000-0001-9266-5146
~~

The above YAML amost looks like a list in Markdown. I need to know the rules about indentation, dashes and colons. If I am including large text blocks I need to know about the pipe character and indentation nuances. Those issues aside this looks more readable and human friendly to me.

Using having our web component extend from textarea means we can support editing even if the web browser used to access our form was Lynx. This does assume that a human an pickup YAML easily.

For web forms that require complex elements like lists I would recommend extending the textarea element and using YAML to encode the structure data. The component implementation can "hide" the unmodified textarea, present a nice UI then update the ocntents of the textarea expressing the structured data in YAML. Likewise when the page loads an pre-populated form the component to read the textarea's content, and decode the YAML to render the UI to manage it.

This approach as in my mind many advantages. I can test the web form processing using curl or other simple HTTP client library of tool. I don't need to learn a framework I only need to learn out to implemenet web components.  The server can handle the textarea content and decode the YAML into an appropriate storage type (e.g. JSON column in a SQL database). All I need for the front end development tools is a text editor, my web browser, a static web server that can process the POST request and display the submitted object.  For backend I just need to be able to parse a standard web form and decode YAML in from those textarea holding the complex data.


