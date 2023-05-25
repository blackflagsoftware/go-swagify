package response

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
)

type (
	Response struct {
		Ref         string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string             `json:"description,omitempty" yaml:"description,omitempty"`
		Content     map[string]Content `json:"content,omitempty" yaml:"content,omitempty"`
	}

	Content struct {
		RefSchema `json:"schema" yaml:"schema"`
	}

	RefSchema struct {
		Ref string `json:"$ref" yaml:"$ref"`
	}
)

/*
	go-swagify

@@response: <name or status code>
@@ref: (optional) schema reference
@@desc: (required, if @@ref not used)
@@content_name: (not required if @@ref is used, else optional) application/json, etc
@@content_ref: (not required if @@ref is used, else optional) schema reference
... can repeat @@content_*
*/
func BuildResponse(comments in.SwagifyComment) map[string]Response {
	responses := make(map[string]Response)
	for name, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			response := &Response{Content: make(map[string]Content)}
			parseResponseLines(lines, response)
			blankOutRef(response)
			responses[name] = *response
		}
	}
	return responses
}

// called by operation
func ParseOperationResponseLines(lines []string) map[string]Response {
	responses := make(map[string]Response)
	response := &Response{Content: make(map[string]Content)}
	reg := regexp.MustCompile("(?P<name>[a-zA-Z_/.]+): *?(?P<value>.+)")
	currentResponseName := ""
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
		case "resp_name":
			if currentResponseName != value {
				if currentResponseName != "" {
					responses[currentResponseName] = *response
				}
				response = &Response{Content: make(map[string]Content)}
			}
			currentResponseName = value
		case "resp_ref":
			response.Ref = "#/components/responses/" + value
		}
	}
	if currentResponseName != "" {
		responses[currentResponseName] = *response
	}
	return responses
}

// called by the comments
func parseResponseLines(lines []string, response *Response) {
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
			response.Ref = "#/components/responses/" + value
		case "desc":
			response.Description = value
		case "content_name":
			if currentContentName != value {
				if currentContentName != "" {
					response.Content[currentContentName] = content
				}
				content = Content{}
			}
			currentContentName = value
		case "content_ref":
			content.Ref = "#/components/schemas/" + value
		}
	}
	if currentContentName != "" {
		response.Content[currentContentName] = content
	}
}

func blankOutRef(response *Response) {
	if response.Ref != "" {
		response.Content = make(map[string]Content)
		response.Description = ""
	}
}
