/**
 * utils.ts modules holds the method related to handling identifiers and form processing. E.g. validation
 * and extraction from the URL pathname. Turning a form request into an object.
 */

/**
 * errorResponse takes a status code and msg and returns a http response.
 * @param {Request} req holds the request to process for the error response.
 * @param {number} status_code holds the HTTP status code.
 * @param {string} msg holds the message to respond with.
 * @return {Response}
 */
export function errorResponse(
  req: Request,
  status_code: number,
  msg: string,
): Response {
  const txt: string = `HTTP ERROR ${status_code}: ${req.method} ${new URL(req.url).pathname} ${msg}`;
  const body: string = `<html>${txt}</html>`;
  console.log(
    `response ${status_code}, ${req.headers.get("content-type")}, ${txt}`,
  );
  return new Response(body, {
    status: status_code,
    headers: { "content-type": "text/html" },
  });
}


/**
 * hasMethodAndPath gets a request method and url for an extact match of pathname and method
 * @param {Request} req holds the request to evaluatte
 * @param {string} method  holds the target method to match against
 * @param {string} prefix holds the path to compare for an exact match.
 * @returns boolean
 */
export function hasMethodAndPath(
  req: Request,
  method: string,
  prefix: string,
): boolean {
  const pathname = new URL(req.url).pathname;
  if (req.method === method.toUpperCase() && prefix === pathname) {
    return true;
  }
  return false;
}

/**
 * hasMethodAndPrefix gets a request method and url for an extact match
 * of method and a pathname value with mathcing prefix.
 * @param {Request} req holds the request to evaluatte
 * @param {string} method  holds the target method to match against
 * @param {string} prefix holds the path to compare for an exact match.
 * @returns boolean
 */
export function hasMethodAndPrefix(
  req: Request,
  method: string,
  prefix: string,
): boolean {
  const pathname = new URL(req.url).pathname;
  if (req.method === method.toUpperCase() && pathname.startsWith(prefix)) {
    return true;
  }
  return false;
}

/**
 * hasMethodAndPathRegExp gets a request method and url for an extact match
 * of method and a regular expression match of pathname.
 * @param {Request} req holds the request to evaluatte
 * @param {string} method  holds the target method to match against
 * @param {string} prefix holds the path to compare for an exact match.
 * @returns boolean
 */
export function hasMethodAndPathRegExp(
  req: Request,
  method: string,
  re: RegExp,
): boolean {
  const pathname = new URL(req.url).pathname;
  if (req.method === method.toUpperCase() && re.test(pathname)) {
    return true;
  }
  return false;
}

/**
 * pathIdentifier extracts the identifier from the last element of the URL pathname.
 * The application is expecting a multipart path and if the first "/" and "/" slash
 * it is presumed the identifier is not available.
 *
 * @param {string} u holds the unparsed URL you want to pull the identifier from.
 * @returns {string} idenitifier as a string, empty string means it could not find the identifier.
 *
 * @example
 * ```
 *    const uri = 'https://localhost:8180/groups/LIGO';
 *    const clgid = pathIdentifier(uri);
 *    console.log("group identifier is", clgid);
 * ```
 */
export function pathIdentifier(u: string): string {
  const pathname: string = new URL(u).pathname;
  const cut_pos = pathname.lastIndexOf("/");
  if (cut_pos != pathname.indexOf("/")) {
    return decodeURI(pathname.slice(cut_pos + 1));
  }
  return "";
}

/**
 * getString takes an object and returns the requested attribute if it exists or the default value provided.
 * @param {object} obj is the object to extract the attribute from
 * @param {string} attr is the attribute name
 * @param {string} default_value holds the default to use if attribute is missing
 * @return {string}
 */
export function getString(
  obj: { [k: string]: any },
  attr: string,
  default_value: string,
): string {
  if (obj.hasOwnProperty(attr) && typeof obj[attr] == "string") {
    return obj[attr];
  }
  return default_value;
}

/**
 * getNumber takes an object and returns the requested attribute if it exists or the default value provided.
 * @param {object} obj is the object to extract the attribute from
 * @param {string} attr is the attribute name
 * @param {number} default_value holds the default to use if attribute is missing
 * @return {number}
 */
export function getNumber(
  obj: { [k: string]: any },
  attr: string,
  default_value: number,
): number {
  if (obj.hasOwnProperty(attr) && typeof obj[attr] == "number") {
    return obj[attr];
  }
  return default_value;
}

/**
 * getBoolean takes an object and returns the requested attribute if it exists or the default value provided.
 * @param {object} obj is the object to extract the attribute from
 * @param {string} attr is the attribute name
 * @param {boolean} default_value holds the default to use if attribute is missing
 * @return {boolean}
 */
export function getBoolean(
  obj: { [k: string]: any },
  attr: string,
  default_value: boolean,
): boolean {
  if (obj.hasOwnProperty(attr) && typeof obj[attr] == "boolean") {
    return obj[attr];
  }
  return default_value;
}

/**
 * formDataToObject turn the form data into a simple object.
 *
 * @param {FormData} form data the form object to process
 * @returns {Object}
 */
export function formDataToObject(form: FormData): object {
  const obj: {
    [k: string]: string | boolean | number;
  } = {};
  for (const v of form.entries()) {
    const key: string = v[0];
    if (key !== "submit") {
      const val: any = v[1];
      if (val === "true") {
        obj[key] = true;
      } else if (val === "false") {
        obj[key] = false;
      } else {
        obj[key] = val;
      }
    }
  }
  /*  NOTE: Make sure we update obj.updated */
  obj["updated"] = new Date().toLocaleDateString("en-US") as string;
  return obj;
}
