package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/newt"

	// 3rd Party packages
	"gopkg.in/yaml.v3"
)

var (
	helpText = `---
title: "{app_name}(1) user manual | {version} {release_hash}"
pubDate: {release_date}
author: "R. S. Doiel"
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] TEMPLATE_NAME [JSON_FNAME]

# DESCRIPTION

**{app_name}** is command line template rendering engine. It is provided so
you can debug your mustache templates without resorting to using cURL. It uses
the same template library that Newt's Mustache uses so the results should
match. **{app_name}}** can read JSON data from a file or from standard 
input. The template name needs to be provided on the command line.
In this way **{app_name}** can be used to process JSON in a POSIX style
pipe line.

# OPTIONS

The following options are supported by **{app_name}**.

-h
: display this help message

-license
: display the software license

-version
: display version information

-page
: points at a YAML file that will be used as elements in the Mustache template.

# EXAMPLE

In this example there is a JSON file called "data.json" and a template called "page.tmpl"
and **{app_name}** is used to run the JSON data through the template. In the first
example the data file is specified as part of the command line arguments in the
second it is read from standard input via file redirection. The Third version
works the same way but data is from a pipe.

~~~shell
{app_name} page.tmpl data.json
{app_name} page.tmpl <data.json
cat data.json | {app_name} page.tmpl
~~~

`

	showHelp    bool
	showLicense bool
	showVersion bool
)

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following variables are set when version.go is generated
	version := newt.Version
	releaseDate := newt.ReleaseDate
	releaseHash := newt.ReleaseHash
	fmtHelp := newt.FmtHelp

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	
	// App option(s)
	fNames := ""
	flag.StringVar(&fNames, "page", "", "YAML file containing elements to pass to the Mustache template")

	// We're ready to process args
	flag.Parse()
	args := flag.Args()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", newt.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}
	pageElements := map[string]interface{}{}
	if fNames != "" {
		src, err := os.ReadFile(fNames)
		if err != nil {
			fmt.Fprintf(eout, "failed to read %s, %s\n", fNames, err)
			os.Exit(1)
		}
		if err = yaml.Unmarshal(src, &pageElements); err != nil {
			fmt.Fprintf(eout, "failed to decode %s, %s\n", fNames, err)
			os.Exit(1)
		}
	}
	os.Exit(newt.RunMustacheCLI(in, out, eout, args, pageElements))
}
