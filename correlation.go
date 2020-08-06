/*
Spearman and pearson correlation in Go
from : https://gist.github.com/barthr/ec5c12033e969b6ca64150f282220b8b
*/

//package stats
package main

import (
	"errors"
	"math"
	"sort"
)

// Correlator represent the contract for a correlation algorithm
// It contains 2 arguments which represent the datasets where the
// correlation has to be calculated on.
// It returns the correlation between the 2 datasets
type Correlator interface {
	Correlate(Float64Data, Float64Data) float64
}

type Float64Data []float64

func (f Float64Data) Len() int { return len(f) }

func (f Float64Data) Get(i int) float64 { return f[i] }

// Pearson calculates the Pearson product-moment correlation coefficient between two variables.
func Pearson(data1, data2 Float64Data) (float64, error) {
	var sum [5]float64

	var n = float64(data1.Len())
	for i := 0; i < data1.Len(); i++ {
		x := data1[i]
		y := data2[i]

		sum[0] += x * y
		sum[1] += x
		sum[2] += y
		sum[3] += math.Pow(x, 2)
		sum[4] += math.Pow(y, 2)
	}

	sqrtX := math.Sqrt(sum[3] - (math.Pow(sum[1], 2) / n))
	sqrtY := math.Sqrt(sum[4] - (math.Pow(sum[2], 2) / n))

	dividend := sum[0] - ((sum[1] * sum[2]) / n)
	divisor := sqrtX * sqrtY

	return dividend / divisor, nil
}

type rank struct {
	X     float64
	Y     float64
	Xrank float64
	Yrank float64
}

func Spearman(data1, data2 Float64Data) (float64, error) {
	if data1.Len() < 3 || data2.Len() != data1.Len() {
		return math.NaN(), errors.New("Invalid size of data")
	}

	ranks := []rank{}

	for index := 0; index < data1.Len(); index++ {
		x := data1.Get(index)
		y := data2.Get(index)
		ranks = append(ranks, rank{
			X: x,
			Y: y,
		})
	}

	sort.Slice(ranks, func(i int, j int) bool {
		return ranks[i].X < ranks[j].X
	})

	for position := 0; position < len(ranks); position++ {
		ranks[position].Xrank = float64(position) + 1

		duplicateValues := []int{position}
		for nested, p := range ranks {
			if ranks[position].X == p.X {
				if position != nested {
					duplicateValues = append(duplicateValues, nested)
				}
			}
		}
		sum := 0
		for _, val := range duplicateValues {
			sum += val
		}

		avg := float64((sum + len(duplicateValues))) / float64(len(duplicateValues))
		ranks[position].Xrank = avg

		for index := 1; index < len(duplicateValues); index++ {
			ranks[duplicateValues[index]].Xrank = avg
		}

		position += len(duplicateValues) - 1
	}

	sort.Slice(ranks, func(i int, j int) bool {
		return ranks[i].Y < ranks[j].Y
	})

	for position := 0; position < len(ranks); position++ {
		ranks[position].Yrank = float64(position) + 1

		duplicateValues := []int{position}
		for nested, p := range ranks {
			if ranks[position].Y == p.Y {
				if position != nested {
					duplicateValues = append(duplicateValues, nested)
				}
			}
		}
		sum := 0
		for _, val := range duplicateValues {
			sum += val
		}
		// fmt.Println(sum + len(duplicateValues))
		avg := float64((sum + len(duplicateValues))) / float64(len(duplicateValues))
		ranks[position].Yrank = avg

		for index := 1; index < len(duplicateValues); index++ {
			ranks[duplicateValues[index]].Yrank = avg
		}

		position += len(duplicateValues) - 1
	}

	xRanked := []float64{}
	yRanked := []float64{}

	for _, rank := range ranks {
		xRanked = append(xRanked, rank.Xrank)
		yRanked = append(yRanked, rank.Yrank)
	}

	return Pearson(xRanked, yRanked)
}
