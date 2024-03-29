package newt

import (
	"fmt"
	"net/mail"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	// 3rd Party packages
	"github.com/shurcooL/github_flavored_markdown"
)


var (
	// DataTypes is a map to the types defined in route_dsl_types.go
	DataTypes = map[string]EvalType{
		"String":   new(TypeString).EvalType,
		"Text":     new(TypeText).EvalType,
		"Integer":  new(TypeInteger).EvalType,
		"Real":     new(TypeReal).EvalType,
		"Boolean":  new(TypeBool).EvalType,
		"Date":     new(TypeDate).EvalType,
		"Year":     new(TypeYear).EvalType,
		"Month":    new(TypeMonth).EvalType,
		"Day":      new(TypeDay).EvalType,
		"Basename": new(TypeBasename).EvalType,
		"Extname":  new(TypeExtname).EvalType,
		"ISBN10":   new(TypeISBN10).EvalType,
		"ISBN13":   new(TypeISBN13).EvalType,
		"ISBN":     new(TypeISBN).EvalType,
		"ISSN":     new(TypeISSN).EvalType,
		"DOI":      new(TypeDOI).EvalType,
		"ISNI":     new(TypeISNI).EvalType,
		"ORCID":    new(TypeORCID).EvalType,
		"ArXiv":    new(TypeArXiv).EvalType,
		"Markdown": new(TypeMarkdown).EvalType,
		"ROR":      new(TypeROR).EvalType,
		"URL":      new(TypeURL).EvalType,
		"EMail":    new(TypeEMail).EvalType,
		"Timestamp": new(TypeTimestamp).EvalType,
	}

	// isRORDomain is the regular expression for the ROR domain per 
	// https://ror.readme.io/docs/ror-identifier-pattern
	isRORDomain = regexp.MustCompile(`^0[a-hj-km-np-tv-z|0-9]{6}[0-9]{2}$`)
)

func lastChar(s string) string {
	l := len(s) - 1
	return s[l:]
}

// DataType is an interface the "data types" need to implement.
type DataType interface {
	// EvalType takes an variable type expression like
	EvalType(string, string) (string, bool)
}

// TypeString implements a string data type, it maps to the SQL VARCHAR
// Max length is 256.
type TypeString struct {
}

func (t TypeString) EvalType(expr string, val string) (string, bool) {
	if len(val) > 256 {
		return val[0:256], false
	}
	return val, true
}

// TypeText represents a text block, this maps to the SQL TEXT
type TypeText struct {
}

func (t TypeText) EvalType(expr string, val string) (string, bool) {
	return val, true
}


// TypeInteger implements an integer data type
type TypeInteger struct {
}

