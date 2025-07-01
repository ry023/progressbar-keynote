package main

type slideData struct {
	enabledIndexes []int
	chapters       []chapter
}

type chapter struct {
	title      string
	startIndex int
	endIndex   int
}

func (s slideData) pageNum(idx int) int {
	for i, eidx := range s.enabledIndexes {
		if eidx == idx {
			return i + 1
		}
	}
	panic("Index not found in enabledIndexes")
}

func (s slideData) totalPages() int {
	return len(s.enabledIndexes)
}
