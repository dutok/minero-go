package material

import (
	"testing"
)

var fromNameTests = []struct {
	in  string
	out Material
}{
	{"", Air},
	{"Air", Air},
	{"Bedrock", Bedrock},
	{"FishingRod", ItemFishingRod},
	{"ItemFishingRod", ItemFishingRod},
	{"Far", DiscFar},
	{"DiscFar", DiscFar},
}

func TestFromName(t *testing.T) {
	for index, tt := range fromNameTests {
		out := FromName(tt.in)
		if tt.out != out {
			t.Fatalf("%d. FromName(%q): Expected %s, got %s.\n", index+1, tt.in, tt.out, out)
		}
	}
}

var fromIdTests = []struct {
	in  int
	out Material
}{
	{-1, Air},
	{0, Air},
	{346, ItemFishingRod},
	{2260, DiscFar},
}

func TestFromId(t *testing.T) {
	for index, tt := range fromIdTests {
		out := FromId(tt.in)
		if tt.out != out {
			t.Fatalf("%d. FromName(%q): Expected %s, got %s.\n", index+1, tt.in, tt.out, out)
		}
	}
}

func TestFromNameRange(t *testing.T) {
	for i := MinBlock; i < MaxBlock; i++ {
		out := FromId(i)
		if !out.IsBlock() {
			t.Fatalf("FromId(%d): Not a block.", i)
		}
	}
	for i := MinItem; i < MaxItem; i++ {
		out := FromId(i)
		if !out.IsItem() {
			t.Fatalf("FromId(%d): Not an item.", i)
		}
	}
	for i := MinDisc; i < MaxDisc; i++ {
		out := FromId(i)
		if !out.IsDisc() {
			t.Fatalf("FromId(%d): Not a disc.", i)
		}
	}
}
