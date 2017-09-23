package golang

import (
	"github.com/lnquy/fugu/modules/global"
	"github.com/lnquy/fugu/languages/base"
)

type Java struct {}

func (j *Java) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	return "Java", nil
}

func (j *Java) OptimizeMemoryAlignment(s *base.Struct, arch global.Architecture) (string, error) {
	return "Java", nil
}
