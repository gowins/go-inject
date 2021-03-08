package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"strings"
)

const (
	IFunc int = iota
	IMethod
)

type ISymbol struct {
	Pkg     *Package
	Func    *ast.FuncDecl
	Fset    *token.FileSet
	Fpath   string
	typ     int
	Content []byte
}

func WalkPackages(pkgs map[string]*Package, pkgPath string) *Package {
	for _, pkg := range pkgs {
		if pkg.ImportPath == pkgPath {
			return pkg
		}
	}
	return nil
}

func WalkSymbols(pkg *Package, name string, typ int) (*ISymbol, error) {
	var (
		content    []byte
		err        error
		parsedFile *ast.File
	)

	for _, gFile := range pkg.GoFiles {

		fset := token.NewFileSet()
		fPath := path.Join(pkg.Dir, gFile)

		if content, err = ioutil.ReadFile(fPath); err != nil {
			return nil, err
		}

		if parsedFile, err = parser.ParseFile(fset, fPath, content, parser.ParseComments); err != nil {
			return nil, err
		}

		for _, v := range parsedFile.Decls {
			fn, ok := v.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if typ == IFunc && fn.Name.Name != name {
				continue
			}

			if typ == IMethod {
				// receiver (methods); or nil (functions)
				if fn.Recv == nil {
					continue
				}

				// E.x. func (*HelloHandler) HelloWorld(_ *gin.Context) pkg.Render
				// name = "${receiver} ${method}"
				// name = "HelloHandler HelloWorld"
				sp := strings.Split(name, " ")
				if len(sp) != 2 {
					continue
				}

				// E.x. func (h *HelloHandler) HelloWorld(_ *gin.Context) pkg.Render
				// fn.Recv.List[0] => "h"
				// fn.Recv.List[1] => "HelloHandler"
				se := fn.Recv.List[len(fn.Recv.List)-1].Type.(*ast.StarExpr)
				ident := se.X.(*ast.Ident)
				if ident.Name != sp[0] {
					continue
				}
			}

			return &ISymbol{pkg, fn, fset, fPath, typ, content}, nil
		}
	}

	return nil, fmt.Errorf("faild: No such symbol")
}

func WalkFuncs(pkg *Package, funcName string) (*ISymbol, error) {
	return WalkSymbols(pkg, funcName, IFunc)
}

func WalkMethods(pkg *Package, methodName string) (*ISymbol, error) {
	return WalkSymbols(pkg, methodName, IMethod)
}
