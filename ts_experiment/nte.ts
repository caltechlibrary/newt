/**
 * newtte.ts provides a simple templating engine service based on the Handlebarjs template langauge.
 */
import {
  appInfo,
  OptionsProcessor,
  defaultHttpPort,
  fmtHelp,
  errorResponse,
  loadAndMergeNewtYAML,
  renderTemplateObject,
  ASTInterface,
} from "./deps.ts";

/**
 * helpText assembles the help information for COLD UI.
 *
 * @param {[k: string]: string} helpOpt holds the help options defined for the app.
 */
function helpText(helpOpt: { [k: string]: string }): string {
  const txt: string[] = [
    `%{app_name}(1) user manual | {version} {release_date}
% R. S. Doiel
% {release_date} {release_hash}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] NEWT_YAML_FILE

# DESCRIPTION

{app_name} provides a handlebarjs template engine as a web service.
{app_name}'s YAML configuration file holds a mapping of URL paths
to handlebar templates. If you make an HTTP GET request the service
will return the unrendered template associated with the URL path.
If you make an HTTP POST it takes the POST the data provided will
be used as the template content. Normally the data you provide is in
the form of JSON. You should set the HTTP header "content-type"
to "application/json" as part of your POST request.

# OPTIONS

`,
  ];
  for (let attr in helpOpt) {
    const msg = helpOpt[attr];
    txt.push(`${attr}
: ${msg}
`);
  }
  txt.push(`
# EXAMPLE

Running {app_name} using the newt YAML file named "app.yaml".

~~~shell
{app_name} app.yaml
~~~

`);
  return txt.join("\n");
}

/**
 * NewtHandlebarsHandler is a function forwards takes POSTS of JSON document and processes with the related
 * related URL's templates path names.
 *
 * @param {Request} req holds the http request recieved from the http server
 * @param {debug: boolean} options holds program options that are made available
 * @returns {Response}
 *
 * @example
 * ```
 *   const options = {
 *      debug: true,
 *   };
 *
 *   const server = Deno.serve({
 *     hostname: "localhost",
 *     port: options.port,
 *   }, (req: Request) => {
 *      return NewtHandlebarsHandler(req, options);
 *   });
 * ```
 */
export async function NewtHandlebarsHandler(
  req: Request,
  config: ASTInterface, //{ [k: string]: any },
): Promise<Response> {
  const content_type = req.headers.get("content-type");

  if (req.method == "POST" && content_type == "application/json") {
    return handleTemplateRenderRequest(req, config);
  }
  //FIXME: display template if "GET"
  if (content_type !== "application/json") {
    return errorResponse(req, 415, `media type - ${content_type} unsupported `);
  }
  return errorResponse(
    req,
    405,
    `${req.method}media type, ${content_type}, method not allowed`,
  );
}

export async function handleTemplateRenderRequest(
  req: Request,
  config: ASTInterface,
): Promise<Response> {
  return await renderTemplateObject(config, req);
}

/**
 * main provide the main program entry point. It handle processing command line
 * options and the environment and of launching the Newt Handlebars template service.
 */
async function main() {
  const appName = appInfo.appName;
  const version = appInfo.verion;
  const releaseDate = appInfo.releaseDate;
  const releaseHash = appInfo.releaseHash;
  const licenseText = appInfo.licenseText;
  const op: OptionsProcessor = new OptionsProcessor();
  const defaultPort: number = 0;

  op.booleanVar("help", false, "display help");
  op.booleanVar("license", false, "display license");
  op.booleanVar("version", false, "display version");
  op.booleanVar("debug", false, "turn on debug logging");
  op.numberVar(
    "port",
    defaultPort,
    `set the port number, default ${defaultPort}`,
  );

  op.parse(Deno.args);

  const options = op.options;
  const args = op.args;

  if (options.help) {
    console.log(fmtHelp(helpText(op.help), appInfo));
    Deno.exit(0);
  }
  if (options.license) {
    console.log(licenseText);
    Deno.exit(0);
  }
  if (options.version) {
    console.log(`${appName} ${version}(${releaseDate}: ${releaseHash})`);
    Deno.exit(0);
  }
  const yaml_name = args.shift();
  if (yaml_name === undefined || yaml_name === "") {
    console.log(`${appName} missing NEWT_YAML_FILENAME`);
    Deno.exit(1);
  }
  console.log(`DEBUG loadAndMergeNewtYAML(${yaml_name}, options)`);
  const config = await loadAndMergeNewtYAML(yaml_name, options);

  console.log(
    `Starting ${appName} HTTP service at http://localhost:${config.applications.template_engine.port}`,
  );
  const server = Deno.serve(
    {
      hostname: "localhost",
      port: config.applications.template_engine.port,
    },
    (req: Request): Promise<Response> => {
      return NewtHandlebarsHandler(req, config);
    },
  );
}

// Learn more at https://deno.land/manual/examples/module_metadata#concepts
if (import.meta.main) {
  await main();
}
