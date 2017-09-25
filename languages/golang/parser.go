package golang

import (
	"fmt"
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/util"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var (
	pkgRegex = regexp.MustCompile(`\A\s*package\s*([a-zA-Z0-9_]+)\s*\z`)
)

type StructVisitor struct {
	src     string
	structs []*base.Struct
}

func (sv *StructVisitor) Visit(node ast.Node) ast.Visitor {
	switch t := node.(type) {
	case *ast.GenDecl:
		if t.Tok == token.TYPE { // Type declaration statement
			for _, spec := range t.Specs { // Items inside type
				tSpec := spec.(*ast.TypeSpec)
				sType, ok := tSpec.Type.(*ast.StructType) // Struct type
				if !ok {
					break
				}
				s := &base.Struct{ // Parse new struct
					Name:   tSpec.Name.Name,
					Fields: make([]*base.Field, 0),
					Info: base.Info{
						Text: fmt.Sprintf("type %s %s", tSpec.Name.Name, sv.src[sType.Pos()-1:sType.End()-1]),
					},
				}

				for _, f := range sType.Fields.List { // Struct fields
					bf := &base.Field{}
					for _, name := range f.Names { // Field name
						bf.Name += name.Name + " "
					}
					bf.Name = strings.TrimSpace(bf.Name)
					bf.Type = strings.TrimSpace(sv.src[f.Type.Pos()-1 : f.Type.End()-1]) // Field type
					s.Fields = append(s.Fields, bf)
				}

				sv.structs = append(sv.structs, s)
			}
		}
	}
	return sv
}

func parseStructs(data string) ([]*base.Struct, error) {
	if !isContainPackageStmt(data) {
		data = fmt.Sprintf("package fugu\n%s", data)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, 0)
	if err != nil {
		return nil, err
	}

	s := &StructVisitor{
		src: data,
	}
	ast.Walk(s, f)
	return s.structs, nil
}

func isContainPackageStmt(data string) bool {
	lines := strings.Split(data, util.LineBreak())
	for _, line := range lines {
		if pkgRegex.MatchString(line) {
			return true
		}
	}
	return false
}
