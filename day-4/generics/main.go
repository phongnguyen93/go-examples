package main

import "fmt"

type Number interface {
	int64 | float64
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var sumNum V
	for _, v := range m {
		sumNum += v
	}
	return sumNum
}

func main() {
	fmt.Println("Generics, sum float: %n", SumNumbers(map[string]float64{"a": 1.3, "b": 2.4}))
	fmt.Println("Generics, sum int: %n", SumNumbers(map[string]int64{"a": 1, "b": 2}))
}
