/**
 * newthandlebarjs.ts provides a simple templating engine service using the Handlebarjs template langauge.
 */
import {
  appInfo,
  OptionsProcessor,
  defaultHttpPort,
  fmtHelp,
  errorResponse,
  loadAndMergeNewtYAML,
  renderTemplateObject
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

{app_name} provides a handlebarjs template engine as a web service

NEW_YAML_FILE is the YAML file for your Newt application with the templates
property and runtime configuration.

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

Running {app_name} using the newt YAML file named "my_app.yaml".

~~~shell
newthandlebars my_app.yaml
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
  config: { [k: string]: any },
): Promise<Response> {
  const content_type = req.headers.get("content-type");

  if (req.method == "POST" && content_type == "application/json") {
    return handleTemplateRequest(req, config);
  }
  if (content_type !== "application/json") {
    return errorResponse(req, 415, `media type - ${content_type} unsupported `);
  }
  return errorResponse(
    req,
    405,
    `${req.method}media type, ${content_type}, method not allowed`,
  );
}

export async function handleTemplateRequest(
  req: Request,
  config: { [k: string]: any },
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
  const defaultPort: number = defaultHttpPort;

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
