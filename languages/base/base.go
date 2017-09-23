package base

import (
	"github.com/lnquy/fugu/modules/global"
	"fmt"
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
		Size    uint8  `json:"size"`
		Index   uint8  `json:"index"`
		Padding uint8  `json:"padding"`
	}

	Info struct {
		Text         string `json:"text"`
		TotalPadding uint16 `json:"total_padding"`
		TotalSize    uint16 `json:"total_size"`
		Optimizable  bool   `json:"optimizable"`
	}
)

func (s *Struct) CalcInfoTotals() {
	s.TotalPadding = 0
	for _, v := range s.Fields {
		s.TotalPadding += uint16(v.Padding)
	}
	s.TotalSize = s.TotalPadding
	for _, v := range s.Fields {
		s.TotalSize += uint16(v.Size)
	}
}

func (s *Struct) CalcOptimizable(arch global.Architecture) {
	s.CalcInfoTotals()

	if s.TotalPadding >= uint16(arch.GetChunkSize()) {
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
