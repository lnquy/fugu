package golang

import (
	"github.com/lnquy/fugu/modules/global"
	"strings"
	"regexp"
	"strconv"
)

var (
	fixedSizes map[string]uint
	sliceRegex = regexp.MustCompile(`\A\s*\[\][a-zA-Z0-9_]+\s*\z`)
	arrayRegex = regexp.MustCompile(`\A\s*\[[0-9]+\](\[([0-9]+)?\])*[a-zA-Z0-9_]+\s*\z`)
)

func init() {
	fixedSizes = make(map[string]uint)
	fixedSizes["bool"] = 1
	fixedSizes["int8"] = 1
	fixedSizes["uint8"] = 1
	fixedSizes["byte"] = 1

	fixedSizes["int16"] = 2
	fixedSizes["uint16"] = 2
	fixedSizes["uint16"] = 2

	fixedSizes["int32"] = 4
	fixedSizes["uint32"] = 4
	fixedSizes["rune"] = 4

	fixedSizes["float32"] = 4
	fixedSizes["uint32"] = 4
	fixedSizes["rune"] = 4

	fixedSizes["float64"] = 8
	fixedSizes["complex64"] = 8
	fixedSizes["uintptr"] = 8

	fixedSizes["string"] = 16
	fixedSizes["complex128"] = 16
}

func getTypeSize(t string, arch global.Architecture) uint {
	if size, ok := fixedSizes[t]; ok {
		return size
	}
	if t == "int" || t == "uint" {
		switch arch {
		case global.I386:
			return 4
		case global.Amd64:
			return 8
		}
	}

	if t == "struct{}" || strings.HasPrefix(t, "[0]") { // TODO: Special cases for nested struct
		return 0
	}
	if strings.HasPrefix(t, "*") || strings.HasPrefix(t, "map[") { // Pointer & map
		return 8
	}
	if strings.HasPrefix(t, "chan") || strings.HasPrefix(t, "<-chan") { // Channel
		return 8
	}
	if strings.HasPrefix(t, "func") { // Function
		return 8
	}
	if sliceRegex.MatchString(t) { // Slice
		return 24
	}
	if arrayRegex.MatchString(t) { // Array: Sizeof([N]Type) = N * Sizeof(Type)
		lb := strings.Index(t, "[")
		rb := strings.Index(t, "]")
		i, err := strconv.Atoi(t[lb+1:rb])
		if err != nil {
			return 0
		}
		return uint(i) * getTypeSize(t[rb+1:], arch)
	}
	return 0
}

/*
| Type | `Sizeof()` bytes |
| ---- | ---------------: |
| `struct{}` | 0 | *
| `[0]Type` | 0 | *
| `bool` | 1 |
| `int8`, `uint8`, `byte` | 1 |
| `int16`, `uint16` | 2 |
| `int32`, `uint32`, `rune` | 4 |
| `float32` | 4 |
| `int`, `uint` | 4 or 8 | *
| `int64`, `uint64` | 8 |
| `float64` | 8 |
| `complex64` | 8 |
| `uintptr` | 8 |
| `*struct{}`, `*Type` | 8 | *
| `map[Type1]Type2` | 8 | *
| `chan Type` | 8 | *
| `func()` | 8 | *
| `string` | 16 |
| `complex128` | 16 |
| `[]Type` | 24 | *
*/
