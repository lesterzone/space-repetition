package main

import (
	"sort"
	"strings"
	"time"
)

// Review stores a given answer for a Card and based on the answer we will
// compute correct/wrong answers
type Review struct {
	// Answer we will use this to compute against Card.Answer and get correct/
	// wrong values instead of adding a Correct: Boolean value to the struct.
	Answer    string `json:"answer,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// Card main object to structure data based on next review and reviews
type Card struct {
	Source     string   `json:"source,omitempty"`
	Target     string   `json:"target,omitempty"`
	Answer     string   `json:"answer,omitempty"` // used only to compute values
	Reviews    []Review `json:"reviews,omitempty"`
	NextReview int64    `json:"next_review"`
}

// Cards is a collection of Card
type Cards []Card

// Sort cards from bigger to smaller in terms of NextReview
func Sort(cards Cards) Cards {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].NextReview > cards[j].NextReview
	})

	return cards
}

// NextCards will sort cards and get reviewed cards first and then, cards
// without any reviews, so we practice with space repetition but include new
// cards after the already reviewed one.
func NextCards(cards Cards, timestamp int64) Cards {
	Sort(cards)
	var result Cards

	for _, card := range cards {
		// cards already reviewed, sorted
		reviewed := card.NextReview <= timestamp && card.NextReview > 0
		if reviewed {
			result = append(result, card)
			continue
		}

		if card.NextReview == 0 {
			result = append(result, card)
			continue
		}

		// let's append the rest of the cards without review
	}

	return result
}

// NextReview will return when a given card needs to be reviewed again.
// this is the heart of the algorithm
func NextReview(card Card, now int64) int64 {
	// if nextReview == zero, then let's add a day to current execution
	if card.NextReview == 0 {
		return AddDay(time.Unix(now, 0).Unix())
	}

	// when there are no reviews, let's just add 24h to the next review
	if len(card.Reviews) == 0 {
		return AddDay(time.Unix(now, 0).Unix())

	}

	review := card.Reviews[len(card.Reviews)-1]
	// correct := strings.ToLower(review.Answer) == strings.ToLower(card.Target)
	correct := strings.EqualFold(review.Answer, card.Target)
	var nextReview int64 = 0

	if correct {
		timeRange := now - review.CreatedAt
		// add 50% of time between last correct answer and current execution.
		// plus the time range between last corrent answer and current one.
		// this will increment the space repetition
		nextReview = (timeRange / 2) + timeRange + now
	}

	if !correct {
		timeRange := now - review.CreatedAt

		// add 25% of time between last answer and current execution.
		// this will ensure we repeat the card soon.
		nextReview = (timeRange / 4) + now
	}

	// what happens if now is less than 0 ?
	// we have corrupted data.

	// by default let's add a day to current execution time
	return nextReview
}

// AddDay will add 24h to given timestamp
func AddDay(timestamp int64) int64 {
	return time.Unix(timestamp, 0).Add(time.Hour * 24).Unix()
}
