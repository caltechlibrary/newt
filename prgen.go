package newt

import (
	"io"

	// 3rd Party Templates
	"github.com/cbroglie/mustache"
)

// postgRESTConf renders a PostgREST configuration given a namespace value.
func postgRESTConf(out io.Writer, namespace string, port string) error {
	txt := `
db-uri = "postgres://{{namespace}}_authenticator:{{secret}}@localhost:{{port}}/postgres"
db-schemas = "{{nampespace}}"
db-anon-role = "{{namespace}}_anonymous"
`
	tmpl, err := mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data := map[string]string{
		"namespace": namespace,
		"port": port,
		"secret": "__change_me_password_goes_here__",
	}
	return tmpl.FRender(out, data)
}
