package config

import (
	"testing"
)

var parseTests = []struct {
	in  string
	out map[string]string
	err error
}{
	{
		in:  "",
		out: map[string]string{},
		err: ErrEmpty,
	},
	{
		in:  "\n",
		out: map[string]string{},
		err: nil,
	},
	{
		in:  "a:\nb:\n",
		out: map[string]string{},
		err: nil,
	},
	{
		in:  "a:\n b:\n",
		out: map[string]string{},
		err: nil,
	},
	{
		in:  "a:\n b:2\nc:3",
		out: map[string]string{"a.b": "2", "c": "3"},
		err: nil,
	},
	{
		in:  "a:\n b:\n  c: 2\n  d: 3\n e:\n  f: 5\ng:\n h: 7",
		out: map[string]string{"a.b.c": "2", "a.b.d": "3", "a.e.f": "5", "g.h": "7"},
		err: nil,
	},
	{
		in:  "#hello!\na:\n b:\n  c: 2 # c is for crafting!",
		out: map[string]string{"a.b.c": "2"},
		err: nil,
	},
	{
		in:  "#hello!\na:\n b:\n  c: -\nd: 2 # d is for destinty!",
		out: map[string]string{"a.b.c": "", "d": "2"},
		err: nil,
	},
}

func TestParse(t *testing.T) {
	for index, tt := range parseTests {
		out := New()
		err := out.Parse(tt.in)

		if err != tt.err {
			t.Fatalf("%d. expecting error: %v, got %v.", index+1, tt.err, err)
		}

		for tk, tv := range tt.out {
			v := out.Get(tk)
			if v != tv {
				t.Fatalf("%d. map[%q] expects %q, got %q.", index+1, tk, tv, v)
			}
		}
	}
}
