package main

import (
	"embed"
	"fmt"
	"os"
	"sort"
)

//go:embed words_alpha.txt
var worldList embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	store, err := loadWordStore("words_alpha.txt")
	if err != nil {
		return err
	}

	query := &wordleQuery{
		present:    append(anyPos("or"), atPos('n', 4), atPos('a', 3)),
		notPresent: append(anyPos("eusgcvb"), atPos('o', 0), atPos('a', 0), atPos('a', 1), atPos('r', 1), atPos('r', 2), atPos('o', 2)),
	}
	result := store.Execute(query)

	for _, w := range result {
		fmt.Println(w)
	}
	fmt.Println("Found", len(result), "words")
	printHist(result, query.present)
	return nil
}

func printHist(word []string, ignore []posChar) {
	type pair struct {
		key   rune
		count int
	}

	data := map[rune]int{}
	for _, w := range word {
		data = hist(w, data)
	}
	pairs := []pair{}
	for k := range data {
		pairs = append(pairs, pair{k, data[k]})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	ignoreMap := map[rune]struct{}{}
	for _, c := range ignore {
		ignoreMap[c.ch] = struct{}{}
	}
	for i := range pairs {
		if _, ok := ignoreMap[pairs[i].key]; ok {
			continue
		}
		fmt.Println(string(pairs[i].key), pairs[i].count)
	}
}
func hist(word string, seed map[rune]int) map[rune]int {
	for _, ch := range word {
		if c, ok := seed[ch]; ok {
			seed[ch] = c + 1
		} else {
			seed[ch] = 1
		}
	}

	return seed
}
