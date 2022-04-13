//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"sort"
)

func main() {
	generateOpenapiErrorCodes(findErrorCodes())
}

// findErrorCodes will look for all declarations of variables
// of type AppError in the file where the generator is executing, such as:
// var ErrorInternal = AppError{
//		message:    "Unexpected server error",
//		code:       "internal_error",
//		statusCode: http.StatusInternalServerError,
//  }
func findErrorCodes() map[string]string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()

	// parse the Go code in the file annotated with a go:generate comment
	file, err := parser.ParseFile(fset, path.Join(cwd, os.Getenv("GOFILE")), nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	errs := map[string]string{}

	// iterate top level declarations
	for _, d := range file.Decls {
		switch decl := d.(type) {
		// we're interested in "generic declaration" nodes, which represents variable declaration among other things
		case *ast.GenDecl:
			// iterate over Specs, which are declarations
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				// we're interested in ValueSpec nodes (represent constant or variable declaration)
				case *ast.ValueSpec:
					// checking if the node is a composite literal of type AppError (AppError{...})
					if len(spec.Values) == 1 && spec.Values[0].(*ast.CompositeLit).Type.(*ast.Ident).Name == "AppError" {
						var code string
						var message string

						// iterate over this composite literal elements (the struct fields: message, code, statusCode)
						for _, el := range spec.Values[0].(*ast.CompositeLit).Elts {
							kv, ok := el.(*ast.KeyValueExpr)
							if !ok {
								continue
							}
							// extract field name
							key := kv.Key.(*ast.Ident).Name

							// extract the value for code and message
							switch key {
							case "code":
								code = kv.Value.(*ast.BasicLit).Value
							case "message":
								message = kv.Value.(*ast.BasicLit).Value
							}
						}

						errs[code] = message
					}
				}
			}
		}
	}

	return errs
}

func stripQuotes(s string) string {
	return s[1 : len(s)-1]
}

func generateOpenapiErrorCodes(codes map[string]string) {
	openapiFile := "openapi.yaml"

	f, err := os.Open(openapiFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	const startComment = "    <!-- ERROR_GENERATOR_START -->"
	const endComment = "    <!-- ERROR_GENERATOR_END -->"

	// skip will be used to signal if the program should copy the lines to the new buffer or not
	skip := false
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// found ERROR_GENERATOR_START: add this line and print the errors
		if scanner.Text() == startComment {
			skip = true

			// write the comment line
			_, err = buf.WriteString(startComment + "\n")
			if err != nil {
				panic(err)
			}

			// extract keys for sorting
			keys := make([]string, 0, len(codes))
			for k := range codes {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// iterate sorted keys and get value from the codes map
			for _, key := range keys {
				to := fmt.Sprintf("\n    `%s`: %s\n", stripQuotes(key), stripQuotes(codes[key]))
				_, err = buf.WriteString(to)
				if err != nil {
					panic(err)
				}
			}
		}

		// reached ERROR_GENERATOR_END: stop skipping lines
		if scanner.Text() == endComment {
			skip = false
		}

		// if skip == true, we're iterating the lines between ERROR_GENERATOR_START and ERROR_GENERATOR_END
		if !skip {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				panic(err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				panic(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// write the new content back to the spec file
	err = os.WriteFile(openapiFile, buf.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}
