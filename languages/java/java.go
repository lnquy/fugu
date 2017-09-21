package golang

import (
	"github.com/lnquy/fugu/modules/global"
)

type Java struct {}

func (g *Java) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	return "Java", nil
}
