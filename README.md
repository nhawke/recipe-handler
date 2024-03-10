# Recipe Server

This is a simple small web server that's meant to serve static recipe pages.  I
made this for myself because I wanted something custom to my needs without
extra bloat and fanciness.

Ideally, it is meant to serve pages rendered from Markdown, but for now it's
only a glorified file server with a few tweaks (such as removing the .md
extension from directory lists).

```
go run main.go /path/to/your/recipe/folder
```
