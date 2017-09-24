package golang

import (
	"encoding/json"
	"fmt"
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/global"
	"go/ast"
	"go/token"
	"sort"
	"strings"
	"go/parser"
	"github.com/lnquy/fugu/modules/util"
	"regexp"
)

var (
	pkgRegex = regexp.MustCompile(`\A\s*package\s*([a-zA-Z0-9_]+)\s*\z`)
)

type (
	Golang struct{}

	StructVisitor struct {
		src     string
		structs []*base.Struct
	}
)

func (g *Golang) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	s, err := parseStructs(data)
	if err != nil {
		return "", err
	}
	for _, v := range s {
		calcFieldSizes(v, arch)
		calcPadding(v, arch)
		v.CalcOptimizable(arch)
	}

	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (g *Golang) OptimizeMemoryAlignment(s *base.Struct, arch global.Architecture) (string, error) {
	os := &base.Struct{
		Name: s.Name,
	}
	chunk := arch.GetChunkSize()
	optm := make([]*base.Field, 0)

	for _, f := range s.Fields {
		if f.Size%chunk == 0 && f.Size >= chunk {
			os.Fields = append(os.Fields, f)
			continue
		}
		optm = append(optm, f)
	}

	// TODO: Not solve the special case when (size > chunk && size % chunk != 0) yet!
	sort.Sort(base.BySize(optm))
	os.Fields = append(os.Fields, optm...)
	calcPadding(os, arch)
	os.CalcOptimizable(arch)
	os.BuildText()

	b, err := json.Marshal(os)
	if err != nil {
		return "", err
	}
	return string(b), nil
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

func calcFieldSizes(s *base.Struct, arch global.Architecture) {
	for _, f := range s.Fields {
		if f.Type == "struct{}" || strings.HasPrefix(f.Type, "[0]") {
			f.Size = 0
			continue
		}
		if strings.Contains(f.Type, "8") || f.Type == "bool" || f.Type == "byte" {
			f.Size = uint8(1)
			continue
		}
		if strings.Contains(f.Type, "16") {
			f.Size = uint8(2)
			continue
		}
		if strings.Contains(f.Type, "32") || f.Type == "rune" {
			f.Size = uint8(4)
			continue
		}
		if strings.Contains(f.Type, "64") ||
			strings.HasPrefix(f.Type, "*") ||
			f.Type == "uintptr" ||
			strings.Contains(f.Type, "chan") ||
			strings.HasPrefix(f.Type, "map") ||
			strings.HasPrefix(f.Type, "func") ||
			strings.HasPrefix(f.Type, "func") {
			f.Size = uint8(8)
			continue
		}
		if strings.Contains(f.Type, "128") || f.Type == "string" {
			f.Size = uint8(16)
			continue
		}
		if strings.HasPrefix(f.Type, "[]") {
			f.Size = uint8(24)
			continue
		}
		if f.Type == "int" || f.Type == "uint" {
			if arch == global.I386 {
				f.Size = uint8(4)
				continue
			}
			if arch == global.Amd64 {
				f.Size = uint8(8)
				continue
			}
		}
	}
}

func calcPadding(s *base.Struct, arch global.Architecture) {
	chunk := arch.GetChunkSize()

	for i, f := range s.Fields {
		lastBits := f.Size % chunk
		if lastBits == 0 {
			lastBits = chunk
		}
		if i == 0 {
			f.Index = 0
		}
		if i == len(s.Fields)-1 {
			f.Padding = chunk - lastBits - f.Index
			continue
		}
		next := s.Fields[i+1]
		if f.Index+f.Size+next.Size > chunk {
			f.Padding = chunk - lastBits - f.Index
			next.Index = 0
		} else {
			f.Padding = 0
			next.Index = f.Index + f.Size
		}
	}
}
