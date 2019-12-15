package main

//go:generate go run ./tools/go-bindata/ -pkg internal -o internal/bindata.go -prefix data data/...

import (
	"fmt"
)

func main() {
	gen := NewGenerator()
	fmt.Println(gen.Generate())
}
