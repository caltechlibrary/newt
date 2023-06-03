package newt

import (
	"os"
	"path"
	"testing"
)

func TestResolveApiURL(t *testing.T) {
	router := new(Router)
	// NOTE: Use some test DB_NAME, DB_USERNAME and PASSWORD to test
	// Data API URL transform. This isn't actually going to contect
	// to anything in the test.
	router.Setenv("DB_USERNAME", "foo")
	router.Setenv("DB_PASSWORD", "bar")
	router.Setenv("DB_NAME", "blog")
	if err := router.ReadCSV(path.Join("testdata", "blog_routes.csv")); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// NOTE: pData holds a test data and expected resulting for
	// our router.
	pData := map[string]string{
		"/blog/2023/05/29/my-post":      `http://foo:bar@localhost:3000/blog/?post_date=2023-05-29&title_slug=my-post`,
		"/blog/2023/06/01/my-post.html": `http://foo:bar@localhost:3000/blog/?post_date=2023-06-01&title_slug=my-post.html`,
	}
	for src, expected := range pData {
		// NOTE: we need to ResolveRoute to get a parsed Route DSL merged
		// with any environment vars.
		no, m, ok := router.ResolveRoute(src, "GET")
		if !ok {
			t.Errorf("expected to resolve the router %q, failed", src)
			t.FailNow()
		}
		res, ok := router.ResolveApiURL(no, m)
		if !ok {
			t.Errorf("expected ResolveApiURL to find router %d, for %q", no, src)
			t.FailNow()
		}
		if res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	}
	// Now we test urls that should fail to resolve ...
	pData = map[string]string{
		"/about.html": ``,
	}
	for src, unexpected := range pData {
		// NOTE: we need to ResolveRoute to get a parsed Route DSL merged
		// with any environment vars.
		no, _, ok := router.ResolveRoute(src, "GET")
		if ok {
			t.Errorf("expected %q to fail to resolve the router, (%d) %q", src, no, unexpected)
			t.FailNow()
		}
	}
}

// TestRequestDataAPI
func TestRequestDataAPI(t *testing.T) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" || dbUsername == "" || dbPassword == "" {
		t.Skipf("TestRequestDataAPI requires DB_NAME, DB_USERNAME, DB_PASSWORD in the environment")
		return
	}

	router := new(Router)
	// Use some test DB_NAME, DB_USERNAME and PASSWORD to test Data API URL transform
	fName := path.Join("testdata", dbName+"-routes.csv")
	router.Setenv("DB_NAME", dbName)
	router.Setenv("DB_USERNAME", dbUsername)
	router.Setenv("DB_PASSWORD", dbPassword)
	if err := router.ReadCSV(fName); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(router.Routes) != 2 {
		t.Errorf("expected 2 routes defined in %q, got %d", fName, len(router.Routes))
		t.FailNow()
	}
	// Map of request to row number
	tMap := map[string]int{
		"http://localhost:3000/sighting":  0,
		"http://localhost:3000/bird_view": 1,
	}
	for apiURL, rNo := range tMap {
		src, _, statusCode := router.RequestDataAPI(rNo, apiURL, []byte{})
		if statusCode != 200 {
			t.Errorf("expected 200 status code, got %d, %s", statusCode, src)
		}
	}
}

// TestPandoc
func TestPandoc(t *testing.T) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" || dbUsername == "" || dbPassword == "" {
		t.Skipf("TestPandoc requires DB_NAME, DB_USERNAME, DB_PASSWORD in the environment")
		return
	}

	router := new(Router)
	// Use some test DB_NAME, DB_USERNAME and PASSWORD to test Data API URL transform
	fName := path.Join("testdata", dbName+"-routes.csv")
	router.Setenv("DB_NAME", dbName)
	router.Setenv("DB_USERNAME", dbUsername)
	router.Setenv("DB_PASSWORD", dbPassword)
	if err := router.ReadCSV(fName); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(router.Routes) != 2 {
		t.Errorf("expected 2 routes defined in %q, got %d", fName, len(router.Routes))
		t.FailNow()
	}
	// Map of request to row number
	tMap := map[string]int{
		"http://localhost:3000/sighting":  0,
		"http://localhost:3000/bird_view": 1,
	}
	dMap := map[string][]byte{}
	for apiURL, rNo := range tMap {
		src, _, statusCode := router.RequestDataAPI(rNo, apiURL, []byte{})
		if statusCode != 200 {
			t.Errorf("expected 200 status code (data api), got %d", statusCode)
		}
		dMap[apiURL] = src
	}
	for k, v := range dMap {
		rNo, ok := tMap[k]
		if ok {
			if len(router.Routes[rNo].PandocTemplate) == 0 {
				t.Errorf("expected pandoc template for route %d, %q", rNo, apiURL)
				continue
			}
			src, _, statusCode := router.RequestPandoc(rNo, v)
			if statusCode != 200 {
				t.Errorf("expected 200 status code (pandoc api), got %d, %s", statusCode, src)
			}
		} else {
			t.Errorf("could not retrieve route number, %q", k)
		}
	}
}
