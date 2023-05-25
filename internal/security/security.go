package security

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
)

type (
	SecurityScheme struct {
		Type        string `json:"type" yaml:"type"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
		Scheme      string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	}
)

/* go-swagify
@@security: <location>
@@name: string
@@scope: (optional) semicolon(;) list of scope names
repeat @@name

@@securityScheme: name
@@type: string
@@scheme: string
@@description: (optional)
*/

func BuildSecurity(comments in.SwagifyComment) map[string][]map[string][]string {
	reg := regexp.MustCompile("(?P<name>[a-zA-Z]+): *?(?P<value>.+)")
	securityMap := make(map[string][]map[string][]string)
	for name, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			security := make(map[string][]string)
			securityName := ""
			for _, line := range lines {
				matches := reg.FindStringSubmatch(line)
				nameIdx := reg.SubexpIndex("name")
				valueIdx := reg.SubexpIndex("value")
				if len(matches) < 2 {
					fmt.Println("parse securitySchema not formatted:", line)
					continue
				}
				value := strings.TrimSpace(matches[valueIdx])
				switch matches[nameIdx] {
				case "name":
					securityName = value
					security[securityName] = []string{}
				case "scope":
					split := strings.Split(value, ";")
					security[securityName] = split
				default:
					perr.AddError(fmt.Sprintf("[Warning] @@security: invalid name option: %s", line))
				}
			}
			securityMap[name] = append(securityMap[name], security)
		}
	}
	return securityMap
}

func BuildSecuritySchemes(comments in.SwagifyComment) map[string]SecurityScheme {
	reg := regexp.MustCompile("(?P<name>[a-zA-Z]+): *?(?P<value>.+)")
	securitySchemeMap := make(map[string]SecurityScheme)
	securityScheme := SecurityScheme{}
	for name, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			for _, line := range lines {
				matches := reg.FindStringSubmatch(line)
				nameIdx := reg.SubexpIndex("name")
				valueIdx := reg.SubexpIndex("value")
				if len(matches) < 2 {
					fmt.Println("parse securitySchema not formatted:", line)
					continue
				}
				value := strings.TrimSpace(matches[valueIdx])
				switch matches[nameIdx] {
				case "type":
					securityScheme.Type = value
				case "scheme":
					securityScheme.Scheme = value
				case "description":
					securityScheme.Description = value
				default:
					perr.AddError(fmt.Sprintf("[Warning] @@schema: invalid name option: %s", line))
				}
			}
			securitySchemeMap[name] = securityScheme
			securityScheme = SecurityScheme{}
		}
	}
	return securitySchemeMap
}
