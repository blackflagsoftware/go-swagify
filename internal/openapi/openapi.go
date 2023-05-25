package openapi

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
	par "github.com/blackflagsoftware/go-swagify/internal/parameter"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
	pat "github.com/blackflagsoftware/go-swagify/internal/path"
	req "github.com/blackflagsoftware/go-swagify/internal/requestBody"
	res "github.com/blackflagsoftware/go-swagify/internal/response"
	sch "github.com/blackflagsoftware/go-swagify/internal/schema"
	sec "github.com/blackflagsoftware/go-swagify/internal/security"
	ser "github.com/blackflagsoftware/go-swagify/internal/server"
)

/* go-swagify
@@openapi: (required) 3.+
@@info.title: (requrired)
@@info.description: (optional)
@@info.termsOfService: (optional)
@@info.contact: (optional)
@@info.license: (optional)
@@info.version: (required)
*/

type (
	OpenApi struct {
		Version    string `json:"openapi" yaml:"openapi"`
		Info       `json:"info" yaml:"info"`
		Servers    []ser.Server          `json:"servers,omitempty" yaml:"servers,omitempty"`
		Paths      map[string]pat.Path   `json:"paths" yaml:"paths"`
		Components Component             `json:"components" yaml:"components"`
		Security   []map[string][]string `json:"security,omitempty" yaml:"security,omitempty"`
	}

	Info struct {
		Title          string `json:"title,omitempty" yaml:"title,omitempty"`
		Description    string `json:"description,omitempty" yaml:"description,omitempty"`
		TermsOfService string `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
		Contact        `json:"contact,omitempty" yaml:"contact,omitempty"`
		License        `json:"license,omittempty" yaml:"license,omitempty"`
		Version        string `json:"version" yaml:"version"`
	}

	Contact struct {
		Name  string `json:"name,omitempty" yaml:"name,omitempty"`
		Url   string `json:"url,omitempty" yaml:"url,omitempty"`
		Email string `json:"email,omitempty" yaml:"email,omitempty"`
	}

	License struct {
		Name string `json:"name,omitempty" yaml:"name,omitempty"`
		Url  string `json:"url,omitempty" yaml:"url,omitempty"`
	}

	Component struct {
		Parameters      map[string]par.Parameter      `json:"parameters" yaml:"parameters"`
		Schemas         map[string]sch.Schema         `json:"schemas" yaml:"schemas"`
		Responses       map[string]res.Response       `json:"responses" yaml:"responses"`
		RequestBodies   map[string]req.RequestBody    `json:"requestBodies" yaml:"requestBodies"`
		SecuritySchemes map[string]sec.SecurityScheme `json:"securitySchemes" yaml:"securitySchemes"`
	}
)

func BuildOpenApi(comments in.SwagifyComment) OpenApi {
	open := &OpenApi{Version: "3.0.0"}
	for _, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			err := parseOpenLines(lines, open)
			if err != nil {
				// nvever going to not nil
				perr.AddError(fmt.Sprintf("Unable to parse lines: %s", err))
				continue
			}
		}
	}
	return *open
}

func parseOpenLines(lines []string, open *OpenApi) error {
	reg := regexp.MustCompile("(?P<name>[a-zA-Z/.]+): *?(?P<value>.+)")
	info := Info{}
	contact := Contact{}
	license := License{}
	for _, line := range lines {
		matches := reg.FindStringSubmatch(line)
		nameIdx := reg.SubexpIndex("name")
		valueIdx := reg.SubexpIndex("value")
		if len(matches) < 2 {
			perr.AddError(fmt.Sprintf("[Warning] @@openapi: bad format for line: %s", line))
			continue
		}
		value := strings.TrimSpace(matches[valueIdx])
		switch matches[nameIdx] {
		case "info.title":
			info.Title = value
		case "info.description":
			info.Description = value
		case "info.termOfService":
			info.TermsOfService = value
		case "info.version":
			info.Version = value
		case "info.contact.name":
			contact.Name = value
		case "info.contact.url":
			contact.Url = value
		case "info.contact.email":
			contact.Email = value
		case "info.license.name":
			license.Name = value
		case "info.license.url":
			license.Url = value
		default:
			perr.AddError(fmt.Sprintf("[Warning] @@openapi: invalid name option: %s", line))
		}
	}
	open.Info = info
	open.Info.Contact = contact
	open.Info.License = license
	return nil
}
