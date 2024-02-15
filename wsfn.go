package newt

import (
	"log"
	"net/http"
	"text/scanner"
	"time"
	"strings"
)

// Logger implementes Newt Project's web logging
type Logger struct {
	handler http.Handler

	// After (defaults to true) logs the request after running the wrapped handler
	After bool
}

//ServeHTTP handles the request by passing it to the real
//handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
	if ! l.After {
    	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
	l.handler.ServeHTTP(w, r)
	if l.After {
    	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

//NewLogger constructs a new Logger middleware handler
func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{
		handler: handlerToWrap,
		After: true,
	}
}

// PatternKeys parses a pattern and returns a list of keys found.
// NOTE: this could be improved to make sure that delimiters are paired and that
// the pattern's names do not contain spaces.
func PatternKeys(p string) []string {
	var s scanner.Scanner
	s.Init(strings.NewReader(p))
	vars := []string{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == '{' {
			s.Scan()
			vName := s.TokenText()
			vars = append(vars, vName)
		}
	}
	return vars
}

