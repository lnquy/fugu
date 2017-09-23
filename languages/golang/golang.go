package golang

import (
	"encoding/json"
	"fmt"
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/global"
	"github.com/lnquy/fugu/modules/util"
	"regexp"
	"strings"
	"time"
	"sort"
)

type Golang struct{}

var (
	sNameRegex  = regexp.MustCompile(`\A\s*type\s+([a-zA-Z0-9_]*\s+)?struct\s*{\s*\z`)
	sFieldRegex *regexp.Regexp
	sEndRegex   = regexp.MustCompile(`\A\s*}\s*\z`)
)

func init() {
	reg := fmt.Sprintf(`\A\s*([a-zA-Z_][a-zA-Z0-9_]*)\s+([a-zA-Z0-9_\[\]\{\}\<\-\>*]+(\s+[a-zA-Z0-9_\[\]\{\}\<\-\>*]+)?)\s*(%s.*)?\s*\z`, "`")
	sFieldRegex = regexp.MustCompile(reg)
}

func (g *Golang) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	s := parseStructs(data)
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
		if f.Size % chunk == 0 && f.Size >= chunk {
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

func parseStructs(data string) []*base.Struct {
	retStructs := make([]*base.Struct, 0)
	s := &base.Struct{}

	lines := strings.Split(data, util.LineBreak())
	for _, line := range lines {
		if sNameRegex.MatchString(line) {
			s.Text = line + "\n"
			s.Name = strings.TrimSpace(sNameRegex.FindStringSubmatch(line)[1])
			if s.Name == "" {
				s.Name = fmt.Sprintf("%v", time.Now().Unix())
			}
			s.Fields = make([]*base.Field, 0)
			continue
		}
		if sEndRegex.MatchString(line) {
			s.Text += "}"
			retStructs = append(retStructs, s)
			s = &base.Struct{}
			continue
		}
		if sFieldRegex.MatchString(line) {
			s.Text += line + "\n"
			m := sFieldRegex.FindStringSubmatch(line)
			s.Fields = append(s.Fields, &base.Field{
				Name: m[1],
				Type: m[2],
			})
		}
	}
	return retStructs
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
			f.Padding = chunk - lastBits -f.Index
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
