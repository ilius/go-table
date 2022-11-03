package table

import "strings"

const alignSep = " "

type Alignment = func(str string, width int) string

func AlignmentLeft(str string, width int) string {
	strWidth := visualWidth(str)
	if strWidth >= width {
		return str
	}
	n := width - strWidth
	return str + strings.Repeat(alignSep, n)
}

func AlignmentRight(str string, width int) string {
	strWidth := visualWidth(str)
	if strWidth >= width {
		return str
	}
	n := width - strWidth
	return strings.Repeat(alignSep, n) + str
}

func AlignmentCenter(str string, width int) string {
	strWidth := visualWidth(str)
	if strWidth >= width {
		return str
	}
	n := width - strWidth
	left := (n-1)/2 + 1
	right := n - left
	return strings.Repeat(alignSep, left) + str + strings.Repeat(" ", right)
}
