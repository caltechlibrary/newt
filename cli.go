/**
 * cli.go an implements runners for the cli of the Newt Project.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"bufio"
	//	"bytes"
	//	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
	// codes and adopt those.
	OK = iota
	CONFIG

	// General failure of a command or service
	INIT_FAIL
	CHECK_FAIL
	MODELER_FAIL
	GENERATOR_FAIL
	BUILD_FAIL
	ROUTER_FAIL
	TEMPLATE_ENGINE_FAIL
	NEWT_FAIL
	SWS_FAIL
	POSTGREST_FAIL

	// Internal service failures
	RESOLVE
	HANDLER
	SERVER_ERROR
	UNSUPPORTED_ACTION
	DATA_ERROR
	READ_ERROR
	DECODE_ERROR
	TEMPLATE_ERROR

	// Default service settings
	ROUTER_PORT                  = 8010
	TEMPLATE_ENGINE_PORT         = 8011
	TEMPLATE_ENGINE_TIMEOUT      = 3 * time.Second
	TEMPLATE_ENGINE_BASE_DIR     = "views"
	TEMPLATE_ENGINE_PARTIALS_DIR = "partials"
	TEMPLATE_ENGINE_EXT_NAME     = ".hbs"
	SWS_PORT                     = 8000
	SWS_HTDOCS                   = "."
	POSTGREST_PORT               = 3000
	POSTGRES_PORT                = 5432
)

// backupFile takes a filename, copies it to filename plus ".bak"
func backupFile(appFName string) error {
	buf, err := os.ReadFile(appFName)
	if err != nil {
		return fmt.Errorf("failed to back up %s, aborting write, %s", appFName, err)
	}
	if err := os.WriteFile(appFName+".bak", buf, 0666); err != nil {
		return fmt.Errorf("failed to write back up %s, aborting write, %s\n", appFName, err)
	}
	return nil
}

// getNewtYamlFName - figure out what the Newt YAML filename should be.
// If no filename is provided use the default "app.yaml".
func getNewtYamlFName(args []string) string {
	fName := ""
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg != "" && !strings.HasPrefix(arg, "-") {
			fName = arg
			break
		}
	}
	if fName == "" {
		fName = "app.yaml"
	}
	return fName
}

// hasArg - review args and see if the use option is in the list. If
// hasArg is found then true is returned if not false.
func hasArg(option string, args []string) bool {
	for _, arg := range args {
		if strings.ToLower(arg) == option {
			return true
		}
	}
	return false
}

func renderTemplate(generator *Generator, tType string, modelID string, action string, fName string) error {
	var err error
	dName := path.Dir(fName)
	if dName != "" && dName != "." {
		if _, err = os.Stat(fName); os.IsNotExist(err) {
			os.MkdirAll(dName, 0775)
		}
	}
	if _, err = os.Stat(fName); err == nil {
		if err = backupFile(fName); err != nil {
			return err
		}
	}
	out, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer out.Close()
	generator.out = out
	if err := generator.Generate(tType, modelID, action); err != nil {
		return err
	}
	return nil
}

// RunBuilder is a runner for take the generating SQL, templates, etc. and
// Generate the validator middleware a well as updating the Postgres+PostgRESTS
// configuration and databases.
func RunBuilder(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	fmt.Fprintf(eout, `Builder not fully implemented yet. 

The generate step should have completed. The remaining steps are not yet implemented.

1 compile the validator middleware using Deno
2. (re)create the necessary database in Postgres
3. Run the SQL files setup.sql and models.sql)

`)
	return BUILD_FAIL
}

// RunGenerator is a runner for generating SQL and templates from our Newt YAML file.
func RunGenerator(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	fName := getNewtYamlFName(args)
	if fName == "" {
		fmt.Fprintf(eout, "missing Newt YAML filename\n")
		return CONFIG
	}

	ast, err := LoadAST(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	generator, err := NewGenerator(ast)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return GENERATOR_FAIL
	}
	//NOTE: I need to generate each of the files needed for Postgres and PostgREST
	for _, fName := range []string{"setup.sql", "models.sql"} {
		action := strings.TrimSuffix(fName, ".sql")
		if err := renderTemplate(generator, "postgres", "", action, fName); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return GENERATOR_FAIL
		}
	}

	fName = "postgrest.conf"
	if err := renderTemplate(generator, "postgrest", "", "", fName); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return GENERATOR_FAIL
	}
	templateEngine := ast.GetService("template_engine")
	if templateEngine != nil {
		//NOTE: For each model generate a set of templates
		for _, modelID := range ast.GetModelIds() {
			// backup and generate {model}_create_form.hbs, {model}_create_response.hbs
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_create_form%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "create_form", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_create_response%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "create_response", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}

			// backup and generate {model}_read.hbs
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_read%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "read", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
			// backup and generate {model}_update_form.hbs, {model}_update_response.hbs
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_update_form%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "update_form", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_update_response%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "update_response", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
			// backup and generate {model}_delete_form.hbs, {model}_delete_response.hbs
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_delete_form%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "delete_form", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_delete_response%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "delete_response", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
			// backup and generate {model}_list.hbs
			fName = path.Join(templateEngine.BaseDir, fmt.Sprintf("%s_list%s", modelID, templateEngine.ExtName))
			if err := renderTemplate(generator, "template", modelID, "list", fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
		}
		// Need to generate Handlebars partials.
		for _, partialType := range []string{"head", "header", "nav", "footer"} {
			fName = path.Join(templateEngine.BaseDir, templateEngine.PartialsDir, partialType+templateEngine.ExtName)
			if err := renderTemplate(generator, "partial_template", "", partialType, fName); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return GENERATOR_FAIL
			}
		}
	} else {
		fmt.Fprintf(eout, "WARNING: templated engine not configured\n")
	}
	return OK
}

// RunTemplateEngine is a runner for a Newt's template engine.
func RunTemplateEngine(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, timeout int, verbose bool) int {
	appName := "Newt Template Engine"
	fName := getNewtYamlFName(args)
	if fName == "" {
		fmt.Fprintf(eout, "missing Newt YAML filename\n")
		return CONFIG
	}
	// Load the Newt YAML syntax file holding the configuration
	// and make sure it conforms.
	ast, err := LoadAST(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if ast.Services == nil {
		fmt.Fprintf(eout, "%s not configed in %s\n", appName, fName)
		return CONFIG
	}
	if ast.Templates == nil {
		fmt.Fprintf(eout, "now templates found in %s\n", fName)
		return CONFIG
	}
	if ast.GetService("template_engine") == nil {
		fmt.Fprintf(eout, "missing template engine configuration in %q", fName)
		return CONFIG
	}
	// Instantiate the specific application with the filename and AST object
	mt, err := NewTemplateEngine(ast)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}

	// If port is not set in the config, set it to the default port.
	if mt.Port == 0 {
		mt.Port = TEMPLATE_ENGINE_PORT
	}
	// See if we have a command line option for port to process.
	if port != 0 {
		mt.Port = port
	}
	if verbose {
		fmt.Fprintf(out, "port set to %d\n", mt.Port)
	}

	if timeout != 0 {
		mt.Timeout = time.Duration(timeout) * time.Second
	}
	if mt.Timeout == 0 {
		mt.Timeout = TEMPLATE_ENGINE_TIMEOUT
	}
	if len(mt.Templates) == 0 {
		fmt.Fprintf(eout, "no templates in configuration\n")
		return CONFIG
	}
	// Create mux for http service
	// Resolve partial templates and build handlers
	for _, tmpl := range mt.Templates {
		if verbose {
			tmpl.Debug = true
		}
		if err := tmpl.ResolveTemplate(); err != nil {
			fmt.Fprintf(eout, "%s failed to resolve, %s\n", tmpl.Template, err)
			return RESOLVE
		}
		if err := tmpl.ResolvePath(); err != nil {
			fmt.Fprintf(eout, "failed to build handler for %q, %s\n", tmpl.Pattern, err)
			return HANDLER
		}
	}
	// Launch web service
	fmt.Printf("starting %s listening on port :%d (press Ctrl-c to exit)\n", appName, mt.Port)
	if err := mt.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return TEMPLATE_ENGINE_FAIL
	}
	return OK
}

// RunRouter is a runner for Newt data router and static file service
func RunRouter(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool, port int, htdocs string, verbose bool) int {
	appName := "Newt Router"
	// You can run Newt Router with just an htdocs directory. If so you don't require a config file.
	var (
		err    error
		router *Router
	)
	fName := getNewtYamlFName(args)
	if fName == "" {
		fmt.Fprintln(eout, "missing Newt YAML filename")
		return CONFIG
	}
	ast, err := LoadAST(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	// Finally Instantiate the router from fName and AST object
	router, err = NewRouter(ast)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if router.Port == 0 {
		router.Port = ROUTER_PORT
	}
	if port != 0 {
		router.Port = port
	}
	if htdocs != "" {
		router.Htdocs = htdocs
	}

	// Are we ready to run service?
	if router.Routes == nil && router.Htdocs == "" {
		fmt.Fprintf(eout, "nether routes or htdocs are set.")
		return CONFIG
	}

	if router.Port == 0 {
		fmt.Fprintf(eout, "port is not set, default is not available\n")
		return CONFIG
	}

	if verbose && router.Routes != nil {
		for _, route := range router.Routes {
			route.Debug = true
		}
	}

	if dryRun {
		fmt.Fprintf(out, "configuration and routes successfully processed\n")
		return OK
	}

	// Launch web services
	fmt.Fprintf(out, "starting %s listening on port :%d (use Ctr-c to exit)\n", appName, router.Port)
	if router.Htdocs != "" {
		dir, err := filepath.Abs(router.Htdocs)
		if err == nil {
			fmt.Fprintf(out, "\tstatic content %s\n", dir)
		} else {
			fmt.Fprintf(out, "\tstatic content %s (warning: %s)\n", router.Htdocs, err)
		}
	}
	if err := router.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return ROUTER_FAIL
	}
	return OK
}

// RunStaticWebServer this provides a localhost for static file content.
func RunStaticWebServer(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, verbose bool) int {
	appName := "Newt Static Web Server"
	if port == 0 {
		port = SWS_PORT
	}
	htdocs := SWS_HTDOCS
	if len(args) > 0 {
		htdocs = args[0]
	}
	fmt.Fprintf(out, "starting %s listening on port :%d (use Ctr-c to exit)\n", appName, port)
	if err := NewtStaticFileServer(port, htdocs, verbose); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return SWS_FAIL
	}
	return OK
}

// NewtRunCheckYAML will load a Newt YAML fiel and make sure it can parse the configuration.
func RunNewtCheckYAML(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	fName := getNewtYamlFName(args)
	if fName == "" {
		fmt.Fprintln(eout, "missing Newt YAML filename")
		return CHECK_FAIL
	}
	ast, err := LoadAST(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s error: %s\n", fName, err)
		return CHECK_FAIL
	}
	if !ast.Check(eout) {
		return CHECK_FAIL
	}
	if verbose {
		postgREST := ast.GetService("postgrest")
		if postgREST != nil {
			fmt.Fprintf(out, "PostgREST configuration is %s\n", postgREST.ConfPath)
			fmt.Fprintf(out, "PostgREST will be run with the command %q\n", strings.Join([]string{
				postgREST.AppPath,
				postgREST.ConfPath,
			}, " "))
		}
		if ast.Models != nil {
			for _, m := range ast.Models {
				fmt.Printf("models %s defined, %d elements\n", m.Id, len(m.Elements))
				if m.Description != "" {
					fmt.Fprintf(out, "\t%s\n\n", m.Description)
				}
			}
		}
		router := ast.GetService("router")
		if router != nil {
			port := ROUTER_PORT
			if router.Port != 0 {
				port = router.Port
			}
			fmt.Fprintf(out, "Newt Router configured, port set to :%d\n", port)
			if router.Htdocs != "" {
				fmt.Fprintf(out, "Static content will be served from %s\n", router.Htdocs)
			}
			if ast.Routes != nil {
				for _, r := range ast.Routes {
					fmt.Fprintf(out, "route %s defined, request path %s, pipeline size %d\n", r.Id, r.Pattern, len(r.Pipeline))
					if r.Description != "" {
						fmt.Fprintf(out, "\t%s\n\n", r.Description)
					}
				}
			}
		}
		templateEngine := ast.GetService("template_engine")
		if templateEngine != nil {
			port := templateEngine.Port
			if port == 0 {
				port = TEMPLATE_ENGINE_PORT
			}
			fmt.Fprintf(out, "Newt template engine configured, port set to :%d\n", port)
			fmt.Fprintf(out, "%d templates are defined\n", len(ast.Templates))
			for _, mt := range ast.Templates {
				tList := []string{
					mt.Template,
				}
				fmt.Fprintf(out, "http://localhost:%d%s points at %s\n", port, mt.Pattern, strings.Join(tList, ", "))
				if mt.Description != "" {
					fmt.Fprintf(out, "\t%s\n\n", mt.Description)
				}
			}
		}
	}
	return OK
}

// RunNewt is a runner that can run Newt's router and template engine plus PostgREST if defined in the Newt YAML file.
func RunNewt(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	appName := path.Base(os.Args[0])
	action := ""
	// Extract our actions from the args
	switch len(args) {
	case 0:
		return NEWT_FAIL
	case 1:
		action, args = args[0], []string{}
	default:
		action, args = args[0], args[1:]
	}

	switch action {
	case "config":
		return RunNewtConfig(in, out, eout, args, verbose)
	case "check":
		return RunNewtCheckYAML(in, out, eout, args, verbose)
	case "model":
		return RunModeler(in, out, eout, args)
	case "generate":
		return RunGenerator(in, out, eout, args)
	case "build":
		if errCode := RunGenerator(in, out, eout, args); errCode != 0 {
			return errCode
		}
		return RunBuilder(in, out, eout, args)
	case "run":
		return RunNewtServices(in, out, eout, args, verbose)
	case "ws":
		return RunStaticWebServer(in, out, eout, args, 0, verbose)
	default:
		fmt.Fprintf(eout, "%s does %q is an unsupported action, see %s -help\n", appName, action, appName)
		return UNSUPPORTED_ACTION
	}
	return OK
}

// RunNewtServices will run the applictions defined in your Newt YAML file.
func RunNewtServices(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	appName := path.Base(os.Args[0])
	// Get the Newt YAML file to run
	fName := getNewtYamlFName(args)
	if fName == "" {
		fmt.Fprintf(eout, "%s expected a Newt YAML filename\n", appName)
		return CONFIG
	}
	ast, err := LoadAST(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s failed to load %q, %s", appName, fName, err)
		return CONFIG
	}
	if ast.Services == nil {
		fmt.Fprintf(eout, "no applications configured in %s", fName)
		return CONFIG
	}
	// Startup PostgREST if configured in the Newt YAML file.
	postgREST := ast.GetService("postgrest")
	if postgREST != nil &&
		postgREST.ConfPath != "" && postgREST.AppPath != "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return POSTGREST_FAIL
		}
		cmd := exec.Command(postgREST.AppPath, postgREST.ConfPath)
		cmd.Dir = cwd
		// Setup the stdin and stdout to be visible from Newt
		cmd.Stdout = out
		cmd.Stderr = eout
		log.Printf("starting %s %s listening on :%d (use Ctrl-c to exit)", postgREST.AppPath, postgREST.ConfPath, postgREST.Port)
		err = cmd.Start()
		if err != nil {
			log.Println(err)
			return POSTGREST_FAIL
		}
		log.Printf("%s running with pid %d in the backround", postgREST.AppPath, cmd.Process.Pid)
		cmd.Process.Release()
	}
	// Setup and start Newt template engine first
	if ast.GetService("template_engine") != nil {
		go func() {
			RunTemplateEngine(in, out, eout, args, 0, 0, verbose)
		}()
	}

	// The router starts up second and is what prevents service from falling through.
	if ast.GetService("router") != nil {
		go func() {
			RunRouter(in, out, eout, args, false, 0, "", verbose)
		}()
	}
	// NOTE: we need to wait for a signal so that our external process and Go routines aan run.

	// Set up channel on which to send signal notifications.
	c := make(chan os.Signal, 1)

	// We're listening for all signals, probably should narrow this down.
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until any signal is received.
	s := <-c
	fmt.Println("exited with signal:", s)
	return OK
}

// RunNewtConfig will initialize a Newt project by creating a Newt YAML file interactively.
func RunNewtConfig(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	var answer string

	ast := &AST{}
	readBuffer := bufio.NewReader(in)
	// Step 1. Figure out what we're going to call our generated Newt YAML file.
	appFName := getNewtYamlFName(args)
	if appFName == "" {
		fmt.Fprintf(eout, "missing Newt YAML Filename\n")
		return INIT_FAIL
	}
	skipPrompts := hasArg("-y", args)
	if skipPrompts {
		answer = "y"
	}
	if !skipPrompts {
		if _, err := os.Stat(appFName); err == nil {
			fmt.Fprintf(out, "Opening %q\n", appFName)
			ast, err = LoadAST(appFName)
		} else if len(args) <= 1 {
			fmt.Fprintf(out, "Creating %q\n", appFName)
		}
	}
	// Step 2. Decide about meta data of your app
	if ast.AppMetadata == nil {
		ast.AppMetadata = new(AppMetadata)
	}
	// Step 3. Decide which services you're going to use (a .Services will need to exist).
	if ast.Services == nil {
		ast.Services = NewServices()
	}
	for {
		//NOTE: Each of these should reflect the current model list in ast.
		setupAppMetadata(ast, readBuffer, out, appFName, skipPrompts)
		setupRouter(ast, readBuffer, out, appFName, skipPrompts)
		setupPostgresAndPostgREST(ast, readBuffer, out, appFName, skipPrompts)
		setupTemplateEngine(ast, readBuffer, out, appFName, skipPrompts)
		/*
		setupEnvironment(ast, readBuffer, out, appFName, skipPrompts)
		setupOptions(ast, readBuffer, out, appFName, skipPrompts)
		*/

		// Now output the YAML
		_, err := ast.Encode()
		if err != nil {
			fmt.Fprintf(eout, "Failed to generate %s, %s\n", appFName, err)
			return INIT_FAIL
		}
		if skipPrompts {
			answer = "y"
		} else {
			fmt.Fprintf(out, "Save and exit (Y/n)? ")
			answer = getAnswer(readBuffer, "y", true)
		}
		if answer == "y" {
			// We're ready to write out result.
			// If file exists make a back up copy
			if err := ast.SaveAs(appFName); err != nil {
				fmt.Fprintf(eout, "failed to write %s, %s\n", appFName, err)
				return INIT_FAIL
			}
			break
		} else {
			fmt.Fprintf(out, "Exit without saving (y/N)? ")
			answer = getAnswer(readBuffer, "n", true)
			if answer == "y" {
				fmt.Fprintf(out, "%s was not saved\n", appFName)
				return INIT_FAIL
			}
		}
	}
	if _, err := os.Stat(".git"); err == nil {
		fmt.Fprintf(out, `It appears your are using the git revision control system.
You should make sure that generated code containing secrets is NOT included in your
repository for %q. It is recommented that add the following to your .gitignore file.

    # Newt Project ignore list.
    *setup*.sql
    postgrest.conf

`, appFName)
	}
	return OK
}

