package newt

import (
	"io/fs"
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

//
// The following is based on the example on the Golang website.
// See https://pkg.go.dev/net/http#example-FileServer-DotFileHiding
//

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fsys dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}


