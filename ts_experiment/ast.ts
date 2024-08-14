export interface ASTInterface {
  applications: {
    [k: string]: ApplicationInterface;
  };
  models: ModelsInterface[];
  routes: RouteInterface[];
  templates: TemplateInterface[];
  options: { [k: string]: string }[];
  vars: string[];
}

export interface ApplicationInterface {
  app_path: string;
  conf_path: string;
  port: number;
  base_dir: string;
  ext_name: string;
  partials_dir: string;
  layouts_dir: string;
  cache_partials: string;
  default_layout: string;
  helpers: string;
  compiler_options: {[k:string]: any};
}

export interface ModelsInterface {
  id: string;
  attributes: { [k: string]: string };
  description: string;
  elements: ElementInterface[];
  title: string;
}

export interface ElementInterface {
  type: string;
  id: string;
  attributes: { [k: string]: string };
  pattern: RegExp;
  options: { [k: string]: string }[];
  is_object_id: boolean;
  label: string;
}

export interface RouteInterface {
  id: string;
  pattern: RegExp;
  description: string;
  pipeline: ServiceInterface[];
  debug: boolean;
}

export interface ServiceInterface {
  service: string;
  descirption: string;
  timeout: number;
}

export interface TemplateInterface {
	id: string;
	pattern: RegExp,
	template: string;
	document: {[k:string]:string},
	vacabularies: string[];
	vars: {[k:string]:string};
	description: string,
	debug: boolean;
}

