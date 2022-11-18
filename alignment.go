package table

import "strings"

const alignSep = " "

type Alignment = func(str string, width uint16) string

func AlignmentLeft(str string, width uint16) string {
	strWidth := visualWidth(str)
	if strWidth >= width {
		return str
	}
	return str + strings.Repeat(alignSep, int(width-strWidth))
}

func AlignmentRight(str string, width uint16) string {
	strWidth := visualWidth(str)
	if strWidth >= width {
		return str
	}
	return strings.Repeat(alignSep, int(width-strWidth)) + str
}

func AlignmentCenter(str string, width uint16) string {
	strWidth := visualWidth(str)
	if strWidth >= width {
		return str
	}
	n := int(width - strWidth)
	left := (n-1)/2 + 1
	right := n - left
	return strings.Repeat(alignSep, left) + str + strings.Repeat(" ", right)
}
