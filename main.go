package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	recipeFolderPath = "/mnt/n/recipes/"
	outFolder        = "./out"
)

func main() {
	http.HandleFunc("/", serve)
	http.ListenAndServe(":8080", nil)

}

func serve(w http.ResponseWriter, req *http.Request) {
	page := req.URL.Path

	if page == "/" {
		dirList(w, req)
		return
	}

	eprintln(page)
	if lpage := strings.ToLower(page); !(strings.HasSuffix(lpage, "/") || strings.HasSuffix(lpage, ".md")) {
		page += ".md"
	}

	fpath := path.Join(recipeFolderPath, page)
	eprintln(fpath)
	http.ServeFile(w, req, path.Clean(fpath))
}

// Listing a dir, modified to exclude the .md suffix and hidden files.
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

// convenience functions

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
