package main

import (
	"encoding/json"
	"testing"
)

func TestCards(t *testing.T) {
	testCases := []struct {
		desc string
		in   int64
		out  int64
	}{
		{desc: "zero input", in: 0, out: 86400},
		{desc: "half day", in: 86400, out: 172800},
	}
	for _, item := range testCases {
		t.Run(item.desc, func(t *testing.T) {
			actual := AddDay(item.in)
			expected := item.out

			if actual != expected {
				t.Fatalf("Expects %d to equal %d", actual, expected)
			}
		})
	}
}

func TestNextReview(t *testing.T) {
	testCases := []struct {
		desc string
		in   Card
		out  int64
	}{
		{
			desc: "next review is zero",
			out:  86410, // expect one day for the next review
			in:   Card{Target: "hi", Reviews: []Review{}},
		},
		{
			desc: "no reviews",
			out:  86410, // expect one day for the next review
			in: Card{
				Target:     "hi",
				NextReview: 1,
				Reviews:    []Review{},
			},
		},
		{
			desc: "correct answer first time",
			out:  25, // expect 50% of the difference between last correct answer
			in: Card{
				Target:     "hola",
				Answer:     "hola",
				NextReview: 1,
				Reviews: []Review{
					{Answer: "hola", CreatedAt: 0},
				},
			},
		},
		{
			desc: "correct answer second time",
			out:  22, // expect 50% of the difference between last correct answer + difference
			in: Card{
				Target:     "hola",
				Answer:     "hola",
				NextReview: 1,
				Reviews: []Review{
					{Answer: "hola", CreatedAt: 0},
					{Answer: "hola", CreatedAt: 2},
				},
			},
		},
		{
			desc: "inccorrect answer first time",
			out:  12, // expect 25% difference between last answer
			in: Card{
				Target:     "hola",
				Answer:     "no",
				NextReview: 1,
				Reviews: []Review{
					{Answer: "no", CreatedAt: 0},
				},
			},
		},
		{
			desc: "inccorrect answer first time",
			out:  11, // expect 25% difference between last answer
			in: Card{
				Target:     "hola",
				Answer:     "no",
				NextReview: 1,
				Reviews: []Review{
					{Answer: "no", CreatedAt: 2},
					{Answer: "no", CreatedAt: 5},
				},
			},
		},
	}

	for _, item := range testCases {
		t.Run(item.desc, func(t *testing.T) {
			expected := item.out
			// 10 is used only to easy test, simulating 'now' in seconds
			actual := NextReview(item.in, 10)

			if actual != item.out {
				t.Fatalf("For %s expected: %d got: %d", item.desc, expected, actual)
			}
		})
	}
}

func TestSortCards(t *testing.T) {
	testCases := []struct {
		desc string
		in   Cards
		out  string
	}{
		{
			desc: "empty cards",
			in:   Cards{{NextReview: 0}},
			out:  `[{"next_review":0}]`,
		},
		{
			desc: "single card",
			in:   Cards{{NextReview: 0, Answer: ""}},
			out:  `[{"next_review":0}]`,
		},
		{
			desc: "multiple cards",
			in: Cards{
				{Source: "", Target: "", Answer: "", NextReview: 0},
				{Source: "", Target: "", Answer: "", NextReview: 1},
				{Source: "", Target: "", Answer: "", NextReview: 3},
			},
			out: `[{"next_review":3},{"next_review":1},{"next_review":0}]`,
		},
	}

	for _, item := range testCases {
		t.Run(item.desc, func(t *testing.T) {
			actual, _ := json.Marshal(Sort(item.in))
			expected := item.out

			if string(actual) != item.out {
				t.Fatalf("For %s expected: %s actual: %s", item.desc, expected, actual)
			}
		})
	}
}

func TestNextCards(t *testing.T) {
	testCases := []struct {
		desc      string
		in        Cards
		timestamp int64
		out       string
	}{
		// {
		// 	desc:      "empty cards",
		// 	in:        Cards{},
		// 	timestamp: 30,
		// 	out:       "",
		// },

		{
			desc:      "single card",
			timestamp: 30,
			in: Cards{
				{Source: "", Target: "", Answer: "", NextReview: 0},
			},
			out: `[{"next_review":0}]`,
		},
		{
			desc:      "one card off",
			timestamp: 30,
			in: Cards{
				{NextReview: 0},
				{NextReview: 31},
				{NextReview: 28},
			},
			out: `[{"next_review":28},{"next_review":0}]`,
		},
		{
			desc:      "one card extra",
			timestamp: 30,
			in: Cards{
				{NextReview: 0},  // should appear after any review > 0
				{NextReview: 0},  // should appear after any review > 0
				{NextReview: 31}, // should not be included since we are using 30
				{NextReview: 28},
			},
			out: `[{"next_review":28},{"next_review":0},{"next_review":0}]`,
		},
	}

	for _, item := range testCases {
		t.Run(item.desc, func(t *testing.T) {
			actual, _ := json.Marshal(NextCards(item.in, item.timestamp))
			expected := item.out

			if string(actual) != item.out {
				t.Fatalf("For %s \nexpected: %s \nactual: %s", item.desc, expected, actual)
			}
		})
	}
}
