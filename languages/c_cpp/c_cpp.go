package golang

import (
	"github.com/lnquy/fugu/modules/global"
)

type C_CPP struct {}

func (g *C_CPP) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	return "C/C++", nil
}
