
# Newt User Manual

## What is in the Newt toolbox?

- [Newt](newt.1.md), is a development tool. It's primary task is to use intepret the YAML modeled data and render code, templates and configuration to run your app. The Newt command can also function as an app runner. As an app runner it is capabile of running Newt's router and template engine and selected JSON data sources such as Postgres+PostgREST.
- [Newt's router](ndr.1.md) (i.e. ndr) is a web service that manages your requests and data pipelines as well as static file hosting
- [Newt's template engine](nte.1.md) (i.e. nte) is a simple stateless template engine as a web service

## Newt YAML syntax (draft)

- [Newt YAML Syntax](newt_yaml_syntax.md)

## Tutorials and reference materials (drafts)

- Getting Started with Newt
  - [Newt Overview](command_overview.md)
  - [Configuring](config_explained.md)
  - [Modeling Data](data_model_explained.md)
  - [Checking your YAML](check_explained.md)
  - [Generating code](generator_explained.md)
- Generating your Newt App
  - [Generated middleware](generated_middleware_explained.md)
  - [Validator Explained](validator_explained.md)
  - [Template Engine Explained](template_engine_explained.md)
- Running your Newt App
  - [Data Router Explained](data_router_explained.md)

## Presentations

- [Newt Prototype 3](presentation3/), (covers third prototype, refining applications and template switch from mustache to handlebars)
- [Newt Prototype 2](presentation2/), SoCal Code4Lib Meetup at USC, April 19, 2024 (covers second prototype)
- Original [Newt Presentation](presentation/), SoCal Code4Lib Meetup at UCLA, July 14, 2023 (covers first prototype)

## Miscellanea

- [prototype 3 ideas](prototype3.md)
- [prototype 2 ideas](prototype2.md)
- [Have more questions?](more_questions.md)
- [Reference Material](reference_material.md), links to prior art
- [Dev notes and ideas](dev_notes_and_ideas.md)

## Installation

- [Installing Newt](INSTALL.md)
- [Installing Dataset and Datasetd](https://caltechlibrary.github.io/dataset/install.html)

### Migth be useful

- [Installing PostgREST](INSTALL-PostgREST.md) from source (developer notes)
- [Installing Pandoc](INSTALL-Pandoc.md) from source (developer notes, Pandoc is needed to compile Newt)

## Reference

- [Evolving Assumptions](evolving_assumptions.md), a timeline of web related technological changes
- [Structured Data Representations](structured_data_representations.md), a musing on how to handle complex data objects in HTML, CSS and JavaScript
