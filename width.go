package table

import (
	"regexp"

	"github.com/ilius/go-lru"
	"github.com/mattn/go-runewidth"
)

// ilius/go-lru		is a fork of dboslee/lru, without loads of depedencies
// dboslee/lru		is almost as fast as an unbounded map[string]int
// dboslee/lru 		is ~%10 faster than	bluele/gcache
// bluele/gcache	is ~%10 faster than	karlseguin/ccache/v3

var ansiEscapeRE = regexp.MustCompile(
	"[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]",
)

var widthCache = lru.New[string, int](lru.WithCapacity(10000))

// runewidth.FillLeft(str, width) or runewidth.FillRight(str, width) do not work

func visualWidth(str string) int {
	// method 1: without considering CJK, emoji, etc:
	//		len(str) - countMatchesLength(ansiEscapeRE, str)
	// method 2: without considering ANSI colors / escape sequences
	//		return runewidth.StringWidth(str)
	// method 3: considering all above, but a bit slow?
	// 		return runewidth.StringWidth(ansiEscapeRE.ReplaceAllString(str, ""))
	w, _ := widthCache.Get(str)
	if w > 0 {
		return w
	}
	w = runewidth.StringWidth(str) - countMatchesLength(ansiEscapeRE, str)
	widthCache.Set(str, w)
	return w
}

func countMatchesLength(re *regexp.Regexp, s string) int {
	// it does not count the first (invisible) character
	// because runewidth library already ignores that
	total := 0
	for start := 0; start < len(s); {
		remaining := s[start:] // slicing the string is cheap
		loc := re.FindStringIndex(remaining)
		if loc == nil {
			break
		}
		// loc[0] is the start index of the match,
		// loc[1] is the end index (exclusive)
		start += loc[1]
		total += loc[1] - loc[0] - 1
	}
	return total
}
