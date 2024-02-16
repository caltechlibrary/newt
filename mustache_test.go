package newt

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestMustacheEngine(t *testing.T) {
	verbose, out := true, os.Stdout
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MustacheHandleRequest(w, r, verbose, out)
	}))
	defer ts.Close()
	data := bytes.NewBuffer([]byte(`{
	"content_type": "text/html",
	"data": {
		"page_title": "Mustache Test Page",
		"title": "Me and my Mustache",
		"byline": "Me, I'm famous",
		"body": "Talking about my famous mustache, yeah!"
	},
	"template": "<DOCTYPE html><html><head><title>{{page_title}}</title></head><article><h1>{{title}}</h1><h2>{{byline}}</h2><p>{{body}}<p></article></html>"
}`))
	expected := []byte(`<DOCTYPE html><html><head><title>Mustache Test Page</title></head><article><h1>Me and my Mustache</h1><h2>Me, I&#39;m famous</h2><p>Talking about my famous mustache, yeah!<p></article></html>`)

	resp, err := http.Post(ts.URL, "application/json", data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
    if resp.StatusCode != http.StatusOK {
		t.Errorf("response status not OK, %d %q", resp.StatusCode, resp.Status)
		t.FailNow()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if bytes.Compare(expected, body) != 0 {
		t.Errorf("expected %q, got %q", expected, body)
	}
}
