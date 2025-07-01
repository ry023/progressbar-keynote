package main

import (
	"fmt"
)

const margin = 2
const inactiveBarHeight = 6
const activeBarHeight = 8

func generateSVG(currentIndex int, s slideData) string {
	widthPerPage := 1000.0 / float64(s.totalPages())
	currentPage := s.pageNum(currentIndex)

	svg := `<svg width="1000" height="100" viewBox="0 0 1000 50" xmlns="http://www.w3.org/2000/svg">`

	for _, chapter := range s.chapters {
		startPage := s.pageNum(chapter.startIndex)
		endPage := s.pageNum(chapter.endIndex)

		barHeight := inactiveBarHeight
		if startPage <= currentPage && currentPage <= endPage {
			barHeight = activeBarHeight
		}

		// Draw background bar
		svg += fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" fill="#ccc" />`,
			int(float64(startPage-1)*widthPerPage),
			23-barHeight/2,
			int(float64(endPage-startPage+1)*widthPerPage)-margin,
			barHeight,
		)

		// Draw active bar
		if startPage <= currentPage && currentPage <= endPage {
			// 再生中
			svg += fmt.Sprintf(
				`<rect x="%d" y="%d" width="%d" height="%d" fill="#EF426D" />`,
				int(float64(startPage-1)*widthPerPage),
				23-barHeight/2,
				int(float64(currentPage-startPage+1)*widthPerPage)-margin,
				barHeight,
			)
		} else if endPage < currentPage {
			// 再生済み
			svg += fmt.Sprintf(
				`<rect x="%d" y="20" width="%d" height="%d" fill="#EF426D" />`,
				int(float64(startPage-1)*widthPerPage),
				int(float64(endPage-startPage+1)*widthPerPage)-margin,
				barHeight,
			)
		}

		// Draw chapter title
		// <text x="10" y="10" font-family="筑紫B丸ゴシック" font-size="10" fill="#ccc">チャプター1</text>
		svg += fmt.Sprintf(
			`<text x="%d" y="10" font-family="筑紫B丸ゴシック" font-size="10" fill="#000">%s</text>`,
			int(float64(startPage-1)*widthPerPage)+10,
			chapter.title,
		)
	}
	// Draw current page indicator
	svg += fmt.Sprintf(
		`<circle cx="%d" cy="24" r="8" fill="#EF426D" />`,
		int(float64(currentPage) * widthPerPage),
	)
	svg += `</svg>`

	return svg
}
