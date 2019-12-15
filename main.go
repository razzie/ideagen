package main

import (
	"fmt"
)

func main() {
	gen := NewGenerator()
	fmt.Println(gen.Generate())
}
