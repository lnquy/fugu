package golang

import (
	"encoding/json"
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/global"
	"sort"
)

type Golang struct{}

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

		f.Index = 0
		f.Padding = chunk - f.Size%chunk
		if f.Size == 0 {
			f.Padding = 0
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

func calcFieldSizes(s *base.Struct, arch global.Architecture) {
	for _, f := range s.Fields {
		f.Size = getTypeSize(f.Type, arch)
	}
}

func calcPadding(s *base.Struct, arch global.Architecture) {
	chunk := arch.GetChunkSize()

	for i, f := range s.Fields {
		lastBits := f.Size % uint(chunk)
		if lastBits == 0 {
			lastBits = uint(chunk)
		}
		if i == 0 {
			f.Index = 0
		}
		if i == len(s.Fields)-1 {
			if f.Size == 0 {
				f.Padding = chunk - f.Index
				continue
			}
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
