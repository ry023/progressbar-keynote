package main

import (
	"fmt"
	"os"
)

func main() {

	s, err := parse()
	if err != nil {
		panic(err)
	}

	if err := deleteProgressBars(); err != nil {
		panic(err)
	}

	if err := os.RemoveAll("tmp"); err != nil {
		panic(fmt.Sprintf("Failed to remove tmp directory: %v", err))
	}

	for _, index := range s.enabledIndexes {
		if err := insertProgressBar(index, s); err != nil {
			panic(fmt.Sprintf("Failed to insert progress bar for slide %d: %v", index, err))
		}
	}
}
