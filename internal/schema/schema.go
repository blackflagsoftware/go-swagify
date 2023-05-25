package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/blackflagsoftware/go-swagify/config"
	in "github.com/blackflagsoftware/go-swagify/internal"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
	"github.com/blackflagsoftware/go-swagify/internal/util"
	"github.com/fatih/structtag"
)

type (
	Schema struct {
		Type           string                    `json:"type,omitempty" yaml:"type,omitempty"`
		Required       []string                  `json:"required,omitempty" yaml:"required,omitempty"`
		Description    string                    `json:"description,omitempty" yaml:"description,omitempty"`
		Example        string                    `json:"example,omitempty" yaml:"example,omitempty"`
		Properties     map[string]SchemaProperty `json:"properties,omitempty" yaml:"properties,omitempty"`
		AddlProperties AdditionalProperty        `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
		Items          map[string]string         `json:"items,omitempty" yaml:"items,omitempty"`
	}

	// TODO: if type is object or array may need to have self reference
	SchemaProperty struct {
		Ref         string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Type        string      `json:"type,omitempty" yaml:"type,omitempty"`
		Description string      `json:"description,omitempty" yaml:"description,omitempty"`
		Example     interface{} `json:"example,omitempty" yaml:"example,omitempty"`
		ExampleStr  string      `json:"-" yaml:"-"`
		Enum        []string    `json:"enum,omitempty" yaml:"enum,omitempty"`
	}

	AdditionalProperty struct {
		Type  string            `json:"type,omitempty" yaml:"type,omitempty"`
		Items map[string]string `json:"items,omitempty" yaml:"items,omitempty"`
	}
)

/* go-swagify
@@schema: <name>
@@type: (required) [object | array]
@@prop_name: <name> (not needed with type => array)
@@prop_ref: <schema ref>
@@prop_req: (optional) add to the list of required in Schema; if false just leave omit
@@prop_type: (string) [object | array | string | number | etc]
@@prop_desc: (optional)
@@prop_ex: (optional)
repeat @@prop_* for object
*/

// or...

/* Schema Sample
// comment included above struct
// without this the parser will not see this struct to be parsed
go-swagify
@@struct: <name to match struct name>

// struct tag for each field
// need either json, yaml or both to match your expected use-case output

// "sw"
// for multiple schema names, separate with ';'
// if it is required for that schema name, append '*'
// required
sw:"ExampleRequest*;ExampleResponse"

// "sw_desc"
// desciption for the field name
// optional; will use the field name, lower case if not present
sw_desc:"here is a summary"

// "sw_ex"
// example for the field name
// optional; will use the field name, lower case if not present
sw_ex:"some example here"
*/

func BuildSchema(comments in.SwagifyComment, schemas map[string]Schema) {
	for name, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			schema := parseSchemaLines(lines)
			schemas[name] = schema
		}
	}
	return
}

func BuildSchemaStruct(myStructs []in.MyStruct) map[string]Schema {
	schemas := make(map[string]Schema)
	for _, m := range myStructs {
		for _, f := range m.Fields {
			parseTag(f.Name, f.Type, f.Tag, schemas)
		}
	}
	return schemas
}

func parseSchemaLines(lines []string) Schema {
	schema := Schema{Properties: make(map[string]SchemaProperty), Items: make(map[string]string)}
	// go through each line and do logic on
	reg := regexp.MustCompile("(?P<name>[a-zA-Z_/.]+): *?(?P<value>.+)")
	currentPropertyName := ""
	schemaProperty := SchemaProperty{}
	for _, line := range lines {
		matches := reg.FindStringSubmatch(line)
		nameIdx := reg.SubexpIndex("name")
		valueIdx := reg.SubexpIndex("value")
		if len(matches) < 2 {
			perr.AddError(fmt.Sprintf("[Warning] @@schema: bad format of line: %s", line))
			continue
		}
		value := strings.TrimSpace(matches[valueIdx])
		switch matches[nameIdx] {
		case "type":
			schema.Type = validateType(value)
		case "desc":
			schema.Description = value
		case "ex":
			schema.Example = value
		case "prop_name":
			if currentPropertyName != value {
				if currentPropertyName != "" {
					// save the current one and start fresh
					if schema.Type == "array" {
						schema.Items["$ref"] = schemaProperty.Ref
					} else {
						schema.Properties[currentPropertyName] = schemaProperty
					}
					schemaProperty = SchemaProperty{}
				}
			}
			currentPropertyName = value
		case "prop_ref":
			schemaProperty.Ref = "#/components/schemas/" + value
			if schema.Type == "array" {
				schema.Items["$ref"] = schemaProperty.Ref
			} else {
				schema.Properties[currentPropertyName] = schemaProperty
			}
		case "prop_type":
			schemaProperty.Type = value
		case "prop_req":
			if value == "true" {
				// just in case they add it with false or anything else
				schema.Required = append(schema.Required, currentPropertyName)
			}
		case "prop_desc":
			schemaProperty.Description = value
		case "prop_ex":
			schemaProperty.Example = exampleConv(schemaProperty.Type, value)
		case "addl_prop_ref":
			ref := "#/components/schemas/" + value
			schema.AddlProperties = AdditionalProperty{Type: "array", Items: map[string]string{"$ref": ref}} // TODO: this is only used to handle a map[string]array
		default:
			perr.AddError(fmt.Sprintf("[Warning] @@schema: invalid name option: %s", line))
		}
	}
	if schema.Type == "array" {
		schema.Items["$ref"] = schemaProperty.Ref
	} else {
		if currentPropertyName != "" {
			schema.Properties[currentPropertyName] = schemaProperty
		}
	}
	blankOutRef(&schema)
	return schema
}

func blankOutRef(schema *Schema) {
	for _, prop := range schema.Properties {
		if prop.Ref != "" {
			prop.Type = ""
			prop.Description = ""
			prop.Example = nil
		}
	}
}

func parseTag(fieldName, fieldType, fieldTag string, schemas map[string]Schema) {
	tags, err := structtag.Parse(fieldTag)
	if err != nil {
		fmt.Println("parseTag", err)
		return
	}
	lowerCaseFieldName := determineFieldName(fieldName, tags)
	sw, err := tags.Get("sw")
	if err != nil {
		// unable to find sw tag, ignore field
		return
	}
	// split possible schema names; i.e. ExampleRequest*;ExampleResponse => [ExampleRequest*, ExampleResponse]
	schemaNames := strings.Split(sw.Name, ";")
	for _, schemaName := range schemaNames {
		name, required := determineRequired(schemaName)
		if _, ok := schemas[name]; !ok {
			schemas[name] = Schema{Type: "object", Required: []string{}, Properties: make(map[string]SchemaProperty)}
		}
		if required {
			schema := schemas[name]
			schema.Required = append(schema.Required, lowerCaseFieldName)
			schemas[name] = schema
		}
		var example interface{}
		docType, desc, ref := "", "", ""
		swRef, errRef := tags.Get("sw_ref")
		if swRef == nil || swRef.Value() == "" || errRef != nil {
			docType, desc, example = parseSwagifyTag(fieldName, fieldType, tags)
		}
		if swRef != nil {
			ref = "#/components/schemas/" + swRef.Value()
		}
		schemaProperty := SchemaProperty{
			Ref:         ref,
			Type:        docType,
			Description: desc,
			Example:     example,
		}
		schemas[name].Properties[lowerCaseFieldName] = schemaProperty
	}
}

func determineRequired(schemaName string) (string, bool) {
	if schemaName[len(schemaName)-1:] == "*" {
		return schemaName[:len(schemaName)-1], true
	}
	return schemaName, false
}

func determineFieldName(fieldName string, tags *structtag.Tags) string {
	altFieldName := util.BuildAlternateFieldName(fieldName, config.AltFieldFormat)
	// name := strings.ToLower(altFieldName)
	tag, err := tags.Get(config.AppOutputFormat)
	if err != nil {
		return altFieldName
	}
	return tag.Name
}

func parseSwagifyTag(fieldName, fieldType string, tags *structtag.Tags) (docType string, desc string, example interface{}) {
	docType = "string"
	switch fieldType {
	case "float32", "float64", "null.Float":
		docType = "number"
	case "int", "int32", "int64", "null.Int":
		docType = "integer"
	case "bool", "null.Bool":
		docType = "boolean"
	}
	if swDesc, errDesc := tags.Get("sw_desc"); errDesc != nil {
		if jsonDesc, err := tags.Get(config.OutputFormat); err != nil {
			desc = strings.ToLower(fieldName)
		} else {
			desc = jsonDesc.Name
		}
	} else {
		desc = swDesc.Name
	}
	if swEx, errEx := tags.Get("sw_ex"); errEx != nil {
		if jsonEx, err := tags.Get(config.OutputFormat); err != nil {
			example = strings.ToLower(fieldName)
		} else {
			example = jsonEx.Name
		}
	} else {
		example = exampleConv(docType, swEx.Name)
		if len(swEx.Options) > 0 && docType == "string" && swEx.Options[0] != "omitempty" {
			// used as an example and may have a ',' in the string
			example = example.(string) + ", " + strings.TrimSpace(strings.Join(swEx.Options, ", "))
		}
	}
	return
}

func exampleConv(docType string, exampleStr string) (example interface{}) {
	example = exampleStr
	switch docType {
	case "number":
		floatEx, err := strconv.ParseFloat(exampleStr, 64)
		if err != nil {
			perr.AddError(fmt.Sprintf("[Warning] @@schema: unable to cast example: %s to float\n", exampleStr))
			floatEx = 0.0
		}
		example = floatEx
	case "integer":
		intEx, err := strconv.Atoi(exampleStr)
		if err != nil {
			perr.AddError(fmt.Sprintf("[Warning] @@schema: unable to cast example: %s to int\n", exampleStr))
			intEx = 0
		}
		example = intEx
	case "boolean":
		example = "true | false"
	}
	return
}

func validateType(t string) string {
	types := map[string]struct{}{"string": {}, "number": {}, "interger": {}, "boolean": {}, "object": {}, "array": {}}
	if _, ok := types[t]; !ok {
		perr.AddError(fmt.Sprintf("[Error] @@schema: invalid type: %s", t))
		return "string"
	}
	return t
}
