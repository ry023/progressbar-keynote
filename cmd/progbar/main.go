package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	arg := os.Args
	if len(arg) != 2 || (arg[1] != "insert" && arg[1] != "remove") {
		fmt.Println("Usage: progbar [insert|remove]")
		os.Exit(1)
	}

	cmd := arg[1]

	if cmd == "remove" {
		remove()
	} else if cmd == "insert" {
		insert()
	}
}

func remove() {
	log.Print("Deleting existing progress bars...")
	if err := deleteProgressBars(); err != nil {
		panic(err)
	}
	log.Println("Complete.")
}

func insert() {
		log.Print("Parsing slides...")
		s, err := parse()
		if err != nil {
			panic(err)
		}
		log.Println("Complete.")

		log.Print("Deleting existing progress bars...")
		if err := deleteProgressBars(); err != nil {
			panic(err)
		}
		log.Println("Complete.")

		log.Print("Inserting progress bars...")
		if err := os.RemoveAll("tmp"); err != nil {
			panic(fmt.Sprintf("Failed to remove tmp directory: %v", err))
		}
		for _, index := range s.enabledIndexes {
			if err := insertProgressBar(index, s); err != nil {
				panic(fmt.Sprintf("Failed to insert progress bar for slide %d: %v", index, err))
			}
		}
		log.Println("Complete.")
}
