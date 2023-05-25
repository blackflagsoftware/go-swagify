package server

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
)

type (
	Server struct {
		Url         string `json:"url" yaml:"url"`
		Description string `json:"description" yaml:"description,omitempty"`
		// TODO:
		// Variables
	}
)

/* go-swagify
@@server: <location>
@@url: (reguired)
@@description: (optional)
@@url: (required)
@@description: (optional)
*/

func BuildServers(comments in.SwagifyComment) map[string][]Server {
	reg := regexp.MustCompile("(?P<name>[a-zA-Z]+): *?(?P<value>.+)")
	serverMap := make(map[string][]Server)
	for name, lineArray := range comments.Comments {
		for _, lines := range lineArray {
			serverMap[name] = []Server{}
			server := Server{}
			foundFirst := false
			for _, line := range lines {
				matches := reg.FindStringSubmatch(line)
				nameIdx := reg.SubexpIndex("name")
				valueIdx := reg.SubexpIndex("value")
				if len(matches) < 2 {
					fmt.Println("parse server not formatted:", line)
					continue
				}
				value := strings.TrimSpace(matches[valueIdx])
				if matches[nameIdx] == "url" {
					if foundFirst {
						serverMap[name] = append(serverMap[name], server)
					}
					server = Server{Url: value}
					foundFirst = true
				}
				if matches[nameIdx] == "description" {
					server.Description = value
				}
			}
			serverMap[name] = append(serverMap[name], server)
		}
	}
	return serverMap
}
