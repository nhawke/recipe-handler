# Recipe Handler

This is a simple HTTP handler library that's meant to serve static recipe
pages.  I made this for myself because I wanted something custom to my needs
without extra bloat and fanciness. It serves recipes formatted in markdown from
a provided directory.

This repo also contains barebones example server whose only handler is the recipe
handler. The exmaple server listens on port `8080`.

```
go run example/main.go /path/to/your/recipe/folder
```
