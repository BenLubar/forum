package main

import (
	"html/template"
)

var tmpl = template.New("")

var header = NewTemplate("header", `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>{{.Title}}</title>
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<link href="/css/bootstrap.css?v2.3.2" rel="stylesheet">
	<link href="/css/fontawesome.css?v3.1.0" rel="stylesheet">
</head>
<body>
`)
var footer = NewTemplate("footer", `
</body>
</html>`)

func NewTemplate(name, content string) *template.Template {
	return template.Must(tmpl.New(name).Parse(content))
}
