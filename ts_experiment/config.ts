/**
 * config.ts holds the application configuration defaults.
 */
import { yaml, ASTInterface, TemplateInterface } from "./deps.ts";

export const defaultHttpPort: number = 3032;

export async function loadAndMergeNewtYAML(
  yaml_name: string,
  options: { [k: string]: any },
): Promise<ASTInterface> {
  const src = await Deno.readTextFile(yaml_name);
  const result: ASTInterface = (await yaml.parse(
    src,
  )) as unknown as ASTInterface;
  console.log("DEBUG AST interface", result);
  if (options.port !== undefined && typeof options.port === "number" && options.port > 0) {
   	result.applications.template_engine.port = options.port;
  }
  console.log("DEBUG AST interface (after merge)", result);
  return result;
}

export function getTemplateFromId(config: ASTInterface, template_id: string) : TemplateInterface | undefined {
	if (config.templates == undefined) {
		return undefined;
	}
	console.log("DEBUG config.templates type ->", typeof config.templates);
	for (let i of config.templates) {
		console.log("DEBUG i of config.templates ->", typeof i, i)
		/*
		const template = config.templates[i];
		if (template.id === template_id) {
			return template;
		}
		*/
	}
	return undefined;
}
