package golang

import (
	"encoding/json"
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/global"
	"reflect"
)

type Golang struct{}

func (g *Golang) CalculateSizeof(data string, arch global.Architecture) (string, error) {
	s, err := parseStructs(data)
	if err != nil {
		return "", err
	}
	for _, v := range s {
		calcFieldSizes(v, arch)
		calcPadding(v, arch)
		v.CalcOptimizable(arch)
	}

	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (g *Golang) OptimizeMemoryAlignment(s *base.Struct, arch global.Architecture) (string, error) {
	os := &base.Struct{
		Name:   s.Name,
		Fields: make([]*base.Field, 0),
	}
	chunk := arch.GetChunkSize()
	grp1 := make([]*base.Field, 0)
	grp2 := make([]*base.Field, 0)

	for _, f := range s.Fields {
		if f.Size%chunk == 0 && f.Size >= chunk {
			os.Fields = append(os.Fields, f)
			continue
		}

		f.Index = 0
		f.Padding = chunk - f.Size%chunk
		if f.Size == 0 {
			f.Padding = 0
		}

		if f.Size < chunk {
			grp2 = append(grp2, f)
			continue
		}
		grp1 = append(grp1, f)
	}

	//sort.Sort(base.ByPadding(grp1))
	//sort.Sort(base.ByPadding(grp2))
	sortByLastBits(grp1, chunk)
	sortByLastBits(grp2, chunk)

	for {
		ss, ok := make([]*base.Field, 0), false
		grp2, ss, ok = subsetSum(grp2, int(chunk), int(chunk))
		if !ok {
			break
		}
		for _, f := range ss {
			os.Fields = append(os.Fields, f)
		}
	}

	for i := range grp1 {
		ss, ok := make([]*base.Field, 0), false
		grp2, ss, ok = subsetSum(grp2, int(grp1[i].Padding), int(chunk))
		if !ok {
			os.Fields = append(os.Fields, grp1[i])
			continue
		}
		os.Fields = append(os.Fields, grp1[i])
		for _, f := range ss {
			os.Fields = append(os.Fields, f)
		}
	}

	for _, f := range grp2 {
		os.Fields = append(os.Fields, f)
	}

	calcPadding(os, arch)
	os.CalcOptimizable(arch)
	os.BuildText()

	b, err := json.Marshal(os)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func subsetSum(set []*base.Field, sum, chunk int) ([]*base.Field, []*base.Field, bool) {
	lset := len(set)
	mat := make([][]bool, lset+1)
	// Fill base 0 column
	for i := 0; i <= lset; i++ {
		mat[i] = make([]bool, sum+1)
		mat[i][0] = true
	}

	// From row 1 to lset-1
	for i := 1; i <= lset; i++ {
		for j := 0; j <= sum; j++ { // Normal rows (1..lset-1)
			if int(set[i-1].Size)%chunk > j {
				mat[i][j] = mat[i-1][j]
				continue
			}
			mat[i][j] = mat[i-1][j] || mat[i-1][j-int(int(set[i-1].Size)%chunk)]
		}
	}

	// Pretty print matrix
	//fmt.Printf("\t")
	//for i := 0; i <= sum; i++ {
	//	fmt.Printf("%v\t\t", i)
	//}
	//fmt.Println()
	//for i := 0; i <= lset; i++ {
	//	if i == 0 {
	//		fmt.Printf("0\t")
	//	} else {
	//		fmt.Printf("%v\t", int(set[i-1].Size)%chunk)
	//	}
	//	for j := 0; j <= sum; j++ {
	//		fmt.Printf("%t\t", mat[i][j])
	//	}
	//	fmt.Println()
	//}

	// If the final result is true so reverse the flow to get the subset
	if mat[lset][sum] {
		i, j := lset, sum
		ssIdx := make([]int, 0)
		for {
			if i == 0 || j == 0 || !mat[i][j] {
				break
			}
			if mat[i][j] {
				if !mat[i-1][j] {
					ssIdx = append(ssIdx, i-1)
					j = j - int(int(set[i-1].Size)%chunk)
					i -= 1
					continue
				}
				i -= 1
			}
		}

		subset := make([]*base.Field, 0)
		for _, v := range ssIdx {
			subset = append(subset, set[v])
		}
		for _, ss := range subset {
			for i, s := range set {
				if reflect.DeepEqual(s, ss) {
					set = append(set[:i], set[i+1:]...)
					break
				}
			}
		}
		return set, subset, true
	}

	return set, nil, false
}

func sortByLastBits(fields []*base.Field, chunk uint) {
	lf := len(fields)
	for i := 0; i < lf -2; i++ {
		for j := lf -1; j > i; j-- {
			if fields[i].Size%chunk > fields[j].Size/chunk {
				fields[i], fields[j] = fields[j], fields[i]
			}
		}
	}
}

func calcFieldSizes(s *base.Struct, arch global.Architecture) {
	for _, f := range s.Fields {
		f.Size = getTypeSize(f.Type, arch)
	}
}

func calcPadding(s *base.Struct, arch global.Architecture) {
	chunk := arch.GetChunkSize()

	for i, f := range s.Fields {
		lastBits := f.Size % uint(chunk)
		if lastBits == 0 {
			lastBits = uint(chunk)
		}
		if i == 0 {
			f.Index = 0
		}
		if i == len(s.Fields)-1 {
			if f.Size == 0 {
				f.Padding = chunk - f.Index
				continue
			}
			f.Padding = chunk - lastBits - f.Index
			continue
		}
		next := s.Fields[i+1]
		if f.Index+f.Size+next.Size > chunk {
			if f.Size > chunk {
				if lastBits+next.Size <= chunk {
					f.Padding = 0
					next.Index = lastBits
					continue
				}
				f.Padding = chunk - lastBits - f.Index
				next.Index = 0
				continue
			}
			f.Padding = chunk - lastBits - f.Index
			next.Index = 0
			continue
		}
		f.Padding = 0
		next.Index = f.Index + f.Size
	}
}
