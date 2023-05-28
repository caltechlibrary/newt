package newt

import (
	"path"
	"testing"
)

func TestResolveApiURL(t *testing.T) {
	router := new(Router)
	// Use some test DB_NAME, DB_USERNAME and PASSWORD to test Data API URL transform
	router.Setenv("DB_USERNAME", "foo")
	router.Setenv("DB_PASSWORD", "bar")
	router.Setenv("DB_NAME", "blog")
	if err := router.ReadCSV(path.Join("testdata", "blog_routes.csv")); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// NOTE: pData holds a test data and expected resulting for 
	// our router.
	pData := map[string]string {
		"/blog/2023/05/29/my-post": `http://foo:bar@localhost:3000/blog/?post_date=2023-05-29`,
	}
	for src, expected := range pData {
		// NOTE: we need to ResolveRoute to get a parsed Route DSL merged
		// with any environment vars.
		no, m, ok := router.ResolveRoute(src, "GET", "text/html")
		if ! ok {
			t.Errorf("expedcted to resolve the router %q, failed", src)
			t.FailNow()
		}
		res, ok := router.ResolveApiURL(no, m);
		if ! ok {
			t.Errorf("expected ResolveApiURL to find router %d, for %q", no, src)
			t.FailNow()
		}
		if res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	}
}
