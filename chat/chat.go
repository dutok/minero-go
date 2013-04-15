package chat

import (
	"regexp"
)

const (
	Black     = "§0" // #000
	DarkBlue  = "§1" // #00a
	DarkGreen = "§2" // #0a0
	DarkCyan  = "§3" // #0aa
	DarkRed   = "§4" // #a00
	Purple    = "§5" // #a0a
	Gold      = "§6" // #fa0
	Gray      = "§7" // #aaa
	DarkGray  = "§8" // #555
	Blue      = "§9" // #55f
	Green     = "§a" // #5f5
	Cyan      = "§b" // #5ff
	Red       = "§c" // #f55
	Pink      = "§d" // #f5f
	Yellow    = "§e" // #ff5
	White     = "§f" // #fff

	Random     = "§k"
	Bold       = "§l"
	Strike     = "§m"
	Underlined = "§n"
	Italic     = "§o"
	Rest       = "§r"
)

const (
	reColor  = "§([0-9a-f])"
	reFormat = "§([klmnor])"
	reBoth   = "§([0-9a-fklmnor])"
)

// IsColor checks if s contains a color code. If len([]rune(s)) != 2 returns
// false.
func IsColor(s string) bool {
	r := []rune(s)
	if len(r) != 2 {
		return false
	}
	a := r[1] >= 48 && r[1] <= 57
	b := r[1] >= 97 && r[1] <= 102
	return a || b
}

// IsColor checks if s contains a format code. If len([]rune(s)) != 2 returns
// false.
func IsFormat(s string) bool {
	r := []rune(s)
	if len(r) != 2 {
		return false
	}
	switch r[1] {
	case 107, 108, 109, 110, 111, 114:
		return true
	}
	return false
}

// Translate translates strings using an alternate control sequence char to a
// string that uses '§' as control sequence char.
func Translate(s, pat string) string {
	r := regexp.MustCompile(pat + "([0-9a-fklmnor])")
	return r.ReplaceAllString(s, "§$1")
}

// StripColor deletes color control sequences.
func StripColor(s string) string {
	r := regexp.MustCompile(reColor)
	return r.ReplaceAllString(s, "")
}

// StripFormat deletes format control sequences.
func StripFormat(s string) string {
	r := regexp.MustCompile(reFormat)
	return r.ReplaceAllString(s, "")
}

// Strip deletes any chat control sequence.
func Strip(s string) string {
	r := regexp.MustCompile(reBoth)
	return r.ReplaceAllString(s, "")
}
