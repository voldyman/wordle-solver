package main

import (
	"reflect"
	"testing"
)

var store *WordStore

func init() {
	var err error
	store, err = loadWordStore("words_alpha.txt")
	if err != nil {
		panic(err)
	}
}

// Old test to check for query is same as Lookup
// lookup method lost to time
// func TestAnyPosContainsQuery(t *testing.T) {
// 	query := &wordleQuery{
// 		present: anyPos('t'), anyPos('a'), anyPos('b'), anyPos('e'),
// 	}
// 	qResult := store.Execute(query)
// 	lResult := store.lookup('t', 'a', 'b', 'e')
// 	if !reflect.DeepEqual(qResult, lResult) {
// 		t.Errorf("Error: query result != lookup result; qResult: %v, lResult: %v", qResult, lResult)
// 	}
// }

func TestAnyPosNotPresentQuery(t *testing.T) {
	query := &wordleQuery{
		present:    anyPos("tabe"),
		notPresent: anyPos("urhgjn"),
	}
	result := store.Execute(query)
	if len(result) != 29 {
		t.Errorf("Error: results returned = %d; expected %d", len(result), 29)
	}
}

func TestUnion(t *testing.T) {
	testcases := []struct {
		lhs    []int
		rhs    []int
		result []int
	}{
		{[]int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
		{[]int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{[]int{}, []int{2, 3, 4}, []int{2, 3, 4}},
		{[]int{1, 2, 3}, []int{2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, tc := range testcases {
		actual := union(tc.lhs, tc.rhs)
		if !reflect.DeepEqual(actual, tc.result) {
			t.Errorf("Error: union(%+v, %+v) = %v; expected %v", tc.lhs, tc.rhs, actual, tc.result)
		}
	}
}
func TestDifference(t *testing.T) {
	testcases := []struct {
		lhs    []int
		rhs    []int
		result []int
	}{
		{[]int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{[]int{1, 2, 3}, []int{2}, []int{1, 3}},
		{[]int{1, 2, 3, 4, 6}, []int{1, 2}, []int{3, 4, 6}},
	}

	for _, tc := range testcases {
		actual := difference(tc.lhs, tc.rhs)
		if !reflect.DeepEqual(actual, tc.result) {
			t.Errorf("Error: difference(%v, %v) = %v; expected %v", tc.lhs, tc.lhs, actual, tc.result)
		}
	}
}

func TestRank(t *testing.T) {
	testcases := []struct {
		input  []string
		output []string
	}{
		{[]string{"table", "tiles", "panic", "manic"}, []string{"tiles", "table", "manic", "panic"}},
	}
	for _, tc := range testcases {
		output := rank(tc.input)
		if !reflect.DeepEqual(output, tc.output) {
			t.Errorf("Error: rank(%+v) = %v; expected: %v", tc.input, output, tc.output)
		}
	}
}
