/**
 * options.ts is a command processing packaged. It is provided as an alternative
 * to Deno's own parse-args package. Main feature is that it supports an Options object
 * that also maps help strings to the command line options. It supports options defined
 * individuall as boolean, number and string.
 */

/**
 * matchType takes a source variable and a target and
 * attempts to return a value from the target that matches
 * the source type. It only works for the simple types of
 * string, number and bool.
 *
 * @param {any} source
 * @param {any} target
 * @returns {any} target value as source type or undefined
 *
 * @example
 *
 * ```
 *   const a = true;
 *   let b = "true";
 *   let c = matchType(a, b);
 * ```
 */
export function matchType(source: any, target: any): any {
  const sourceType = typeof source;
  const targetType = typeof target;
  if (sourceType === targetType) {
    return target;
  }
  if (sourceType === "string") {
    return target as string;
  }
  if (sourceType == "number") {
    return (new Number(target)).valueOf();
  }
  if (sourceType === "boolean") {
    if (typeof target == "string") {
      const val = target.toLocaleLowerCase();
      if (
        val === "true" || val === "t" || val === "1" || val === "yes" ||
        val === "ok"
      ) {
        return true;
      }
      return false;
    }
    if (typeof target == "number") {
      if (target > 0) {
        return true;
      }
      return false;
    }
    return target as unknown as boolean;
  }
  return undefined;
}

/**
 * firstAndRest splits a string or the delimiter return an array of two strings.
 * If a value isn't available then the string will be an empty string.
 *
 * @param {string} source is the source string to be split into first and rest strings
 * @param {RegExp} re is a JavaScript regular expression object used to make the split.
 * @returns [string, string] composed of to strings first and rest where an empty string
 * means there was no avaialble value.
 *
 * @example
 *
 * ```
 *   const name = "fred thomas zip";
 *   let name_parts: [string, string] = firstAndRest(name, / /);
 *   console.log(name_parts);
 * ```
 */
export function firstAndRest(source: string, re: RegExp): [string, string] {
  const parts = source.split(re);
  if (parts.length > 1) {
    let first = parts.shift();
    if (first === undefined) {
      first = "";
    }
    const rest = source.replace(first, "").replace(re, "");
    return [first, rest];
  }
  if (parts.length == 1) {
    return [parts[0], ""];
  }
  return ["", ""];
}

/**
 * OptionsProcessor is a class for processing command line parameters or arrays that contain parameters.
 * The assumptions is that the array of parameters is structure with "options" coming before
 * ordered parameters. Options start with a delimiter (e.g. "-"). Option values maybe passed
 * as the next parameter in the list or joined to the option by an equal sign.
 *
 * Remaining ordered args returned in the object's args array. Option
 * settings are returns in the options object as key/value pairs.
 *
 * The Options object let's you define valid options using the {type_name}Var methods.
 * After the parse method is called the Options object attributes of options and args will be
 * set. Options.options holds the key/value pairs with simple type enfocement, the Options.args
 * is an array of ordered paremeters.
 *
 * You may display a list of supported options by iterating over Options.defaults
 *
 * @example
 *
 * ```
 *  op = new OptionsProcessor();
 *  op.boolVar("help", false, "display help message");
 *  op.stringVar("url", "http://localhost:8000", "set the URL");
 *  op.numberVar("retry", 3, "set the default number of retries");
 *  op.parse(Deno.args)
 *  if (op.options.help) {
 *     // ... show the help page ...
 *  }
 *  for (argv in op.args) {
 *     console.log("ordered argv: ", argv);
 *  }
 *  console.log("url: " + op.options.url);
 *  console.log("retry count: ", op.options.retry);
 * ```
 */
export class OptionsProcessor {
  defaults: { [k: string]: any } = {};
  help: { [k: string]: string } = {};
  readonly options: { [k: string]: any } = {};
  readonly args: string[] = [];

  /**
   * booleanVar adds a boolean option to be processed. It takes the option text (e.g. "help", "license", "version")
   * and tracks the default value and type associating a help message with the specific option.
   * These are used by `Options.parse({string[]})` in processing a list of strings that contain options and ordered parameters.
   *
   * @param {string} option holds the string used to identify the option. It will become the attribute name in the
   * options attribute of the Options object.
   * @param {boolean} this is the default value to use if the option is NOT present when processed.
   * @param {string} msg holds the help text association with the option.
   *
   * @example
   * ```
   *   op = new OptionsProcessor();
   *   op.booleanVar('help', false, 'display help')
   *   op.parse(Deno.args);
   *
   *   if (op.options.help) {
   *      * ... display help page ...
   *   }
   * ```
   */
  booleanVar(option: string, defaultVal: boolean, msg: string) {
    this.defaults[option] = defaultVal;
    this.help[option] = msg;
  }

