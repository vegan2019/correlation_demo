package main

import "fmt"

func main() {
	var data1 Float64Data = Float64Data{3.2, 2.6, 3.8}
	var data2 Float64Data = Float64Data{6.4, 5.2, 7.6}
	var data3 Float64Data = Float64Data{6.4, 5.2, 7.6 + 2}

	corr12, _ := Pearson(data1, data2)
	corr13, _ := Pearson(data1, data3)

	fmt.Println(corr12)
	fmt.Println(corr13)
}
