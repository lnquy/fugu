package base

import (
	"fmt"
	"github.com/lnquy/fugu/modules/global"
)

type (
	Struct struct {
		Name   string   `json:"name"`
		Fields []*Field `json:"fields"`
		Info   `json:"info"`
	}

	Field struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Size    uint `json:"size"`
		Index   uint  `json:"index"`
		Padding uint  `json:"padding"`
	}

	Info struct {
		Text         string `json:"text"`
		TotalPadding uint `json:"total_padding"`
		TotalSize    uint `json:"total_size"`
		Optimizable  bool   `json:"optimizable"`
	}
)

func (s *Struct) CalcInfoTotals() {
	s.TotalPadding = 0
	for _, v := range s.Fields {
		s.TotalPadding += v.Padding
	}
	s.TotalSize = s.TotalPadding
	for _, v := range s.Fields {
		s.TotalSize += v.Size
	}
}

func (s *Struct) CalcOptimizable(arch global.Architecture) {
	s.CalcInfoTotals()

	if s.TotalPadding >= arch.GetChunkSize() {
		s.Optimizable = true
		return
	}
	s.Optimizable = false
}

func (s *Struct) BuildText() {
	s.Text = fmt.Sprintf("type %s struct {\n", s.Name)
	for _, f := range s.Fields {
		s.Text += fmt.Sprintf("\t%s %s\n", f.Name, f.Type)
	}
	s.Text += "}"
}

type BySize []*Field

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size }
