package renderer

import "strings"

const glyphHeight = 5
const digitGap = 2 // space between two digit glyphs
const colonGap = 0 // extra space when either neighbour is a colon (it has built-in padding)

// gapBefore returns the number of spaces to insert before the current character
// given the previous character.
func gapBefore(prev, cur rune) int {
	if prev == ':' || cur == ':' {
		return colonGap
	}
	return digitGap
}

// BigDigits renders a string of digits and colons into a slice of 5 strings,
// one per row. Characters not present in the font are skipped.
func BigDigits(s string) []string {
	var rows [glyphHeight]strings.Builder

	prev := rune(0)
	for _, ch := range s {
		glyph, ok := glyphs[ch]
		if !ok {
			continue
		}
		if prev != 0 {
			gap := gapBefore(prev, ch)
			for r := range rows {
				rows[r].WriteString(strings.Repeat(" ", gap))
			}
		}
		prev = ch
		for r, line := range glyph {
			rows[r].WriteString(line)
		}
	}

	result := make([]string, glyphHeight)
	for r, b := range rows {
		result[r] = b.String()
	}
	return result
}
