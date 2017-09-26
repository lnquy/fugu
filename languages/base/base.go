package base

import (
	"fmt"
	"github.com/lnquy/fugu/modules/global"
	"github.com/prometheus/common/log"
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
		Size    uint   `json:"size"`
		Index   uint   `json:"index"`
		Padding uint   `json:"padding"`
	}

	Info struct {
		Text         string `json:"text"`
		TotalPadding uint   `json:"total_padding"`
		TotalSize    uint   `json:"total_size"`
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
	chunk := arch.GetChunkSize()
	if s.TotalPadding < chunk {
		s.Optimizable = false
		return
	}
	s.Optimizable = s.checkPossibleAdd(chunk)
}

func (s *Struct) checkPossibleAdd(chunk uint) bool {
	lf := len(s.Fields)
	for i := 0; i < lf-2; i++ {
		for j := lf - 1; j > i; j-- {
			if (s.Fields[i].Padding != 0 && s.Fields[j].Padding != 0) &&
				(s.Fields[i].Padding+s.Fields[j].Padding <= chunk) &&
				((s.Fields[i].Size < chunk && s.Fields[j].Size < chunk) ||
					(s.Fields[i].Size > chunk && s.Fields[j].Size < chunk) ||
					(s.Fields[i].Size < chunk && s.Fields[j].Size > chunk)) {
				log.Info(s.Fields[i], s.Fields[j])
				return true
			}
		}
	}
	return false
}

func (s *Struct) BuildText() {
	s.Text = fmt.Sprintf("type %s struct {\n", s.Name)
	for _, f := range s.Fields {
		s.Text += fmt.Sprintf("\t%s %s\n", f.Name, f.Type)
	}
	s.Text += "}"
}

type ByPadding []*Field

func (a ByPadding) Len() int           { return len(a) }
func (a ByPadding) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPadding) Less(i, j int) bool { return a[i].Padding < a[j].Padding }
