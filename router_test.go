package newt

import (
	"path"
	"testing"
)

func TestRouter(t *testing.T) {
	_, err := RouterFromCSV(path.Join("testdata", "blog_routes.csv"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
