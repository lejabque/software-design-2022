package main

import "github.com/lejabque/software-design-2022/parser/internal/calculator"

func main() {
	input := "(23 + 10) * 5 - 3 * (32 + 5) * (10 - 4 * 5) + 8 / 2"
	calculator := calculator.Calculator{}
	printed, res, err := calculator.Calculate(input)
	if err != nil {
		panic(err)
	}
	println(printed, "=", res)
}
