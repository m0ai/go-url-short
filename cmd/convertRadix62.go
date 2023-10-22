package main

import (
	"fmt"
	"go-url-short/internal/shorten"
)

func main() {
	input := int64(1541815603606036480)
	output := shorten.ConvertRadix62(input)
	fmt.Printf("input: %d, output: %s\n", input, output)
}
