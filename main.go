package main

//go:generate go run ./tools/go-bindata/ -prefix data data/...

import (
	"fmt"
)

func main() {
	gen := NewGenerator()
	fmt.Println(gen.Generate())
}
