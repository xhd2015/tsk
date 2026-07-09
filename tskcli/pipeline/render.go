package pipeline

import "strings"

// Render returns the compact pipeline diagram. When plain is true, ASCII boxes
// are used instead of Unicode box-drawing characters.
func Render(plain bool) string {
	art := UnicodeArt
	if plain {
		art = ASCIIArt
	}
	if art == "" {
		return "\n"
	}
	return strings.TrimRight(art, "\n") + "\n"
}