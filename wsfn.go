package newt

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"text/scanner"
	"time"
)

// Logger implementes Newt Project's web logging
type Logger struct {
	// handler holds the handler (mux) you're wrapping with your logger.
	handler http.Handler

	// After (defaults to true) logs the request after running the wrapped handler
	After bool

	// Verbose (defaults to false) show the contents of a GET, POST, PUT and DELETE in log output
	Verbose bool
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	buf := []byte{}
	header := []byte{}
	if l.Verbose {
		header, _ = json.MarshalIndent(r.Header, "", "  ")
		switch r.Method {
			case http.MethodPost :
				buf, _ = io.ReadAll(r.Body)
			case http.MethodGet:
				values := r.URL.Query()
				buf = []byte(fmt.Sprintf("%+v", values))
			case http.MethodPut:
				buf, _ = io.ReadAll(r.Body)
			case http.MethodPatch:
				buf, _ = io.ReadAll(r.Body)
			default:
				buf = []byte("body not captured for method "+r.Method)
		}
	}
	if !l.After {
		if l.Verbose {
			log.Printf("%s %s %v\nHeaders ->\n%s\n\n Body ->\n%s\n\n", r.Method, r.URL.Path, time.Since(start), header, buf)
		} else {
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		}
	}
	if l.handler != nil && l.handler.ServeHTTP != nil {
		l.handler.ServeHTTP(w, r)
	} else {
		log.Printf("%s %s %v handler or ServiceHTTP is nil", r.Method, r.URL.Path, time.Since(start))
	}
	if l.After {
		if l.Verbose {
			log.Printf("%s %s %v\nHeaders ->\n%s\n\n Body ->\n%s\n\n", r.Method, r.URL.Path, time.Since(start), header, buf)
		} else {
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		}
	}
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{
		handler: handlerToWrap,
		After:   true,
		Verbose: false,
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

func NewtStaticFileServer(port int, htdocs string, verbose bool) error {
	if htdocs == "" {
		htdocs = "."
	}
	if port == 0 {
		port = 8000
	}
	mux := http.NewServeMux()
	fsys := dotFileHidingFileSystem{http.Dir(htdocs)}
	mux.Handle("/", http.FileServer(fsys))
	// Wrap our mux in a logger, it will become the wrapper of our mux handed off to the server
	logger := NewLogger(mux)
	logger.Verbose = verbose
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = fmt.Sprintf(":%d", port)
	svr.Handler = logger
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
