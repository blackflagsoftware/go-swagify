package util

import (
	"bytes"
	"fmt"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig/v3"
)

type (
	Format struct {
		Name string
	}
)

// name => to be formatted
// mode => what to change to: snakeCase, kebabCase, pascalCase, camelCase, lowerCase, upperCase
func BuildAlternateFieldName(name, mode string) string {
	f := Format{Name: name}
	var t *template.Template
	var err error
	switch mode {
	case "snakeCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | snakecase}}")
	case "kebabCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | kebabcase}}")
	case "camelCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | camelcase}}")
	case "pascalCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | camelcase}}")
	case "upperCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | upper}}")
	default:
		// lowerCase
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | lower}}")
	}
	if err != nil {
		fmt.Println(err)
		return name
	}
	b := bytes.NewBufferString("")
	errE := t.Execute(b, f)
	if errE != nil {
		fmt.Println(errE)
		return name
	}
	if mode == "camelCase" {
		// finish off the camel case functionality
		n := []rune(b.String())
		return string(append([]rune{unicode.ToLower(n[0])}, n[1:]...))
	}
	return b.String()
}
