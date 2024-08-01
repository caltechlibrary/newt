/**
 * newt information
 */
export const appInfo: {[k: string]: string} = {
  // appName holds the application/package name
  appName: "newt",

  // Version number of release
  version: "0.0.9",

  // ReleaseDate, the date version.ts was generated
  releaseDate: "2024-08-01",

  // ReleaseHash, the Git hash when version.go was generated
  releaseHash: "f466728",

  // licenseText holds a copy of the application license text.
  licenseText: `
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

`
}

/**
 * fmtHelp lets you process a text block with simple curly brace markup.
 * @param {src} src holds the help text to be processed with curly brace
 * refernces for app_name, version, release_date, release_hash.
 * @returns {string} returns the rendered help text.
 *
 * @example
 * ```
 *   import { appInfo, fmtHelp } from './version.ts';
 *
 *   const helpText = ` ... this is where you document your program in Pandoc manpage format `;
 *   console.log(fmtHelp(helpText, appInfo));
 * ```
 */
export function fmtHelp(src: string, appInfo: {[k: string]: string}): string {
  const terms: { [k: string]: string } = {
    app_name: appInfo.appName,
    version: appInfo.version,
    release_date: appInfo.releaseDate,
    release_hash: appInfo.releaseHash
  };
  for (let key in terms) {
    const val: string = terms[key];
    const varname: string = ['{', key, '}'].join('');
    if (src.indexOf(varname) > -1) {
      src = src.replaceAll(varname, val)
    }
  }
  return src
}
