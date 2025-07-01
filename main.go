package main

import (
	"fmt"
)

func main() {
	s := slideData{
		enabledIndexes: []int{0, 1, 2, 3, 30, 100},
		chapters: []chapter{
			{title: "Chapter 1", startIndex: 0, endIndex: 1},
			{title: "Chapter 2", startIndex: 2, endIndex: 3},
			{title: "Chapter 3", startIndex: 30, endIndex: 100},
		},
	}

	svg := generateSVG(1, s)
	fmt.Println(svg)
}
