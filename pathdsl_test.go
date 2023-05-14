package newt

import (
	"path"
	"strings"
	"testing"
	"time"
)

func TestBlogPathEndingWithString(t *testing.T) {
	pdsl, err := NewPathDSL(`/blog/{year Year}/{month Month}/{day Day}/{title-slug String}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Year", func(expr string, val string) (interface{}, bool) {
		dt, err := time.Parse(`2006`, val)
		if err != nil {
			return "", false
		}
		return dt.Format(`2006`), true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Month", func(expr string, val string) (interface{}, bool) {
		dt, err := time.Parse(`01`, val)
		if err != nil {
			return "", false
		}
		return dt.Format(`01`), true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Day", func(expr string, val string) (interface{}, bool) {
		dt, err := time.Parse(`02`, val)
		if err != nil {
			return "", false
		}
		return dt.Format(`02`), true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("String", func(expr string, val string) (interface{}, bool) {
		return val, true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	testMap := map[string]bool{
		"/":                             false,
		"/my-post.html":  false,
		"/blog/2023/05/13/my-post.html": true,
	}
	for p, expected := range testMap {
		val, ok := pdsl.Eval(p)
		if ok != expected {
			t.Errorf("expected (%q) %t, got (%T) %+v %t", p, expected, val, val, ok)
		}
	}
}


func TestBlogPathEndingWithExt(t *testing.T) {
	pdsl, err := NewPathDSL(`/blog/{year Year}/{month Month}/{day Day}/{title-slug Basename}{ext Extname}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Year", func(expr string, val string) (interface{}, bool) {
		dt, err := time.Parse(`2006`, val)
		if err != nil {
			return "", false
		}
		return dt.Format(`2006`), true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Month", func(expr string, val string) (interface{}, bool) {
		dt, err := time.Parse(`01`, val)
		if err != nil {
			return "", false
		}
		return dt.Format(`01`), true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Day", func(expr string, val string) (interface{}, bool) {
		dt, err := time.Parse(`02`, val)
		if err != nil {
			return "", false
		}
		return dt.Format(`02`), true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("String", func(expr string, val string) (interface{}, bool) {
		return val, true
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Basename", func(expr string, val string) (interface{}, bool) {
		ext := path.Ext(val)
		return strings.TrimSuffix(val, ext), true
	})
	err = pdsl.RegisterType("Extname", func(expr string, val string) (interface{}, bool) {
		return path.Ext(val), true
	})

	testMap := map[string]bool{
		"/":                             false,
		"/my-post.html":  false,
		"/blog/2023/05/13/my-post.html": true,
	}
	for p, expected := range testMap {
		val, ok := pdsl.Eval(p)
		if ok != expected {
			t.Errorf("expected (%q) %t, got (%T) %+v %t", p, expected, val, val, ok)
		}
	}
}
