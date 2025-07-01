package main

import (
	"fmt"
)

func main() {

	s, err := parse()
	if err != nil {
		panic(err)
	}

	svg := generateSVG(s.enabledIndexes[5], s)
	fmt.Println(svg)
}
