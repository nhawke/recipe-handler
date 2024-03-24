package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"unicode"

	_ "embed"
)

var (
	recipeFolderPath string

	//go:embed recipe.html.tmpl
	recipeTmplText string
	recipeTmpl     *template.Template
)

func main() {
	if len(os.Args) < 2 {
		exit("Must provide recipe folder as first argument\n")
	}
	recipeFolderPath = os.Args[1]

	var err error
	recipeTmpl, err = template.New("recipe").Parse(recipeTmplText)
	if err != nil {
		exit("Failed to parse reipce Template: %v\n", err)
	}

	eprintf("Serving recipes in folder: %v\n", recipeFolderPath)
	http.HandleFunc("/", serve)

	err = http.ListenAndServe(":8080", nil)
	eprintln(err)
}

func serve(w http.ResponseWriter, req *http.Request) {
	page := req.URL.Path

	if page == "/" {
		dirList(w, req)
		return
	}

	eprintf("Request path: %v\n", page)
	if lpage := strings.ToLower(page); !(strings.HasSuffix(lpage, "/") || strings.HasSuffix(lpage, ".md")) {
		page += ".md"
	}

	fpath := path.Join(recipeFolderPath, page)
	eprintf("Serving page: %v\n", fpath)
	serveFile(w, req, path.Clean(fpath))
}

// dirList lists the contents of a directory, excluding the .md
// suffix and hidden files with a . prefix.
func dirList(w http.ResponseWriter, r *http.Request) {
	dir, err := os.ReadDir(recipeFolderPath)

	if err != nil {
		errString := fmt.Sprintf("Error reading recipe folder %q: %v\n", recipeFolderPath, err)
		http.Error(w, errString, 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")

	for _, dirE := range dir {
		name := strings.TrimSuffix(dirE.Name(), ".md")

		// skip hidden files
		if name[0] == '.' {
			continue
		}

		// Properly serialize URL
		url := url.URL{Path: name}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), name)
	}

	fmt.Fprintf(w, "</pre>\n")
}

func serveFile(w http.ResponseWriter, r *http.Request, path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, fs.ErrPermission) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	page := Page{
		Title: pageTitle(path),
		Body:  string(wrapLines(b)),
	}

	if err := recipeTmpl.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Page struct {
	Title string
	Body  string
}

func pageTitle(p string) string {
	file := path.Base(p)
	name, _, _ := strings.Cut(file, ".")
	return name
}

const lineLength = 80

func wrapLines(b []byte) []byte {
	sb := string(b)
	out := bytes.Buffer{}

	lpos := 0
	for _, c := range sb {
		if c == '\n' {
			lpos = 0
		} else if lpos >= 80 && unicode.IsSpace(c) {
			// Wrap line only if at the end of the word.
			out.WriteRune('\n')
			lpos = 0
		}

		if lpos == 0 && c != '\n' && unicode.IsSpace(c) {
			// Don't start a created line with a space.
			continue
		}

		out.WriteRune(c)
		lpos++
	}

	return out.Bytes()
}

func eprintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func eprintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func exit(format string, args ...any) {
	eprintf(format, args...)
	os.Exit(1)
}
