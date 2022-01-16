package main

import (
	"bufio"
	"fmt"
	"sort"
)

type WordStore struct {
	all   []int // array of all results indexes
	words []string
	index [][][]int // maps char ascii code -> list of positions in a word -> position in words
}

func loadWordStore(wordFilePath string) (*WordStore, error) {
	f, err := worldList.Open(wordFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open word file: %w", err)
	}

	store := newWordStore()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if acceptableWord(line) {
			store.add(line)
		}
	}
	if err = s.Err(); err != nil {
		return nil, fmt.Errorf("unable to read word file: %w", err)
	}
	return store, nil
}

var charsetSize = ord('z') - ord('A') + 1

func acceptableWord(word string) bool {
	lookup := make([]bool, charsetSize)
	for _, ch := range word {
		if lookup[ord(ch)] {
			return false
		} else {
			lookup[ord(ch)] = true
		}
	}
	return len(word) == 5
}

func newWordStore() *WordStore {
	return &WordStore{
		words: []string{},
		index: make([][][]int, 'z'-'A'+1),
	}
}

func (s *WordStore) add(word string) {
	wordIdx := len(s.words)
	s.words = append(s.words, word)

	for i, ch := range word {
		charCode := ord(ch)
		posArray := s.index[charCode]
		if len(posArray) == 0 {
			// wordle is 5 characters long
			s.index[charCode] = make([][]int, 5)
			posArray = s.index[charCode]
		}
		posArray[i] = append(posArray[i], wordIdx)
		s.all = append(s.all, wordIdx)
	}
}

func (s *WordStore) Execute(q Query) []string {
	idxs := q.Eval(s)
	result := make([]string, len(idxs))
	for i, wordIdx := range idxs {
		result[i] = s.words[wordIdx]
	}

	return rank(result)
}

type posChar struct {
	ch  rune
	pos int
}

func atPoss(ch rune, poses ...int) []posChar {
	result := []posChar{}
	for _, pos := range poses {
		result = append(result, posChar{
			ch:  ch,
			pos: pos,
		})
	}
	return result
}
func atPos(ch rune, pos int) posChar {
	return posChar{
		ch:  ch,
		pos: pos,
	}
}
func anyPos(chs string) []posChar {
	result := make([]posChar, len(chs))
	for i, ch := range chs {
		result[i] = posChar{
			ch:  ch,
			pos: -1,
		}
	}
	return result
}

type Query interface {
	Eval(s *WordStore) []int
}

type wordleQuery struct {
	present    []posChar
	notPresent []posChar
}

func (q *wordleQuery) Eval(s *WordStore) []int {
	result := s.all
	if len(q.present) > 0 {
		result = retrieve(s, s.all, q.present, intersect)
	}
	if len(q.notPresent) > 0 {
		notPresent := retrieve(s, []int{}, q.notPresent, union)
		result = difference(result, notPresent)
	}
	return result

}

type combineFn func([]int, []int) []int

func retrieve(s *WordStore, start []int, target []posChar, combine combineFn) []int {
	for _, pch := range target {
		charCode := ord(pch.ch)
		posArray := s.index[charCode]

		posting := []int{}
		// no pos specified
		if pch.pos < 0 {
			for _, p := range posArray {
				posting = union(posting, p)
			}
		} else {
			posting = posArray[pch.pos]
		}
		start = combine(start, posting)

	}
	return start
}

func rank(words []string) []string {
	// group words by char freq
	// sort groups by char pos frequency
	// since each word has a char only once, sorting by char pos frequency should do the both above simultaneously

	freqs := [5][]int{}
	updateFreq := func(word string, i int) {
		if len(freqs[i]) == 0 {
			freqs[i] = make([]int, charsetSize)
		}
		freqs[i][ord(rune(word[i]))] += 1
	}

	for _, word := range words {
		for i := 0; i < 5; i++ {
			updateFreq(word, i)
		}
	}
	score := func(word string) int {
		return freqs[0][bord(word[0])] +
			freqs[1][bord(word[1])] +
			freqs[2][bord(word[2])] +
			freqs[3][bord(word[3])] +
			freqs[4][bord(word[4])]
	}
	sort.Slice(words, func(i, j int) bool {
		return score(words[i]) <= score(words[j])
	})
	return words
}

func bord(ch byte) int {
	return int(ch - 'A')
}
func ord(ch rune) int {
	return int(ch - 'A')
}

func unord(o int) rune {
	return rune(o + 'A')
}

func intersect(a []int, b []int) []int {
	if len(b) > len(a) {
		a, b = b, a
	}

	result := []int{}
	for _, aVal := range a {
		bIdx := sort.SearchInts(b, aVal)
		if bIdx < len(b) && b[bIdx] == aVal {
			result = append(result, aVal)
		}
	}

	return result
}

func union(a []int, b []int) []int {
	aIdx := 0
	bIdx := 0

	result := make([]int, 0, min(len(a), len(b)))
	for aIdx < len(a) && bIdx < len(b) {
		if a[aIdx] < b[bIdx] {
			result = append(result, a[aIdx])
			aIdx++
			continue
		}
		if a[aIdx] > b[bIdx] {
			result = append(result, b[bIdx])
			bIdx++
			continue
		}
		if a[aIdx] != b[bIdx] {
			panic(fmt.Sprintf("assertion failed: a[aIdx](%d) == b[bIdx](%d)", a[aIdx], b[bIdx]))
		}
		result = append(result, a[aIdx])
		aIdx++
		bIdx++
	}

	for aIdx < len(a) {
		result = append(result, a[aIdx])
		aIdx++
	}
	for bIdx < len(b) {
		result = append(result, b[bIdx])
		bIdx++
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func difference(a []int, b []int) []int {
	moveHead := func(val int, posting []int) ([]int, bool) {
		idx := sort.SearchInts(posting, val)
		found := false
		if idx < len(posting) && posting[idx] == val {
			found = true
			b = b[idx:]
		}
		if idx == len(b) {
			b = []int{}
		}
		return b, found
	}
	result := []int{}
	for _, aVal := range a {
		foundInB := false
		if len(b) > 0 {
			b, foundInB = moveHead(aVal, b)
		}
		if !foundInB {
			result = append(result, aVal)
		}
	}
	return result
}
