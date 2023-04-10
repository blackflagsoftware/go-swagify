package parameter

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
	sch "github.com/blackflagsoftware/go-swagify/internal/schema"
)

type (
	Parameter struct {
		Name        string              `json:"name" yaml:"name"`
		In          string              `json:"in" yaml:"in"`
		Description string              `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool                `json:"required,omitempty" yaml:"required,omitempty"`
		Schema      *sch.SchemaProperty `json:"schema,omitempty" yaml:"schema,omitempty"`
	}

	ParameterRef struct {
		Ref string `json:"$ref" yaml:"$ref"`
	}
)

/* Parameter Sample
go-swagify
@@parameter: <name of parameter>
@@name: <match inline name>
@@in: query | header | path | cookie
@@description: (optional)
@@required: (optional) true | false(default)
// schema (optional) see schema.go/SchemaProperty
@@schema_type: string
@@schema_description: string
@@schema_example: interface{}
*/

func BuildParameters(comments in.SwagifyComment) map[string]Parameter {
	parameter := make(map[string]Parameter)
	for name, lines := range comments.Comments {
		Parameter, err := parseParameterLines(lines)
		if err != nil {
			// will never be not nil
			continue
		}
		parameter[name] = Parameter
	}
	return parameter
}

func parseParameterLines(lines []string) (Parameter, error) {
	var schemaProperty *sch.SchemaProperty
	Parameter := Parameter{}
	// go through each line and do logic on
	reg := regexp.MustCompile("(?P<name>[a-zA-Z_]+): *?(?P<value>.+)")
	lastName := ""
	for _, line := range lines {
		matches := reg.FindStringSubmatch(line)
		nameIdx := reg.SubexpIndex("name")
		valueIdx := reg.SubexpIndex("value")
		if len(matches) < 2 {
			perr.AddError(fmt.Sprintf("[Warning] @@parameter: bad format of line: %s", line))
			continue
		}
		lastName = matches[nameIdx]
		value := strings.TrimSpace(matches[valueIdx])
		if strings.Index(lastName, "schema") == 0 && schemaProperty == nil {
			schemaProperty = &sch.SchemaProperty{}
			Parameter.Schema = schemaProperty
		}
		switch matches[nameIdx] {
		case "name":
			Parameter.Name = value
		case "in":
			Parameter.In = value
		case "description":
			Parameter.Description = value
		case "required":
			if value == "true" {
				Parameter.Required = true
			}
		case "schema_type":
			schemaProperty.Type = value
		case "schema_description":
			schemaProperty.Description = value
		case "schema_example":
			// TODO: check for type and cast if appropriate, probably do this in schema.go
			schemaProperty.ExampleStr = value
		default:
			perr.AddError(fmt.Sprintf("[Warning] @@parameter: invalid name option: %s", line))
		}
	}
	Parameter.ValidateIn()
	return Parameter, nil
}

func (p *Parameter) ValidateIn() {
	if p.In == "" {
		perr.AddError(fmt.Sprintf("[Error] parameter In is required for %s", p.Name))
		return
	}
	validIn := map[string]struct{}{"query": {}, "header": {}, "path": {}, "cookie": {}}
	if _, ok := validIn[p.In]; !ok {
		perr.AddError(fmt.Sprintf("[Error] parameter In is invalid for %s; expected [query | header | path | cookie]", p.Name))
	}
}