func (t TypeInteger) EvalType(expr string, val string) (string, bool) {
	var layout string
	layout = "%d"
	if strings.Contains(expr, " ") {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	x, err := strconv.Atoi(val)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf(layout, x), true
}

// TypeReal implements a decimal data type
type TypeReal struct {
}

func (t TypeReal) EvalType(expr string, val string) (string, bool) {
	var layout string
	layout = "%f"
	if strings.Contains(expr, " ") {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	x, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf(layout, x), true
}

// TypeBool implements a boolean data type
type TypeBool struct {
}

func (t TypeBool) EvalType(expr string, val string) (string, bool) {
	var layout string
	layout = "%t"
	if strings.Contains(expr, " ") {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf(layout, b), true
}

    
    

// Type Date implements a date type.
type TypeDate struct {
}

func (t TypeDate) EvalType(expr string, val string) (string, bool) {
	var layout string
	if ! strings.Contains(expr, " ") {
		layout = "2006-01-02"
	} else {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	dt, err := time.Parse(layout, val)
	if err != nil {
		return "", false
	}
	return dt.Format(layout), true
}

// TypeYear implements a variation of a Date type for working with the
// year component of a Go Date.
type TypeYear struct {
}

func (t TypeYear) EvalType(expr string, val string) (string, bool) {
	var layout string
	if ! strings.Contains(expr, " ") {
		layout = "2006"
	} else {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	dt, err := time.Parse(layout, val)
	if err != nil {
		return "", false
	}
	return dt.Format(layout), true
}

// TypeMonth implements a variation of a Date type for working with the
// month component of a Go Date.
type TypeMonth struct {
}

func (t TypeMonth) EvalType(expr string, val string) (string, bool) {
	var layout string
	if ! strings.Contains(expr, " ") {
		layout = "01"
	} else {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	dt, err := time.Parse(layout, val)
	if err != nil {
		return "", false
	}
	return dt.Format(layout), true
}

// TypeDay implements a variation of a Date type for working with the
// day component of a Go Date.
type TypeDay struct {
}

func (t TypeDay) EvalType(expr string, val string) (string, bool) {
	var layout string
	if ! strings.Contains(expr, " ") {
		layout = "02"
	} else {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			layout = parts[1]
		}
	}
	dt, err := time.Parse(layout, val)
	if err != nil {
		return "", false
	}
	return dt.Format(layout), true
}

// TypeBasename is a type specific to a path's filename.
type TypeBasename struct {
}

func (t TypeBasename) EvalType(expr string, val string) (string, bool) {
	ext := path.Ext(val)
	return strings.TrimSuffix(val, ext), true
}

// TypeExtname is a type specific to a path's filename's extension.
type TypeExtname struct {
}

func (t TypeExtname) EvalType(expr string, val string) (string, bool) {
	return path.Ext(val), true
}

// TypeISBN10 implements an 10 digit ISBN type.
type TypeISBN10 struct {
}

func isISBN10(val string) bool {
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

func (t TypeISBN10) EvalType(expr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isISBN10(val) {
		return "", false
	}
	return val, true
}

// TypeISBN13 implements an 13 digit ISBN type.
type TypeISBN13 struct {
}

func isISBN13(val string) bool {
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

func (t TypeISBN13) EvalType(extr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isISBN13(val) {
		return "", false
	}
	return val, true
}

// TypeISBN implements both a 10 digit and 13 digit ISBN
type TypeISBN struct {
}

func (t TypeISBN) EvalType(extr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !(isISBN10(val) || isISBN13(val)) {
		return "", false
	}
	return val, true
}

// TypeISSN implements the ISSN data type.
type TypeISSN struct {
}

func isISSN(val string) bool {
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

func (t TypeISSN) EvalType(expr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isISSN(val) {
		return "", false
	}
	return val, true
}

// TypeDOI implements a DOI data type
type TypeDOI struct {
}

func (t TypeDOI) EvalType(expr string, val string) (string, bool) {
	doiRE := regexp.MustCompile(`doi:\s*|(?:https?://)?(?:dx\.)?doi\.org/)?(10\.\d+(.\d+)*/.+)$)`)
	if doiRE.MatchString(val) {
		return val, true
	}
	return "", false
}

// TypeISNI implements an ISNI data type
type TypeISNI struct {
}

func isISNI(val string) bool {
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

func (t TypeISNI) EvalType(expr string, val string) (string, bool) {
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if !isISNI(val) {
		return "", false
	}
	return val, true
}

// TypeORCID impements the ORCID data type
type TypeORCID struct {
}

func (t TypeORCID) EvalType(expr string, val string) (string, bool) {
	if strings.HasPrefix(val, "https://orcid.org/") {
		val = strings.TrimPrefix(val, "https://orcid.org/")
	}
	val = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", ""))
	if isISNI(val) {
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

// TypeMarkdown implements GitHub Flavored Markdown data type using
// the Go package github_flavored_markdown. This is useful many when
// processing form data (e.g. like when processing a POST). The parameter
// passed after "Markdown" should be the name of a varaible that holds
// the processed markup.
type TypeMarkdown struct {
}

func (t TypeMarkdown) GetTargetVarname(expr string) (string, bool) {
	if strings.Contains(expr, " ") {
		parts := strings.SplitN(expr, " ", 2)
		if len(parts) == 2 {
			return parts[1], true
		}
	}
	return "", false
}

func (t TypeMarkdown) EvalType(expr string, val string) (string, bool) {
	return fmt.Sprintf("%s", github_flavored_markdown.Markdown([]byte(val))), true
}

// TypeArXiv implements an arXiv identifier. 
// See https://info.arxiv.org/help/arxiv_identifier.html
type TypeArXiv struct {
}

func (t TypeArXiv) EvalType (expr string, val string) (string, bool) {
	s := strings.ToLower(val)
	if strings.HasPrefix(val, "arxiv:") {
		s = strings.TrimPrefix(strings.ToLower(s), "arxiv:")
	}
	// Find the period after four character YYMM
	if strings.Index(s, ".") != 4 {
		return "", false
	}
	yy, mm := s[0:2], s[2:4]
	_, err := strconv.Atoi(yy)
	if err != nil {
		return "", false
	}
	m, err := strconv.Atoi(mm)
	if (err != nil) || (m < 1) || (m > 12) {
		return "", false
	}
    serial := s[6:]
	vPos := strings.Index(serial, "v")
    if vPos >= 4 {
	   serial = serial[0:vPos]
	}
    if len(serial) < 4 || len(serial) > 5 {
		return "", false
	}
	return strings.TrimPrefix(strings.ToLower(val), "arxiv:"), true
}

// TypeROR implements a naive ROR validator
type TypeROR struct {
}

func (t TypeROR) EvalType (expr string, val string) (string, bool) {
	// FIXME: This is a naive implementation based on https://ror.org/about/faqs/
	// and https://ror.readme.io/docs/ror-identifier-pattern
	rorDomain := val
	if strings.HasPrefix(expr, "https://ror.org/") {
		rorDomain = strings.TrimPrefix(expr, "https://ror.org/")
	} 
	if isRORDomain.MatchString(rorDomain) {
		return "https://ror.org/" + rorDomain, true
	}
	return "", false
}


// TypeURL implements a niave URL validator using net/url's Parser
type TypeURL struct {
}

func (t TypeURL) EvalType (expr string, val string) (string, bool) {
	_, err := url.Parse(val)
	if err != nil {
		return "", false
	}
	return val, true
}

// TypeEMail implements a simlpe EMail validator using Go's mail package
type TypeEMail struct {
}

func (t TypeEMail) EvalType(expr string, val string) (string, bool) {
    _, err := mail.ParseAddress(val)
	if err != nil {
		return "", false
	}
    return val, true
}

// TypeTimestamp implements a simple timestamp format based on Go's time package
// using the RFC3339 format for parsing and writing out.
type TypeTimestamp struct {
}

func (t TypeTimestamp) EvalType(expr string, val string) (string, bool) {
	ts, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return "", false
	}
	return ts.Format(time.RFC3339), true
}
