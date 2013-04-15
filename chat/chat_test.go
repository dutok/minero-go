package chat

import (
	"testing"
)

var stripTests = []struct {
	in, out string
}{
	{"", ""},
	{"§0Black", "Black"},
	{"§1DarkBlue", "DarkBlue"},
	{"§2DarkGreen", "DarkGreen"},
	{"§3DarkCyan", "DarkCyan"},
	{"§4DarkRed", "DarkRed"},
	{"§5Purple", "Purple"},
	{"§6Gold", "Gold"},
	{"§7Gray", "Gray"},
	{"§8DarkGray", "DarkGray"},
	{"§9Blue", "Blue"},
	{"§aGreen", "Green"},
	{"§bCyan", "Cyan"},
	{"§cRed", "Red"},
	{"§dPink", "Pink"},
	{"§eYellow", "Yellow"},
	{"§fWhite", "White"},
	{"§cRed§r and §aGreen", "Red§r and Green"},
}

func TestStrip(t *testing.T) {
	for index, tt := range stripTests {
		out := StripColor(tt.in)
		if out != tt.out {
			t.Fatalf("%d. Failed. Expecting %q not %q.\n", index, tt.out, out)
		}
	}
}

var isColorTests = []struct {
	in  string
	out bool
}{
	{"", false},
	{"§0", true},
	{"§1", true},
	{"§2", true},
	{"§3", true},
	{"§4", true},
	{"§5", true},
	{"§6", true},
	{"§7", true},
	{"§8", true},
	{"§9", true},
	{"§a", true},
	{"§b", true},
	{"§c", true},
	{"§d", true},
	{"§e", true},
	{"§f", true},
	{"§k", false},
	{"§l", false},
	{"§m", false},
	{"§n", false},
	{"§o", false},
	{"§r", false},
}

func TestIsColor(t *testing.T) {
	for index, tt := range isColorTests {
		out := IsColor(tt.in)
		if out != tt.out {
			t.Fatalf("%d. Failed. IsColor(%s) => %v not %v.\n", index, tt.in, tt.out, out)
		}
	}
}

var isFormatTests = []struct {
	in  string
	out bool
}{
	{"", false},
	{"§0", false},
	{"§1", false},
	{"§2", false},
	{"§3", false},
	{"§4", false},
	{"§5", false},
	{"§6", false},
	{"§7", false},
	{"§8", false},
	{"§9", false},
	{"§a", false},
	{"§b", false},
	{"§c", false},
	{"§d", false},
	{"§e", false},
	{"§f", false},
	{"§k", true},
	{"§l", true},
	{"§m", true},
	{"§n", true},
	{"§o", true},
	{"§r", true},
}

func TestIsFormat(t *testing.T) {
	for index, tt := range isFormatTests {
		out := IsFormat(tt.in)
		if out != tt.out {
			t.Fatalf("%d. Failed. IsFormat(%s) => %v not %v.\n", index, tt.in, tt.out, out)
		}
	}
}

var translateTests = []struct {
	in, r, out string
}{
	{"&lHola &r&cmundo", "&", "§lHola §r§cmundo"},
	{"-lHola ---a todo--- -r-cmundo", "-", "§lHola --§a todo--- §r§cmundo"},
}

func TestTranslate(t *testing.T) {
	for index, tt := range translateTests {
		out := Translate(tt.in, tt.r)
		if out != tt.out {
			t.Fatalf("%d. Failed. Expected %q not %q.\n", index, tt.out, out)
		}
	}
}
