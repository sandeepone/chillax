package templates

import (
	"bytes"
	"os"
	"testing"
)

type GoTemplateTestStruct struct {
	GoTemplate
}

func (gt *GoTemplateTestStruct) String() string {
	return "hello {{.Name}}!"
}

// Used as test data
type Person struct {
	Name string
}

func TestTemplateSrc(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	gt := &GoTemplateTestStruct{}
	gt.Name = "GoTemplateTestStruct"

	if gt.String() != "hello {{.Name}}!" {
		t.Errorf("Template source string is incorrect. gt.Src(): ", gt.String())
	}
}

func TestTemplateExecute(t *testing.T) {
	os.Setenv("CHILLAX_ENV", "test")

	gt := &GoTemplateTestStruct{}
	gt.Name = "GoTemplateTestStruct"
	gt.Src = gt.String()

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
