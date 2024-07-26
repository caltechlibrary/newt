/**
 * render.ts holds the page handlebarjs rendering functions.
 */
import {
  Handlebars,
  HandlebarsConfig,
} from "./deps.ts";

/**
 * Default uses this config:
 */
const DEFAULT_HANDLEBARS_CONFIG: HandlebarsConfig = {
  baseDir: "views",
  extname: ".hbs",
  partialsDir: "partials/",
  layoutsDir: "layouts/",
  cachePartials: false,
  defaultLayout: "",
  helpers: undefined,
  compilerOptions: undefined,
};

const handle = new Handlebars(DEFAULT_HANDLEBARS_CONFIG);

/**
 * renderPage takes a template path and a page object and returns a Response object.
 *
 * @param {string} template: this name of the template in the views folder
 * @param {{Object} page_object: the page object to apply to template
 * @returns {Promise<Response>} returns a response once everything is ready.
 */
export async function renderPage(
  template: string,
  page_object: { [k: string]: string | object | boolean | undefined },
): Promise<Response> {
  let body: string = await handle.renderView(template, page_object);
  if (body !== undefined) {
    return new Response(body, {
      status: 200,
      headers: { "content-type": "text/html" },
    });
  }
  body =
    `<doctype html>\n<html lang="en">something went wrong, failed to render ${template}.</html>`;
  return new Response(body, {
    status: 501,
    headers: { "content-type": "text/html" },
  });
}

/**
 * renderTemplateObject takes a configuration and req build the page data object and
 * calles renderPage with the results.
 *
 * @param {ASTInterface} config object
 * @param {Request} request object
 * @returns {Promise<Response>} returns a response 
 */
export async function reanderTemplate(config: ASTInterface, req: request): Promise<Response> {
	const template_id = pathIdentifier(req.url);
	if (template_id === undefined) {
		return errorResponse(req, 404, `${template_id} not found`);
	}
	return errorResponse(req, 501, `renderTemplateObject(${JSON.stringify(config)}, ${JSON.stringify(req)}) has not been implemented`);
}

/**
 * makePage takes a template path and a page object and returns a Response object.
 *
 * @param {string} template: this name of the template in the views folder
 * @param {object} page_object: the page object to apply to template
 * @returns {Promise<string>} returns a string once everything is ready.
 */
export async function makePage(
  template: string,
  page_object: { [k: string]: string | object },
): Promise<string> {
  let body = await handle.renderView(template, page_object);
  if (body !== undefined) {
    return body;
  }
  return `<doctype html>\n<html lang="en">something went wrong, failed to render ${template}.</html>`;
}
