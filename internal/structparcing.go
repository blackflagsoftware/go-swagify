package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
)

type (
	MyStruct struct {
		Name   string
		Fields []MyField
	}

	MyField struct {
		Name string
		Type string
		Tag  string
	}
)

/*
this will take a list of "marked" (through comments) to parse and create a MyStruct structure
based on the struct's content, used by "schema", see internal/schema/schema.go
*/
func ParseDirForStructs(directory string, comments SwagifyComment) (myStructs []MyStruct) {
	// only want to deal with .go files
	dirItems, errReadDir := os.ReadDir(directory)
	if errReadDir != nil {
		fmt.Println("Error reading diretory:", directory)
		return
	}
	for _, di := range dirItems {
		file := path.Join(directory, di.Name())
		if di.IsDir() {
			myStructs = append(myStructs, ParseDirForStructs(file, comments)...)
			continue
		}
		if path.Ext(di.Name()) != ".go" {
			continue
		}
		src, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("readfile:", err)
			break
		}
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, file, src, 0)
		if err != nil {
			fmt.Println("Error in parsing file:", file)
			return
		}
		ast.Inspect(parsedFile, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.TypeSpec:
				if s, ok := t.Type.(*ast.StructType); ok {
					if _, foundStruct := comments.Comments[t.Name.Name]; foundStruct {
						myStruct := MyStruct{Name: t.Name.Name}
						for _, field := range s.Fields.List {
							if len(field.Names) > 0 && field.Tag != nil {
								myStruct.Fields = append(myStruct.Fields, MyField{
									Name: field.Names[0].Name,
									Type: string(src[field.Type.Pos()-1 : field.Type.End()-1]),
									Tag:  string(src[field.Tag.Pos() : field.Tag.End()-2]),
								})
							}
						}
						myStructs = append(myStructs, myStruct)
					}
				}
			}
			return true
		})
	}
	return
}
