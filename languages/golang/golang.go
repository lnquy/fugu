package golang

import (
	"encoding/json"
	"github.com/lnquy/fugu/languages/base"
	"github.com/lnquy/fugu/modules/global"
	"reflect"
	"sort"
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
	return g.OptimizeMemoryAlignment2(s, arch)

	//os := &base.Struct{
	//	Name: s.Name,
	//}
	//chunk := arch.GetChunkSize()
	//optm := make([]*base.Field, 0)
	//
	//for _, f := range s.Fields {
	//	if f.Size%chunk == 0 && f.Size >= chunk {
	//		os.Fields = append(os.Fields, f)
	//		continue
	//	}
	//
	//	f.Index = 0
	//	f.Padding = chunk - f.Size%chunk
	//	if f.Size == 0 {
	//		f.Padding = 0
	//	}
	//	optm = append(optm, f)
	//}
	//
	//// TODO: Not solve the special case when (size > chunk && size % chunk != 0) yet!
	//sort.Sort(base.BySize(optm))
	//os.Fields = append(os.Fields, optm...)
	//calcPadding(os, arch)
	//os.CalcOptimizable(arch)
	//os.BuildText()
	//
	//b, err := json.Marshal(os)
	//if err != nil {
	//	return "", err
	//}
	//return string(b), nil
}

func (g *Golang) OptimizeMemoryAlignment2(s *base.Struct, arch global.Architecture) (string, error) {
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

	sort.Sort(base.ByPadding(grp1))
	sort.Sort(base.ByPadding(grp2))

	for {
		ss, ok := make([]*base.Field, 0), false
		grp2, ss, ok = subsetSum(grp2, int(chunk))
		if !ok {
			break
		}
		for _, f := range ss {
			os.Fields = append(os.Fields, f)
		}
	}

	for i := range grp1 {
		ss, ok := make([]*base.Field, 0), false
		grp2, ss, ok = subsetSum(grp2, int(chunk - grp1[i].Padding))
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

func subsetSum(set []*base.Field, sum int) ([]*base.Field, []*base.Field, bool) {
	lset := len(set)
	mat := make([][]bool, lset+1)
	// Fill base 0 column
	for i := 0; i <= lset; i++ {
		mat[i] = make([]bool, sum+1)
		mat[i][0] = true
	}

	// Traverse the matrix and fill up values
	//for j := 0; j <= sum; j++ { // Special case for first row
	//	if int(set[0].Padding) >= j {
	//		mat[0][j] = true
	//		continue
	//	}
	//	mat[0][j] = false
	//}

	// From row 1 to lset-1
	for i := 1; i <= lset; i++ {
		//if i == lset { // Special case for last row. TODO: May remove this
		//	mat[i][sum] = mat[i-1][sum] || mat[i-1][sum - int(set[i-1].Padding)]
		//	continue
		//}

		for j := 0; j <= sum; j++ { // Normal rows (1..lset-1)
			if int(set[i-1].Padding) > j {
				mat[i][j] = mat[i-1][j]
				continue
			}
			mat[i][j] = mat[i-1][j] || mat[i-1][j-int(set[i-1].Padding)]
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
	//		fmt.Printf("%v\t", set[i-1].Padding)
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
					j = j - int(set[i-1].Padding)
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
			if f.Size > chunk && (f.Size % chunk + next.Size <= chunk) {
				f.Padding = 0
				next.Index = f.Size % chunk
				continue
			}
			if f.Size < chunk {
				f.Padding = chunk - lastBits - f.Index
				next.Index = 0
			}
		} else {
			f.Padding = 0
			next.Index = f.Index + f.Size
		}
	}
}
