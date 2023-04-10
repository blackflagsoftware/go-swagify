package operation

import (
	"fmt"
	"regexp"
	"strings"

	in "github.com/blackflagsoftware/go-swagify/internal"
	// srv "github.com/blackflagsoftware/go-swagify/internal/server"
	par "github.com/blackflagsoftware/go-swagify/internal/parameter"
	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
	req "github.com/blackflagsoftware/go-swagify/internal/requestBody"
	res "github.com/blackflagsoftware/go-swagify/internal/response"
)

type (
	// helper struct
	OperationBuild struct {
		Operations map[string]Operation
	}

	Operation struct {
		Summary     string             `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string             `json:"description,omitempty" yaml:"description,omitempty"`
		Parameters  []par.ParameterRef `json:"parameters,omitempty" yaml:"parameters,omitempty"`
		RequestBody req.ReqSchema      `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
		// Servers     []srv.Server    `json:"servers" yaml:"servers"`
		Response map[string]res.Response `json:"responses,omitempty" yaml:"responses,omitempty"`
	}
)

/* go-swagify
@@operation: <path url>
@@method: get|put|post|delete|options|head|patch|trace
@@summary: (optional)
@@description: (optional)
@@parameters.ref: (optional) semicolon(;) list of ref parameter names
@@resp_name: (required) 200, 300, 4xx, etc
@@resp_ref: (required) name of the response reference
... @@resp_name, resp_ref can repeat
*/
func BuildOperations(comments in.SwagifyComment) map[string]OperationBuild {
	operations := make(map[string]OperationBuild)
	operationBuild := &OperationBuild{Operations: make(map[string]Operation)}
	for name, lines := range comments.Comments {
		err := parseOperationLines(lines, operationBuild)
		if err != nil {
			// will never be not nil
			continue
		}
		operations[name] = *operationBuild
	}
	return operations
}

func parseOperationLines(lines []string, operationBuild *OperationBuild) error {
	operation := Operation{}
	// go through each line and do logic on
	reg := regexp.MustCompile("(?P<name>[a-zA-Z_/.]+): *?(?P<value>.+)")
	method := ""
lines_loop:
	for i, line := range lines {
		matches := reg.FindStringSubmatch(line)
		nameIdx := reg.SubexpIndex("name")
		valueIdx := reg.SubexpIndex("value")
		if len(matches) < 2 {
			perr.AddError(fmt.Sprintf("[Warning] @@operation: bad format of line: %s", line))
			continue
		}
		value := strings.TrimSpace(matches[valueIdx])
		switch matches[nameIdx] {
		case "method":
			method = value
		case "summary":
			operation.Summary = value
		case "description":
			operation.Description = value
		case "parameters.ref":
			split := strings.Split(value, ";")
			parameters := []par.ParameterRef{}
			for i := range split {
				parameters = append(parameters, par.ParameterRef{Ref: fmt.Sprintf("#/components/parameters/%s", split[i])})
			}
			operation.Parameters = parameters
		case "req_ref":
			operation.RequestBody = req.ReqSchema{Ref: fmt.Sprintf("#/components/requestBodies/%s", value)}
		case "resp_name":
			// hand off all the rest of the lines to responses
			operation.Response = res.ParseOperationResponseLines(lines[i:])
			break lines_loop
		default:
			perr.AddError(fmt.Sprintf("[Warning] @@operation: invalid name option: %s", line))
		}
	}
	if method == "" {
		perr.AddError(fmt.Sprintf("[Error] @@operation: no method specified"))
		return nil
	}
	validMethods := map[string]struct{}{"get": {}, "put": {}, "post": {}, "delete": {}, "options": {}, "head": {}, "patch": {}, "trace": {}}
	if _, ok := validMethods[method]; !ok {
		perr.AddError(fmt.Sprintf("[Warning] @@operation: invalid method: %s", method))
		return nil
	}
	operationBuild.Operations[method] = operation
	return nil
}
