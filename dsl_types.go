package newt

import (
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// RdslTypes is a map to the types defined in route_dsl_types.go
	RouteTypes = map[string]EvalType{
		"String":   new(RdslString).EvalType,
		"Year":     new(RdslYear).EvalType,
		"Month":    new(RdslMonth).EvalType,
		"Day":      new(RdslDay).EvalType,
		"Basename": new(RdslBasename).EvalType,
		"Extname":  new(RdslExtname).EvalType,
		"Isbn10":   new(RdslIsbn10).EvalType,
		"Isbn13":   new(RdslIsbn13).EvalType,
		"Isbn":     new(RdslIsbn).EvalType,
		"Issn":     new(RdslIssn).EvalType,
		"DOI":      new(RdslDOI).EvalType,
		"Isni":     new(RdslIsni).EvalType,
		"ORCID":    new(RdslORCID).EvalType,
	}
)

func lastChar(s string) string {
	l := len(s) - 1
	return s[l:]
}

// DSLType
type DSLType interface {
	// EvalType takes an variable type expression like
	EvalType(string, string) (string, bool)
}

// Route DSL types
type RdslString struct {
}

func (str RdslString) EvalType(expr string, val string) (string, bool) {
	return val, true
}

type RdslYear struct {
}

func (year RdslYear) EvalType(expr string, val string) (string, bool) {
	dt, err := time.Parse(`2006`, val)
	if err != nil {
		return "", false
	}
	return dt.Format(`2006`), true
}

type RdslMonth struct {
}

func (month RdslMonth) EvalType(expr string, val string) (string, bool) {
	dt, err := time.Parse(`01`, val)
	if err != nil {
		return "", false
	}
	return dt.Format(`01`), true
}

type RdslDay struct {
}

func (day RdslDay) EvalType(expr string, val string) (string, bool) {
	dt, err := time.Parse(`02`, val)
	if err != nil {
		return "", false
	}
	return dt.Format(`02`), true
}

type RdslBasename struct {
}

func (basename RdslBasename) EvalType(expr string, val string) (string, bool) {
	ext := path.Ext(val)
	return strings.TrimSuffix(val, ext), true
}

type RdslExtname struct {
}

func (extname RdslExtname) EvalType(expr string, val string) (string, bool) {
	return path.Ext(val), true
}

type RdslIsbn10 struct {
}

func isIsbn10(val string) bool {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if len(val) != 10 {
		return false
	}
	r := 0
	for i := 0; i < len(val); i++ {
		x, err := strconv.Atoi(val[i:1])
		if err != nil {
			return false
		}
		r += (10 - i) * x
	}
	if !((r % 11) == 0) {
		return false
	}
	return true
}

func (isbn10 RdslIsbn10) EvalType(expr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isIsbn10(val) {
		return "", false
	}
	return val, true
}

type RdslIsbn13 struct {
}

func isIsbn13(val string) bool {
	if len(val) != 13 {
		return false
	}
	y, err := strconv.Atoi(lastChar(val))
	if err != nil {
		return false
	}
	r := 0
	for i := 0; i < len(val); i++ {
		x, err := strconv.Atoi(val[i : i+1])
		if err != nil {
			return false
		}
		if (i % 2) == 0 {
			// Even just sum the value
			r += x
		} else {
			// Odd so weight the value by 3
			r += (x * 3)
		}
	}
	// Finalize the check sum and compare with the final digit's value
	chk := ((10 - r) % 10)
	if chk != y {
		return false
	}
	return true
}

func (isbn13 RdslIsbn13) EvalType(extr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isIsbn13(val) {
		return "", false
	}
	return val, true
}

type RdslIsbn struct {
}

func (isbn RdslIsbn) EvalType(extr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !(isIsbn10(val) || isIsbn13(val)) {
		return "", false
	}
	return val, true
}

type RdslIssn struct {
}

func isIssn(val string) bool {
	if len(val) != 8 {
		return false
	}
	r := 0
	for i := 0; i < len(val); i++ {
		x, err := strconv.Atoi(val[i : i+1])
		if err != nil {
			return false
		}
		r += (8 - i) * x
	}
	if !((r % 11) == 0) {
		return false
	}
	return true
}

func (issn RdslIssn) EvalType(expr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isIssn(val) {
		return "", false
	}
	return val, true
}

type RdslDOI struct {
}

func (doi RdslDOI) EvalType(expr string, val string) (string, bool) {
	doiRE := regexp.MustCompile(`doi:\s*|(?:https?://)?(?:dx\.)?doi\.org/)?(10\.\d+(.\d+)*/.+)$)`)
	if doiRE.MatchString(val) {
		return val, true
	}
	return "", false
}

type RdslIsni struct {
}

func isIsni(val string) bool {
	if len(val) != 16 {
		return false
	}
	y, err := strconv.Atoi(lastChar(val))
	if err != nil {
		return false
	}
	r := 0
	for x := 0; x < len(val); x++ {
		x, err := strconv.Atoi(val[x : x+1])
		if err != nil {
			return false
		}
		r = (r + x) * 2
	}
	if chk := (12 - r%11) % 11; chk != y {
		return false
	}
	return true
}

func (isni RdslIsni) EvalType(expr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isIsni(val) {
		return "", false
	}
	return val, true
}

type RdslORCID struct {
}

func (orcid RdslORCID) EvalType(expr string, val string) (string, bool) {
	if strings.HasPrefix(val, "https://orcid.org/") {
		val = strings.TrimPrefix(val, "https://orcid.org/")
	}
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if isIsni(val) {
		// Trim the check sum digit
		val = val[0 : len(val)-1]
		chk, err := strconv.Atoi(val)
		if err != nil {
			return "", false
		}
		if (chk >= 15000000) && (chk <= 35000000) {
			return val, true
		}
	}
	return "", false
}
