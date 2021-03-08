package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func InjectInstr(sl *ISymbol, instr string) error {
	offset := sl.Fset.Position(sl.Func.Body.Lbrace).Offset + 1

	topHalf := make([]byte, offset)
	copy(topHalf, sl.Content[:offset])

	instr = fmt.Sprintf("\n%s\n", instr)

	bottomHalf := make([]byte, len(sl.Content)-offset)
	copy(bottomHalf, sl.Content[offset:])

	var newContent []byte
	newContent = append(newContent, topHalf...)
	newContent = append(newContent, []byte(instr)...)
	newContent = append(newContent, bottomHalf...)

	sl.Content = newContent
	sl.Fset = token.NewFileSet()

	return nil
}

func InjectImports(sl *ISymbol, imports string) error {
	defer func() { GoFmt(sl.Fpath) }()

	var (
		parsedFile *ast.File
		err        error
	)

	if parsedFile, err = parser.ParseFile(sl.Fset, sl.Fpath, sl.Content, parser.ParseComments); err != nil {
		return err
	}

	// E.x.
	// --imports='fmt'
	// --imports='fmt, net/http'
	// --imports='f fmt, net/http'
	// --imports='f fmt, h net/http'
	// --imports='f fmt, h net/http, _ time'
	newImports := ""
	sp := strings.Split(imports, ",")

	for _, imp := range sp {
		_sp := strings.Split(imp, " ")

		if len(_sp) > 2 || len(_sp) < 1 {
			return fmt.Errorf("failed to inject import, its format is wrong: %s", imp)
		}

		if len(_sp) == 1 {
			newImports += fmt.Sprintf("\"%s\"\n", _sp[0])
			continue
		}

		newImports += fmt.Sprintf("%s \"%s\"\n", _sp[0], _sp[1])
	}

	if newImports != "" {
		newImports = fmt.Sprintf("import (\n %s)\n", newImports)
		offset := sl.Fset.Position(parsedFile.Name.End()).Offset + 1

		topHalf := make([]byte, offset)
		copy(topHalf, sl.Content[:offset])

		bottomHalf := make([]byte, len(sl.Content)-offset)
		copy(bottomHalf, sl.Content[offset:])

		var newContent []byte
		newContent = append(newContent, topHalf...)
		newContent = append(newContent, []byte(newImports)...)
		newContent = append(newContent, bottomHalf...)
		sl.Content = newContent
	}

	return ioutil.WriteFile(sl.Fpath, sl.Content, os.ModePerm)
}

// Depend on gofmt bin
func GoFmt(path string) error {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("gofmt -s -w %s", path))
	return cmd.Run()
}
