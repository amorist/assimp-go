package main

import (
	"flag"

	"github.com/amorist/assimp-go/assimp"
)

func main() {
	a := assimp.NewAssimp()
	var filename string
	flag.StringVar(&filename, "f", "", "模型路径")
	flag.Parse()
	a.Export(filename)
}
