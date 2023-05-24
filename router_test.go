package newt

import (
	"path"
	"testing"
)

func TestRouter(t *testing.T) {
	router := new(Router)
	if err := router.ReadCSV(path.Join("testdata", "blog_routes.csv")); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
