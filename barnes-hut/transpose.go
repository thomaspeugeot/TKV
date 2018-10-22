package barneshut

import "fmt"

type row []float64
type matrix []row

func printMatrix(m matrix) {
	for _, s := range m {
		fmt.Println(s)
	}
}

func transpose(m matrix) matrix {
	r := make(matrix, len(m[0]))
	for x := range r {
		r[x] = make(row, len(m))
	}
	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}

func transposeFloat64(m [][]float64) [][]float64 {
	r := make([][]float64, len(m[0]))
	for x := range r {
		r[x] = make(row, len(m))
	}
	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}
