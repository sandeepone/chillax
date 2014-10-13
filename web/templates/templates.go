package templates

import (
	html_template "html/template"
)

type GoTemplate struct {
	Name string
}

func (gt *GoTemplate) Parse() (*html_template.Template, error) {
	return html_template.New(gt.Name).Parse(gt.Src())
}

func (gt *GoTemplate) Src() string {
	return ""
}
