package golang

import (
	"github.com/lnquy/fugu/modules/global"
	"strings"
	"github.com/lnquy/fugu/modules/util"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type Golang struct {}

func (g *Golang) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	parseData(data)
	log.Infof("Map: %v", s)
	return "Go", nil
}

var (
	s = make(map[string][]string)
)

func parseData(data string) {
	lines := strings.Split(data, util.LineBreak())
	var curStruct string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.Trim(line, "\t")
		if strings.HasPrefix(line, "type ") { // struct type
			name := line[5:strings.Index(line, "struct")]
			if name == "" {
				name = xid.New().String() // TODO
			}
			curStruct = name
			s[curStruct] = make([]string, 0)
			continue
		}
		if line == "}" {
			continue
		}
		if curStruct != "" {
			s[curStruct] = append(s[curStruct], line)
		}
	}
}
