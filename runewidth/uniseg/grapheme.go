package uniseg

import "unicode/utf8"

// Graphemes implements an iterator over Unicode grapheme clusters, or
// user-perceived characters. While iterating, it also provides information
// about word boundaries, sentence boundaries, line breaks, and monospace
// character widths.
//
// After constructing the class via [NewGraphemes] for a given string "str",
// [Graphemes.Next] is called for every grapheme cluster in a loop until it
// returns false. Inside the loop, information about the grapheme cluster as
// well as boundary information and character width is available via the various
// methods (see examples below).
//
// Using this class to iterate over a string is convenient but it is much slower
// than using this package's [Step] or [StepString] functions or any of the
// other specialized functions starting with "First".
type Graphemes struct {
	// The original string.
	original string

	// The remaining string to be parsed.
	remaining string

	// The current grapheme cluster.
	cluster string

	// The byte offset of the current grapheme cluster relative to the original
	// string.
	offset int

	// The current boundary information of the [Step] parser.
	boundaries int

	// The current state of the [Step] parser.
	state int
}

// NewGraphemes returns a new grapheme cluster iterator.
func NewGraphemes(str string) *Graphemes {
	return &Graphemes{
		original:  str,
		remaining: str,
		state:     -1,
	}
}

// Next advances the iterator by one grapheme cluster and returns false if no
// clusters are left. This function must be called before the first cluster is
// accessed.
func (g *Graphemes) Next() bool {
	if len(g.remaining) == 0 {
		// We're already past the end.
		g.state = -2
		g.cluster = ""
		return false
	}
	g.offset += len(g.cluster)
	g.cluster, g.remaining, g.boundaries, g.state = StepString(g.remaining, g.state)
	return true
}

// Runes returns a slice of runes (code points) which corresponds to the current
// grapheme cluster. If the iterator is already past the end or [Graphemes.Next]
// has not yet been called, nil is returned.
func (g *Graphemes) Runes() []rune {
	if g.state < 0 {
		return nil
	}
	return []rune(g.cluster)
}

// Str returns a substring of the original string which corresponds to the
// current grapheme cluster. If the iterator is already past the end or
// [Graphemes.Next] has not yet been called, an empty string is returned.
func (g *Graphemes) Str() string {
	return g.cluster
}

// Positions returns the interval of the current grapheme cluster as byte
// positions into the original string. The first returned value "from" indexes
// the first byte and the second returned value "to" indexes the first byte that
// is not included anymore, i.e. str[from:to] is the current grapheme cluster of
// the original string "str". If [Graphemes.Next] has not yet been called, both
// values are 0. If the iterator is already past the end, both values are 1.
func (g *Graphemes) Positions() (int, int) {
	if g.state == -1 {
		return 0, 0
	} else if g.state == -2 {
		return 1, 1
	}
	return g.offset, g.offset + len(g.cluster)
}

// Width returns the monospace width of the current grapheme cluster.
func (g *Graphemes) Width() int {
	if g.state < 0 {
		return 0
	}
	return g.boundaries >> ShiftWidth
}

// Reset puts the iterator into its initial state such that the next call to
// [Graphemes.Next] sets it to the first grapheme cluster again.
func (g *Graphemes) Reset() {
	g.state = -1
	g.offset = 0
	g.cluster = ""
	g.remaining = g.original
}

// The number of bits the grapheme property must be shifted to make place for
// grapheme states.
const shiftGraphemePropState = 4

// FirstGraphemeClusterInString is like [FirstGraphemeCluster] but its input and
// outputs are strings.
func FirstGraphemeClusterInString(str string, state int) (cluster, rest string, width, newState int) {
	// An empty string returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := utf8.DecodeRuneInString(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		var prop int
		if state < 0 {
			prop = property(graphemeCodePoints, r)
		} else {
			prop = state >> shiftGraphemePropState
		}
		return str, "", runeWidth(r, prop), grAny | (prop << shiftGraphemePropState)
	}

	// If we don't know the state, determine it now.
	var firstProp int
	if state < 0 {
		state, firstProp, _ = transitionGraphemeState(state, r)
	} else {
		firstProp = state >> shiftGraphemePropState
	}
	width += runeWidth(r, firstProp)

	// Transition until we find a boundary.
	for {
		var (
			prop     int
			boundary bool
		)

		r, l := utf8.DecodeRuneInString(str[length:])
		state, prop, boundary = transitionGraphemeState(state&maskGraphemeState, r)

		if boundary {
			return str[:length], str[length:], width, state | (prop << shiftGraphemePropState)
		}

		if r == vs16 {
			width = 2
		} else if firstProp != prExtendedPictographic && firstProp != prRegionalIndicator && firstProp != prL {
			width += runeWidth(r, prop)
		} else if firstProp == prExtendedPictographic {
			if r == vs15 {
				width = 1
			} else {
				width = 2
			}
		}

		length += l
		if len(str) <= length {
			return str, "", width, grAny | (prop << shiftGraphemePropState)
		}
	}
}
