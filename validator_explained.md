
# Validator Explained

## Overview

It is important for data quality and security to vet submitted data early. Ideally validation happens at least at two points. In browser and in the first stage of a pipeline processing a form submission. Newt accomplishes this by generating a TypeScript validation service and JavaScript to include with the generated form. Deno is used as a runtime server side for the validation code.

TypeScript was choosen because it is a typed language but also can be transpiled to JavaScript. For JS enabled browers this means the same code for validation can be used both server and browser side. That simplifies the code which may need to be further customized such as when someone needs a data type not supported by Newt. Deno was chosen as the runtime because it is robust and much more secure than NodeJS. In fact the validation service can be easily restricted to only respond to network requests in the pipeline.

## How it works

FIXME: To be written
