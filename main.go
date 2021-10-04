package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	config := Config{
		PersonalFolder: os.Getenv("PERSONAL_FOLDER"),
		DatabasePath:   "/memcards.json",
	}

	// if path is provided, import the CSV
	if len(os.Args) == 2 {
		csvPath := os.Args[1]

		if csvPath != "" {
			importCards(config, csvPath)
			return
		}
	}

	var cards Cards
	readDB(config, &cards)
	allCards := memoryStore(cards)

	cards = NextCards(cards, time.Now().Unix())
	store := memoryStore(cards)

	newStore := present(config, store)

	for key, value := range newStore {
		allCards[key] = value
	}

	updateDB(config, allCards)
}

type Config struct {
	PersonalFolder string
	DatabasePath   string
}

func readDB(config Config, cards *Cards) {
	data, err := ioutil.ReadFile(config.PersonalFolder + config.DatabasePath)

	if err != nil {
		fmt.Println("if not file found, let's handle the creation of the file")
		fmt.Print(err)
		return
	}

	err = json.Unmarshal(data, &cards)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func importCards(config Config, csvPath string) {
	var cards Cards
	readDB(config, &cards)

	file, err := os.Open(csvPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	store := memoryStore(cards)
	toUpdate := make(map[string]Card)

	// remember to close the file at the end of the program
	defer file.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		key := strings.ToLower(row[0] + row[1])
		_, ok := store[key]

		if !ok {
			toUpdate[key] = Card{Source: row[0], Target: row[1]}
		}
	}

	if len(toUpdate) == 0 {
		return
	}

	for key, value := range toUpdate {
		_, ok := store[key]
		if !ok {
			store[key] = value
		}
	}

	updateDB(config, store)
}

func updateDB(config Config, store map[string]Card) {
	file := config.PersonalFolder + config.DatabasePath
	_, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("if not file found, let's handle the creation of the file")
		fmt.Print(err)
		return
	}

	var list Cards
	for _, value := range store {
		list = append(list, value)
	}

	result, err := json.Marshal(list)

	if err != nil {
		fmt.Println("error marshaling list")
		return
	}

	err = ioutil.WriteFile(file, result, 0644)
	if err != nil {
		fmt.Println("error writing to file")
		return
	}
}

// present
// will display cards one by one waiting for user input.
// to quit, we can use q
// TODO: quit with ctrl + c and update store.
func present(config Config, store map[string]Card) map[string]Card {
	newStore := make(map[string]Card)

	for _, card := range store {
		fmt.Println(card.Source)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nYour answer?\n")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", 1)

		if text == "q" {
			updateDB(config, store)
			return newStore
		}

		now := time.Now().Unix()

		card.Reviews = append(card.Reviews, Review{
			CreatedAt: now, Answer: text,
		})

		card.NextReview = NextReview(card, now)
		key := strings.ToLower(card.Source + card.Target)
		newStore[key] = card
	}

	return newStore
}
