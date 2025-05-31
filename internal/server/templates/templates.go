package templates

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed all:templates
var templates embed.FS

func inArray(val int32, slice []int32) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func Render(name string, data interface{}) ([]byte, error) {
	filename := "templates/" + name + ".html"
	file, err := templates.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.
		New(filename).
		Funcs(template.FuncMap{
			"inArray": inArray,
		}).
		Parse(string(file))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, filename, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
