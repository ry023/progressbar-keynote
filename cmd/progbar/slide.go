package main

import (
	"fmt"
	"os"
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
          activate
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
	script := fmt.Sprintf(`
			tell application "Keynote"
			    activate
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

func deleteProgressBars() error {
	script := `
tell application "Keynote"
  activate
  set slideCount to count of slides of front document
  repeat with i from 1 to slideCount
    set theSlide to slide i of front document
	  set imageCount to count of images of theSlide
	  repeat with j from imageCount to 1 by -1
      set theImage to image j of theSlide
      set imageName to file name of theImage
      if imageName contains "progress_" then
        delete theImage
      end if
	  end repeat
  end repeat
end tell
	`
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func insertProgressBar(currentIndex int, s slideData) error {
	// create svg file
	os.MkdirAll("tmp", os.ModePerm)
	fname := fmt.Sprintf("progress_%02d.svg", currentIndex)
	svg := generateSVG(currentIndex, s)
	file, err := os.Create(fmt.Sprintf("tmp/%s", fname))
	if err != nil {
		return fmt.Errorf("failed to create SVG file: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(svg)
	if err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}

	script := fmt.Sprintf(`
tell application "Keynote"
    activate
    set theDoc to front document
    set slideWidth to width of theDoc
    set slideHeight to height of theDoc
    set barY to slideHeight - 120
    set theSlide to slide %d of theDoc
    set svgPath to POSIX file "%s" as alias
    tell theSlide
        set newImage to make new image with properties {file:svgPath, position:{10, barY}}
        set width of newImage to slideWidth - 20
    end tell
end tell
`, currentIndex, fmt.Sprintf("%s/tmp/%s", os.Getenv("PWD"), fname))
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to insert progress bar: %w", err)
	}
	return nil
}
