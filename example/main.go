package main

import (
	"fmt"
	"net/http"
	"os"

	recipes "github.com/nhawke/recipe-handler"

	_ "embed"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Must provide recipe folder as first argument\n")
		os.Exit(1)
	}
	recipeFolderPath := os.Args[1]

	fmt.Fprintf(os.Stderr, "Serving recipes in folder: %v\n", recipeFolderPath)
	http.Handle("/", recipes.NewHandler(recipes.Config{
		Path: recipeFolderPath,
	}))

	err := http.ListenAndServe(":8080", nil)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
