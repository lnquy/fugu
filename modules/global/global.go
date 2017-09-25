package global

import (
	"strconv"
)

type (
	Language     uint16
	Architecture uint8
)

const (
	EnvironmentPrefix = "fugu"
)

const (
	// Programming languages
	UndefinedLanguage Language = iota
	Go
	C
	Java
)

const (
	// Architectures
	UndefinedArchitecture Architecture = iota
	I386
	Amd64
)

var (
	langStrings = []string{"undefined", "go", "c/c++", "java"}
	archStrings = []string{"undefined", "i386", "amd64"}
)

func (l Language) String() string {
	i := uint16(l)
	switch {
	case i <= uint16(Java):
		return langStrings[i]
	default:
		return strconv.Itoa(int(i))
	}
}

func LanguageEnum(in string) Language {
	for i, v := range langStrings {
		if v == in {
			return Language(i)
		}
	}
	return UndefinedLanguage
}

func (a Architecture) GetChunkSize() uint {
	switch a {
	case I386:
		return 4
	case Amd64:
		return 8
	default:
		return 0
	}
}

func (a Architecture) String() string {
	i := uint8(a)
	switch {
	case i <= uint8(Amd64):
		return archStrings[i]
	default:
		return strconv.Itoa(int(i))
	}
}

func ArchitectureEnum(in string) Architecture {
	for i, v := range archStrings {
		if v == in {
			return Architecture(i)
		}
	}
	return UndefinedArchitecture
}

func GetProgrammingLanguages() []string {
	ret := make([]string, len(langStrings))
	copy(ret, langStrings)
	return ret
}
