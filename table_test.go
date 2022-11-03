package table

import (
	"strconv"
	"testing"

	"github.com/ilius/is/v2"
)

// Fg wraps an 8-bit foreground color code in the ANSI escape sequence
func Fg(code int) string {
	return "\x1b[38;5;" + strconv.Itoa(code) + "m"
}

// Bg wraps an 8-bit background color code in the ANSI escape sequence
func Bg(code int) string {
	return "\x1b[48;5;" + strconv.Itoa(code) + "m"
}

func Test_visualWidth(t *testing.T) {
	is := is.New(t)
	test := func(width int, str string) {
		actualWidth := visualWidth(str)
		is.AddMsg(
			"str=%#v",
			str,
		).Equal(actualWidth, width)
	}

	test(8, "さの.png")
	test(8, Fg(1)+"さの"+Fg(2)+".png"+reset)
	test(8, "たき.png")
	test(10, "いざわ.png")
	test(10, Fg(15)+"いざわ.png"+reset)
	test(10, "うらべ.png")
	test(10, "かずお.png")
	test(10, "きすぎ.png")
	test(10, "さわだ.png")
	test(10, "じとう.png")
	test(10, "そうだ.png")
	test(10, "つばさ.png")
	test(10, "にった.png")
	test(10, "まさお.png")
	test(10, "みさき.png")
	test(10, "みすぎ.png")
	test(12, "いしざき.png")
	test(12, Fg(3)+"いしざき"+Bg(4)+".png"+reset)
	test(12, "そりまち.png")
	test(12, "たかすぎ.png")
	test(12, "ひゅうが.png")
	test(12, "まつやま.png")
	test(12, "もリさき.png")
	test(14, "わかしまづ.png")
	test(14, "わかばやし.png")
}

func TestAlign(t *testing.T) {
	is := is.New(t)
	test := func(alignment Alignment, width int, str string, out string) {
		actualOut := alignment(str, width)
		is.AddMsg(
			"str=%#v",
			str,
		).Equal(actualOut, out)
	}

	test(AlignmentLeft, 11, "さの.png", "さの.png   ")
	{
		str := Fg(1) + "さの" + Fg(2) + Fg(3) + Fg(4) + ".png" + reset
		test(AlignmentLeft, 11, str, str+"   ")
	}
	test(AlignmentRight, 11, "さの.png", "   さの.png")
	{
		str := Fg(1) + "さの" + Fg(2) + Fg(3) + Fg(4) + ".png" + reset
		test(AlignmentRight, 11, str, "   "+str)
	}
	test(AlignmentCenter, 11, "さの.png", "  さの.png ")
	{
		str := Fg(1) + "さの" + Fg(2) + Fg(3) + Fg(4) + ".png" + reset
		test(AlignmentCenter, 11, str, "  "+str+" ")
	}
}
