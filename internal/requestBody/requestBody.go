package requestBody

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
)

type (
	RequestBody struct {
		Ref         string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string             `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool               `json:"required,omitempty" yaml:"required,omitempty"`
		Content     map[string]Content `json:"content,omitempty" yaml:"content,omitempty"`
	}

	Content struct {
		ReqSchema `json:"schema" yaml:"schema"`
	}

	ReqSchema struct {
		Ref string `json:"$ref" yaml:"$ref"`
	}
)

/*
	go-swagify

@@requestBody: <name or status code>
@@ref: (optional) schema reference
@@desc: (required, if @@ref not used)
@@required: (optional) true/false
@@content_name: (not required if @@ref is used, else optional) application/json, etc
@@content_ref: (not required if @@ref is used, else optional) schema reference
*/
func BuildRequestBody(comments in.SwagifyComment) map[string]RequestBody {
	requestBodies := make(map[string]RequestBody)
	for name, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			requestBody := &RequestBody{Content: make(map[string]Content)}
			parseRequestBodyLines(lines, requestBody)
			blankOutRef(requestBody)
			requestBodies[name] = *requestBody
		}
	}
	return requestBodies
}

// called by the comments
func parseRequestBodyLines(lines []string, requestBody *RequestBody) {
	content := Content{}
	reg := regexp.MustCompile("(?P<name>[a-zA-Z_/.]+): *?(?P<value>.+)")
	currentContentName := ""
	for _, line := range lines {
		matches := reg.FindStringSubmatch(line)
		nameIdx := reg.SubexpIndex("name")
		valueIdx := reg.SubexpIndex("value")
		if len(matches) < 2 {
			perr.AddError(fmt.Sprintf("[Warning] @@responses: bad format of line: %s", line))
			continue
		}
		value := strings.TrimSpace(matches[valueIdx])
		switch matches[nameIdx] {
		case "ref":
			requestBody.Ref = "#/components/requestBodies/" + value
		case "desc":
			requestBody.Description = value
		case "required":
			if value == "true" {
				requestBody.Required = true
			}
		case "content_name":
			if currentContentName != value {
				if currentContentName != "" {
					requestBody.Content[currentContentName] = content
				}
				content = Content{}
			}
			currentContentName = value
		case "content_ref":
			content.Ref = "#/components/schemas/" + value
		}
	}
	if currentContentName != "" {
		requestBody.Content[currentContentName] = content
	}
}

func blankOutRef(requestBody *RequestBody) {
	if requestBody.Ref != "" {
		requestBody.Content = make(map[string]Content)
		requestBody.Description = ""
	}
}