func RunModeler(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	var (
		answer string
	)
	ast := &AST{}
	readBuffer := bufio.NewReader(in)
	// Step 1. Figure out what we're going to call our generated Newt YAML file.
	appFName := getNewtYamlFName(args)
	if appFName == "" {
		fmt.Fprintf(eout, "missing Newt YAML Filename\n")
		return MODELER_FAIL
	}
	if _, err := os.Stat(appFName); err == nil {
		ast, err = LoadAST(appFName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return MODELER_FAIL
		}
	} else {
		fmt.Fprintf(out, "Create %q (Y/n)? ", appFName)
		answer = getAnswer(readBuffer, "y", true)
		if answer != "y" {
			fmt.Fprintf(eout, "aborting creation of %q\n", appFName)
			return MODELER_FAIL
		}
	}
	// Step 2. If "applications" not configured, generate default configuration.
	if ast.Services == nil ||
		ast.GetService("router") == nil ||
		ast.GetService("postgres") == nil ||
		ast.GetService("postgrest") == nil ||
		ast.GetService("template_engine") == nil {
		ast.Services = NewServices()
		skipPrompts := true
		setupAppMetadata(ast, readBuffer, out, appFName, skipPrompts)
		setupRouter(ast, readBuffer, out, appFName, skipPrompts)
		setupPostgresAndPostgREST(ast, readBuffer, out, appFName, skipPrompts)
		setupTemplateEngine(ast, readBuffer, out, appFName, skipPrompts)
		/*
		setupEnvironment(ast, readBuffer, out, appFName, skipPrompts)
		setupOptions(ast, readBuffer, out, appFName, skipPrompts)
		*/
		bName := path.Base(appFName)
		ext := path.Ext(appFName)
		mName := strings.TrimSuffix(bName, ext)
		if len(args) == 0 {
			args = append(args, mName)
		} else if len(args) == 1 {
			args[0] = mName
		}
	}
	// Step 3. build our lists of models and manage them
	if err := modelerTUI(ast, in, out, eout, appFName, args); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return MODELER_FAIL
	}
	return OK
}
