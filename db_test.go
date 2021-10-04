package main

import "testing"

func TestMemoryStore(t *testing.T) {
	testCases := []struct {
		desc string
		in   Cards
		out  int
	}{
		{
			desc: "Empty cards",
			in:   Cards{},
			out:  0,
		},
		{
			desc: "Single card",
			in:   Cards{{NextReview: 0}},
			out:  1,
		},
		{
			desc: "Multiple cards",
			in: Cards{
				{Source: "first", Target: "target"},
				{Source: "second", Target: "target"},
			},
			out: 2,
		},
		{
			desc: "Cards duplicated",
			in: Cards{
				{Source: "first", Target: "target"},
				{Source: "first", Target: "target"},
				{Source: "second", Target: "target"},
			},
			out: 2,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := len(memoryStore(tC.in))
			expected := tC.out

			if actual != expected {
				t.Fatalf("For %s, expexted: %d got: %d", tC.desc, expected, actual)
			}
		})
	}
}
