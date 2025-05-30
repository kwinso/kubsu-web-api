package templates

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed all:templates
var templates embed.FS

func Render(name string, data interface{}) ([]byte, error) {
	filename := "templates/" + name + ".tmpl"
	file, err := templates.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(filename).Parse(string(file))
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
