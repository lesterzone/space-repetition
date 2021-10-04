package main

import "strings"

func memoryStore(cards Cards) map[string]Card {
	store := make(map[string]Card)
	for _, item := range cards {
		key := strings.ToLower(item.Source + item.Target)
		store[key] = item
	}
	return store
}
