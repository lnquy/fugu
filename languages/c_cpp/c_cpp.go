package golang

import (
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/global"
)

type C_CPP struct{}

func (c *C_CPP) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	return "C/C++", nil
}

func (c *C_CPP) OptimizeMemoryAlignment(s *base.Struct, arch global.Architecture) (string, error) {
	return "C/C++", nil
}
