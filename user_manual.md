
# Newt User Manual

## What is in the Newt toolbox?

- [Newt](newt.1.md), is a development tool. It's primary task is as a code generator and secondary task is a app runner. As an app runner it is capabile of running Newt's router and template engine and selected JSON data sources such as datasetd or PostgREST.
- [Newt's router](newtrouter.1.md) is a web service that manages your requests and data pipelines as well as static file hosting
- [Newt's template engine](newthandlebars.1.md) is a web service providing a light weight template engine supporting the [Handlebars](https://handlebarsjs.com) template language.

## Newt YAML syntax (draft)

- [Newt YAML Syntax](newt_yaml_syntax.md)

## Tutorials and reference materials (drafts)

- Getting Started with Newt
    - [Newt Overview](newt_command_overview.md)
    - [Configuring](newt_config_explained.md)
    - [Modeling Data](newt_model_explained.md)
    - [Checking your YAML](newt_check_explained.md)
    - [Generating code](newt_generate_explained.md)
- Running your Newt App
    - [Newt Router Explained](newtrouter_explained.md)
    - [Newt Handlebars Explained](newthandlebars_explained.md)
- [Have more questions?](more_questions.md)
- [Reference Material](reference_material.md), links to prior art

## Presentations

- [Newt Prototype 3](presentation3/), (covers third prototype, refining applications and template switch from mustache to handlebars)
- [Newt Prototype 2](presentation2/), SoCal Code4Lib Meetup at USC, April 19, 2024 (covers second prototype)
- Original [Newt Presentation](presentation/), SoCal Code4Lib Meetup at UCLA, July 14, 2023 (covers first prototype)

## Installation

- [Installing Newt](INSTALL.md)
- [Installing Dataset and Datasetd](https://caltechlibrary.github.io/dataset/INSTALL.md)
- [Installing PostgREST](INSTALL-PostgREST.md) from source (developer notes)
- [Installing Pandoc](INSTALL-Pandoc.md) from source (developer notes, Pandoc is needed to compile Newt)

## Reference

- [Evolving Assumptions](evolving_assumptions.md), a timeline of web related technological changes
- [Structured Data Representations](structured_data_representations.md), a musing on how to handle complex data objects in HTML, CSS and JavaScript
