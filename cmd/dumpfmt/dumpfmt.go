// dumpfmt project main.go
package main

import (
	"flag"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"go/parser"

	"github.com/davecgh/go-spew/spew"
)

type FmtString string

func (s *FmtString) IsFmt() bool {
	return strings.Contains(string(*s), "%")
}

func loadFile(fname string) (fset *token.FileSet, r *ast.File) {
	fset = token.NewFileSet()
	r, err := parser.ParseFile(fset, fname, nil, parser.AllErrors)
	if err != nil {
		log.Fatal("Failed to parse file: ", err)
	}
	return fset, r
}

func main() {
	////
	srcpath := flag.String("file", `C:\_go\wrk\src\golang.org\x\tools\go\ssa\print.go`, "Go source file")
	flag.Parse()

	if srcpath == nil || *srcpath == "" {
		log.Fatal("No file specified")
		return
	}

	////
	_, stx := loadFile(*srcpath)

	typeFmts := map[string][]FmtString{}
	var typeName, funcName string

	ast.Inspect(stx, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcName = x.Name.Name
			if x.Recv == nil || len(x.Recv.List) == 0 {
				typeName = funcName
				break
			}

			switch r := x.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				typeName = r.X.(*ast.Ident).Name
			default:
				log.Printf("Unk receiver: %#v", typeName)
				typeName = ""
				funcName = ""
			}

		case *ast.BasicLit:
			if x.Kind != token.STRING {
				break
			}

			s := FmtString(x.Value)
			if !s.IsFmt() && funcName != "String" {
				break
			}
			if typeName != "" {
				fmts := typeFmts[typeName]
				if fmts == nil {
					fmts = []FmtString{}
				}
				typeFmts[typeName] = append(fmts, s)
			}
		}
		return true
	})
	spew.Dump(typeFmts)
}
