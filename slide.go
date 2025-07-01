package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

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

func parse() (slideData, error) {
	slideCount, err := getSlideCount()
	if err != nil {
		return slideData{}, err
	}

	s := slideData{
		enabledIndexes: make([]int, 0, slideCount),
		chapters:       make([]chapter, 0),
	}

	for i := 1; i <= slideCount; i++ {
		slideNote, err := getSlideNote(i)
		if err != nil {
			return s, err
		}

		if slideNote.Skip {
			continue
		}

		s.enabledIndexes = append(s.enabledIndexes, i)

		if slideNote.StartChapter != "" {
			chapter := chapter{
				title:      slideNote.StartChapter,
				startIndex: i,
			}
			s.chapters = append(s.chapters, chapter)
		} else if len(s.chapters) > 0 {
			s.chapters[len(s.chapters)-1].endIndex = i
		}
	}

	if len(s.chapters) > 0 && s.chapters[len(s.chapters)-1].endIndex == 0 {
		s.chapters[len(s.chapters)-1].endIndex = slideCount
	}

	return s, nil
}

func getSlideCount() (int, error) {
	script := `
      tell application "Keynote"
          count of slides of front document
      end tell
    `
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(out)))
}

type SlideNote struct {
	StartChapter string
	Skip         bool
}

func getSlideNote(slideIndex int) (*SlideNote, error) {
	// タイトルは問答無用でスキップ
	if slideIndex == 0 {
		return &SlideNote{
			Skip: true,
		}, nil
	}
	script := fmt.Sprintf(`
			tell application "Keynote"
			    presenter notes of slide %d of front document as string
			end tell
		`, slideIndex)

	rawNote, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, err
	}

	note := SlideNote{}
	lines := strings.Split(strings.TrimSpace(string(rawNote)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Chapter:") {
			note.StartChapter = strings.TrimSpace(strings.TrimPrefix(line, "Chapter:"))
		}
		if strings.HasPrefix(line, "Skip: true") {
			note.Skip = true
		}
	}
	return &note, nil
}
