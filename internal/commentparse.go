package internal

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	perr "github.com/blackflagsoftware/go-swagify/internal/parseerror"
)

type (
	Component struct {
		Types map[string]SwagifyComment
	}

	SwagifyComment struct {
		Comments map[string][]string
	}
)

/* go-swagify
@@<type>: <name>
@@<name>: <value>
...
*/
// or
/* go-swagify
@@<type>: <name> @@<name>: <value> ...
*/

func ParseDirForComments(directory string) (comments []string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, directory, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error in parsing directory:", directory)
		return
	}
	for _, pa := range pkgs {
		for n, f := range pa.Files {
			if path.Ext(n) == ".go" {
				for _, c := range f.Comments {
					for _, l := range c.List {
						comments = append(comments, l.Text)
					}
				}
			}
		}
	}

	directoryItems, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Error in getting folder items for directory:", directory)
		return
	}
	for _, di := range directoryItems {
		if di.IsDir() {
			nextDirectory := path.Join(directory, di.Name())
			comments = append(comments, ParseDirForComments(nextDirectory)...)
		}
	}
	return
}

/*
this will return something like this:
{
map[parameter]: {map[<name>]: [line1, line2, ...], map[<name>]: [line1, line2, ...]}
map[schema]: {map[<name>]: [line1, line2, ...], map[<name>]: [line1, line2, ...]}
}
*/
func ParseSwagifyComment(comments []string) Component {
	reg := regexp.MustCompile("(?P<comp_type>[a-zA-Z]+): *?(?P<name>.+)")
	component := Component{Types: make(map[string]SwagifyComment)}
	for _, c := range comments {
		// check if the first 20 characters contain "go-swagify"
		l := len(c)
		if l > 20 {
			l = 20
		}
		check := string(c[:l])
		if idx := strings.Index(check, "go-swagify"); idx > -1 {
			// remove "/* go-swagify" and trailing "*/"
			start := idx + 10
			end := len(c) - 2
			justCommments := strings.TrimSpace(string(c[start:end]))
			splitComment := strings.Split(justCommments, "@@") // splitting this will give the frist index of "", just ignore
			cleanComment := strings.TrimSpace(splitComment[1])
			matches := reg.FindStringSubmatch(cleanComment)
			compTypeIdx := reg.SubexpIndex("comp_type")
			nameIdx := reg.SubexpIndex("name")
			if len(matches) < 2 {
				// no match
				perr.AddError(fmt.Sprintf("[Warning] bad format of line: %s", cleanComment))
				continue
			}
			compType := strings.TrimSpace(matches[compTypeIdx])
			name := strings.TrimSpace(matches[nameIdx])
			// save off the rest of the lines per map name
			_, ok := component.Types[matches[compTypeIdx]]
			if !ok {
				component.Types[compType] = SwagifyComment{Comments: make(map[string][]string)}
			}
			comments := component.Types[compType].Comments
			// remove \n and rejoin with @@
			cleanedComments := []string{}
			for i := 2; i < len(splitComment); i++ {
				cleanedComment := strings.ReplaceAll(splitComment[i], "\n", " ")
				cleanedComment = strings.TrimSpace(cleanedComment)
				cleanedComments = append(cleanedComments, cleanedComment)
			}
			comments[name] = cleanedComments
		}
	}
	return component
}
