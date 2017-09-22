package golang

import (
	"fmt"
	"github.com/lnquy/fugu/modules/global"
	"github.com/lnquy/fugu/modules/util"
	"regexp"
	"strings"
	"time"
)

type Abc_Def struct {
}

type (
	Golang struct{}

	Field struct {
		Name    string
		Type    string
		Size    uint8
		Index   uint8
		Padding uint8
	}
)

func (g *Golang) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	s := parseData(data)
	for _, v := range s {
		getFieldSizes(v, arch)
		calcPadding(v, arch)
	}

	for k, v := range s {
		fmt.Printf("%v\n", k)
		for _, f := range v {
			fmt.Printf("    %v\n", f)
		}
	}
	return "Go", nil
}

var (
	sNameRegex  = regexp.MustCompile(`\A\s*type\s+([a-zA-Z0-9_]*\s+)?struct\s*{\s*\z`)
	sFieldRegex = regexp.MustCompile(`\A\s*([a-zA-Z_][a-zA-Z0-9_]*)\s+([a-zA-Z0-9_\[\]\{\}\<\-\>*]+(\s+[a-zA-Z0-9_\[\]\{\}\<\-\>*]+)?)\s*\z`)
	sEndRegex   = regexp.MustCompile(`\A\s*}\s*\z`)
)

func parseData(data string) map[string][]*Field {
	s := make(map[string][]*Field)
	curStruct := ""

	lines := strings.Split(data, util.LineBreak())
	for _, line := range lines {
		if sNameRegex.MatchString(line) {
			curStruct = sNameRegex.FindStringSubmatch(line)[1]
			if curStruct == "" {
				curStruct = fmt.Sprintf("%v", time.Now().Unix())
			}
			s[curStruct] = make([]*Field, 0)
			continue
		}
		if sEndRegex.MatchString(line) {
			curStruct = ""
			continue
		}
		if sFieldRegex.MatchString(line) {
			m := sFieldRegex.FindStringSubmatch(line)
			s[curStruct] = append(s[curStruct], &Field{
				Name: m[1],
				Type: m[2],
			})
		}
	}
	return s
}

func getFieldSizes(fields []*Field, arch global.Architecture) {
	for _, f := range fields {
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

func calcPadding(fields []*Field, arch global.Architecture) {
	chunk := uint8(4)
	if arch == global.Amd64 {
		chunk = uint8(8)
	}

	for i, f := range fields {
		lastBits := f.Size % chunk
		if lastBits == 0 {
			lastBits = chunk
		}
		if i == 0 {
			f.Index = 0
		}
		if i == len(fields)-1 {
			f.Padding = chunk - lastBits
			continue
		}
		next := fields[i+1]
		if f.Index+f.Size+next.Size > chunk {
			f.Padding = chunk - lastBits - f.Index
			next.Index = 0
		} else {
			f.Padding = 0
			next.Index = f.Index + f.Size
		}
	}
}