package templates

import (
	"bytes"
	"testing"
)

type GoTemplateTestStruct struct {
	GoTemplate
}

func (gt *GoTemplateTestStruct) Src() string {
	return "hello {{.Name}}!"
}

// Used as test data
type Person struct {
	Name string
}

func TestTemplateSrc(t *testing.T) {
	gt := &GoTemplateTestStruct{}
	gt.Name = "GoTemplateTestStruct"

	if gt.Src() != "hello {{.Name}}!" {
		t.Errorf("Template source string is incorrect. gt.Src(): ", gt.Src())
	}
}

func TestTemplateExecute(t *testing.T) {
	gt := &GoTemplateTestStruct{}
	gt.Name = "GoTemplateTestStruct"

	templ, err := gt.Parse()

	if err != nil {
		t.Errorf("Unable to parse template. Error: %v", err)
	}

	data := Person{Name: "Mary"}

	var buffer bytes.Buffer

	err = templ.Execute(&buffer, data)
	if err != nil {
		t.Errorf("Unable to execute template. Error: %v", err)
	}

	html := buffer.String()

	if html == "" {
		t.Errorf("Generated HTML should not be empty. HTML: %v", html)
	}
}