  /**
   * stringVar adds a string option to be processed. It takes the option text (e.g. "f", "file", "o", "output")
   * and tracks the default value and type associating a help message with the specific option.
   * These are uses by `Options.parse({string[]})` in processing a list of strings that contain options and
   * ordered parameters.
   *
   * @param {string} option holds the option flag. This also becomes the attribute name.
   * @param {string} defaultVal holds the default string value to use when the option is NOT present on the command line.
   * @param {string} msg holds the help text associated with the option.
   *
   * @example
   *
   * ```
   *   op = new OptionsProcessor();
   *   op.stringVar('i', '', 'set the input filename')
   *   op.stringVar('o', '', 'set the output filename')
   *   op.parse(Deno.args);
   *   console.log(`Input file is ${op.i}`);
   *   console.log(`output file is ${op.o}`);
   * ```
   */
  stringVar(option: string, defaultVal: string, msg: string) {
    this.defaults[option] = defaultVal;
    this.help[option] = msg;
  }

  /**
   * numberVar adds a number option to be processed. It takes the option text (e.g. "retry", "port")
   * and tracks the default value and type associating a help message with the specific option.
   * These are uses by `Options.parse({string[]})` in processing a list of strings that contain options and
   * ordered parameters.
   *
   * @param {string} option holds the option flag. This also becomes the attribute name.
   * @param {string} defaultVal holds the default number value to use when the option is NOT present on the command line.
   * @param {string} msg holds the help text associated with the option.
   *
   * @example
   *
   * ```
   *   op = new OptionsProcessor();
   *   op.numberVar('retry', 3, 'set the number of retries');
   *   op.numberVar('port', 8800, 'set the port number');
   *   op.parse(Deno.args);
   *   console.log(`Retry count ${op.retry}`);
   *   console.log(`Port number ${op.port}`);
   * ```
   */
  numberVar(option: string, defaultVal: number, msg: string) {
    this.defaults[option] = defaultVal;
    this.help[option] = msg;
  }

  /**
   * parse takes an array of string such as those supplied by `Deno.args` and process them seperating
   * options and ordered parameters.
   *
   * @param {string[]} args is the list of strings to process for options and ordered parameters
   *
   * @example
   *
   * ```
   *   let cmd_options: string[] = [ '-help', 'one', 'two', 'three' ];
   *   op = new OptionsProcessor();
   *   op.boolVar("help", false, "display help");
   *   op.parse(cmd_options);
   *   if (op.options.help) {
   *      // ... display help page ...
   *   }
   *   console.log("all options are ", op.options);
   *   console.log("ordered parameters are ", op.args);
   * ```
   */
  parse(args: string[]): string[] {
    let optionalParameters: boolean = true;
    let errors: string[] = [];
    for (let i = 0; i < args.length; i++) {
      const argv = args[i];
      if (optionalParameters && argv.startsWith("-")) {
        // handle case where arg contains equal sign before value
        let parts = firstAndRest(argv.replace(/^-+/, ""), /=/);
        let optname = parts.shift();
        let rest = parts.shift();
        // handle case where arg's value is next arg
        if ((rest === undefined) || (rest === "")) {
          if (
            ((i + 1) < args.length) && (args[i + 1].startsWith("-") === false)
          ) {
            rest = args[i + 1];
            i++;
          }
        }
        if ((optname !== undefined) && this.defaults.hasOwnProperty(optname)) {
          const defval = this.defaults[optname];
          let optval = matchType(defval, rest);
          if (rest === "") {
            if (typeof defval === "boolean") {
              // Handle the case where a boolean flag turns a value true
              optval = !defval;
            }
          }
          if (optval !== undefined) {
            this.options[optname] = optval;
          } else {
            errors.push(
              `-${optname} cannot be converted to type ${typeof defval}, skipping`,
            );
          }
        } else {
          // Handle case where arg value comes as separate arg
          errors.push(
            `option no. ${i}: "${
              args[i]
            }" -> "${argv}" is an unsupported option, ignoring`,
          );
        }
      } else {
        // If we're in the ordered args then we stop processing the optional
        // parameters.
        optionalParameters = false;
        this.args.push(argv);
      }
    }
    // Finally we need to set our default values if the options doesn't contain them.
    for (let attr in this.defaults) {
      if (!this.options.hasOwnProperty(attr)) {
        this.options[attr] = this.defaults[attr];
      }
    }
    return errors;
  }
}
