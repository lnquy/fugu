package languages

import (
	"fmt"
	"github.com/lnquy/fugu/modules/global"
	log "github.com/sirupsen/logrus"

	cpl "github.com/lnquy/fugu/languages/c_cpp"
	gopl "github.com/lnquy/fugu/languages/golang"
	javapl "github.com/lnquy/fugu/languages/java"
)

type ProgrammingLanguage interface {
	CalculateSizeof(data string, arch global.Architecture) (string, error)
}

var (
	pLangs = make(map[string]ProgrammingLanguage)
)

func init() {
	langs := global.GetProgrammingLanguages()
	for _, lang := range langs {
		register(lang)
	}
}

func CalcSizeOf(data string, lang global.Language, arch global.Architecture) (string, error) {
	log.Infof("plangs: %v - Lang: %s", pLangs, lang.String())
	pl, ok := pLangs[lang.String()]
	if !ok {
		return "", fmt.Errorf("fugu: %s language is not supported yet", lang.String())
	}
	return pl.CalculateSizeof(data, arch)
}

func register(lang string) {
	el := global.LanguageEnum(lang)
	switch el {
	case global.Go:
		pLangs[lang] = &gopl.Golang{}
	case global.C:
		pLangs[lang] = &cpl.C_CPP{}
	case global.Java:
		pLangs[lang] = &javapl.Java{}
	default:
	}
}
