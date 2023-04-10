package parseerror

import (
	"fmt"
)

var parseErrors []string

func init() {
	parseErrors = []string{}
}

func AddError(e string) {
	parseErrors = append(parseErrors, e)
}

func PrintErrors() {
	if len(parseErrors) > 0 {
		fmt.Println("Messages while parsing")
		for _, e := range parseErrors {
			fmt.Printf("\t%s\n", e)
		}
	}
}
