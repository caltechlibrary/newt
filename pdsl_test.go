package newt

import (
	"testing"
)

func TestBlogPathEndingWithString(t *testing.T) {
	pdsl, err := NewPathDSL(`/blog/{yr Year}/{mo Month}/{dy Day}/{title-slug String}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var (
		year pdslYear
		month pdslMonth
		day pdslDay
		str pdslString
	)
	err = pdsl.RegisterType("Year", year)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Month", month)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Day", day)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("String", str)
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
	pdsl, err := NewPathDSL(`/blog/{yr Year}/{mo Month}/{dy Day}/{title-slug Basename}{ext Extname}`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var (
		year pdslYear
		month pdslMonth
		day pdslDay
		str pdslString
		basename pdslBasename
		extname pdslExtname
	)

	err = pdsl.RegisterType("Year", year)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Month", month)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Day", day)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("String", str)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Basename", basename)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = pdsl.RegisterType("Extname", extname)
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
