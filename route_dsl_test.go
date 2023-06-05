package newt

import (
	//"fmt" // DEBUG
	//"os" // DEBUG
	"testing"
)

func TestBlogPathEndingWithString(t *testing.T) {
	rdsl, err := NewRouteDSL(`/blog/${yr Year}/${mo Month}/${dy Day}/${title-slug String}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testMap := map[string]bool{
		"/":                        false,
		"/not-a-post.html":         false,
		"/also-not-a-post":         false,
		"/blog/2023/05/13/my-post": true,
	}
	for p, expected := range testMap {
		//fmt.Fprintf(os.Stderr, "DEBUG p -> %q\n", p)
		val, ok := rdsl.Eval(p)
		//fmt.Fprintf(os.Stderr, "DEBUG val -> %q, ok -> %t\n", val, ok)
		if ok != expected {
			t.Errorf("expected (%q) %t, got (%T) %+v %t", p, expected, val, val, ok)
		}
	}
}

func TestBlogPathEndingWithExt(t *testing.T) {
	rdsl, err := NewRouteDSL(`/blog/${year Year}/${month Month}/${day Day}/${title-slug Basename}${ext Extname}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testMap := map[string]bool{
		"/":                             false,
		"/not-a-post.html":              false,
		"/blog/2023/05/13/my-post.html": true,
	}
	for p, expected := range testMap {
		//fmt.Fprintf(os.Stderr, "DEBUG p -> %q\n", p)
		val, ok := rdsl.Eval(p)
		if ok != expected {
			t.Errorf("expected (%q) %t, got (%T) %+v %t", p, expected, val, val, ok)
		}
	}
}

func TestEvalAndResolve(t *testing.T) {
	rdsl, err := NewRouteDSL(`/blog/${yr Year}/${mo Month}/${dy Day}/${title-slug String}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	pData := map[string]string{
		"/blog/2023/05/29/newt-presentation": "/blog?post_date=2023-05-29&title_slug=newt-presentation",
	}
	for p, expected := range pData {
		m, ok := rdsl.Eval(p)
		if !ok {
			t.Errorf("expected to eval p %q, failed", p)
			t.FailNow()
		} else {
			src := rdsl.Resolve(m, "/blog?post_date=${yr}-${mo}-${dy}&title_slug=${title-slug}")
			if src != expected {
				t.Errorf("expected %q, got %q", expected, src)
			}
		}
	}
}
