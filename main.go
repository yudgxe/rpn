package main

import (
	"fmt"
	"log"

	lrpn "github.com/yudgxe/lexer/rpn"
)

func main() {
	test := "a > "
	rpn, err := lrpn.NewParser(test).Parse()
	if err != nil {
		log.Panic(err)
	}

	for _, v := range rpn {
		fmt.Print(v.Kind)
		fmt.Print(" ")
		fmt.Print(v.Value)

		fmt.Println()
	}
	values := map[string]int{
		"a": 10,
		"b": 20,
	}
	result, err := lrpn.Execute[int](rpn, values, func(i int) float64 { return float64(i) })
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}
