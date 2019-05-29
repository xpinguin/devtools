// dumpfmt project main.go
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/iancoleman/strcase"

	"go/parser"

	"golang.org/x/tools/go/ast/astutil"
	_ "golang.org/x/tools/go/ast/astutil"
)

type FmtString string

func (s *FmtString) IsFmt() bool {
	return strings.Contains(string(*s), "%")
}

type FmtCall struct {
	recvTy string

	fmt  FmtString
	args []string
	name string

	closure       string
	closureParams []string

	recv string
}

////
func loadFile(fname string) (fset *token.FileSet, r *ast.File, text string) {
	fset = token.NewFileSet()
	r, err := parser.ParseFile(fset, fname, nil, parser.AllErrors)
	if err != nil {
		log.Fatal("Failed to parse file: ", err)
	}

	data, _ := ioutil.ReadFile(fname)
	text = string(data)
	return fset, r, text
}

////
func nodeSrc(fset *token.FileSet, ftext string, n ast.Node) string {
	pos := fset.Position(n.Pos())
	endpos := fset.Position(n.End())
	return ftext[pos.Offset:endpos.Offset]
}

func posSrc(fset *token.FileSet, ftext string, s, e token.Pos) string {
	return ftext[fset.Position(s).Offset : fset.Position(e).Offset+1]
}

////
func collectFmtCalls(stx ast.Node, fset *token.FileSet, ftext string) (calls []FmtCall) {
	var fmtCall FmtCall

	astutil.Apply(
		stx,
		func(c *astutil.Cursor) bool {
			switch x := c.Node().(type) {
			case *ast.FuncDecl:
				fmtCall.closure = x.Name.Name
				if x.Recv != nil && len(x.Recv.List) > 0 {
					if fmtCall.closure != "String" {
						return false
					}
					fmtCall.recvTy = strcase.ToSnake(x.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name)
					if ns := x.Recv.List[0].Names; len(ns) > 0 {
						fmtCall.recv = ns[0].Name
					}
				}
				for _, par := range x.Type.Params.List {
					fmtCall.closureParams = append(fmtCall.closureParams, nodeSrc(fset, ftext, par))
				}

			case *ast.CallExpr:
				fmtCall.name = nodeSrc(fset, ftext, x.Fun)
				fmtCall.args = []string{}
				fmtCall.fmt = ""
				for _, arg := range x.Args {
					argStr := strings.TrimSpace(nodeSrc(fset, ftext, arg))
					if s := FmtString(argStr); fmtCall.fmt == "" && s.IsFmt() {
						fmtCall.fmt = s
					} else if fmtCall.fmt != "" {
						argStr = strings.ReplaceAll(argStr, fmtCall.recv+".", "|_|.")
						argStr = strings.ReplaceAll(argStr, ", "+fmtCall.recv+",", ", |_|,")
						argStr = strings.ReplaceAll(argStr, ", "+fmtCall.recv+")", ", |_|)")
						argStr = strings.ReplaceAll(argStr, "("+fmtCall.recv+",", "(|_|,")

						fmtCall.args = append(fmtCall.args, argStr)
					}
				}
				if fmtCall.fmt != "" {
					calls = append(calls, fmtCall)
					return false
				}
			}
			return true
		},
		func(c *astutil.Cursor) bool {
			switch c.Node().(type) {
			case *ast.FuncDecl:
				fmtCall = FmtCall{}
			}
			return true
		},
	)
	return calls
}

////
func main() {
	////
	srcpath := flag.String("file", `C:\_go\wrk\src\golang.org\x\tools\go\ssa\print.go`, "Go source file")
	flag.Parse()

	if srcpath == nil || *srcpath == "" {
		log.Fatal("No file specified")
		return
	}

	////
	fset, stx, flines := loadFile(*srcpath)

	scopes := map[string]FmtCall{}
	cnt, total := 0, 0
	for _, c := range collectFmtCalls(stx, fset, flines) {
		scopeName := c.closure
		if c.recvTy != "" {
			scopeName = c.recvTy + "." + scopeName
			///
			fmt.Printf("%s(%s) .. %s(%s)\n\n",
				c.recvTy, strings.TrimSpace(string(c.fmt)),
				c.recvTy, strings.Join(c.args, ", "))
			cnt++
		}
		total++
		///
		if _, ok := scopes[scopeName]; !ok {
			scopes[scopeName] = c
		}
	}

	fmt.Println("UNIQUE: ", len(scopes))
	fmt.Println("WITH-RECEIVER: ", cnt)
	fmt.Println("TOTAL: ", total)

	return

}
