package config

import (
	"testing"
)

var parseTests = []struct {
	in  string
	out map[string]string
	err bool
}{
	{
		in:  "",
		out: map[string]string{},
		err: false,
	},
	{
		in:  "\n",
		out: map[string]string{},
		err: false,
	},
	{
		in:  "a:\nb:\n",
		out: map[string]string{},
		err: false,
	},
	{
		in:  "a:\n b:\n",
		out: map[string]string{},
		err: false,
	},
	{
		in:  "a:\n b:2\nc:3",
		out: map[string]string{"a.b": "2", "c": "3"},
		err: false,
	},
	{
		in:  "a:\n b:\n  c: 2\n  d: 3\n e:\n  f: 5\ng:\n h: 7",
		out: map[string]string{"a.b.c": "2", "a.b.d": "3", "a.e.f": "5", "g.h": "7"},
		err: false,
	},
}

func TestParse(t *testing.T) {
	for index, tt := range parseTests {
		out := New()
		out.Parse(tt.in)

		for tk, tv := range tt.out {
			v := out.Get(tk)
			if v != tv {
				t.Fatalf("%d. map[%q] expects %q, got %q.", index, tk, tv, v)
			}
		}
	}
}
