package main

import (
	"reflect"
	"testing"
)

var productTests = []struct {
	src [][]string
	dst [][]int
}{
	{
		src: [][]string{
			{"a", "b", "c"},
			{},
			{"1", "2"},
		},
		dst: [][]int{},
	},
	{
		src: [][]string{
			{"a", "b", "c"},
		},
		dst: [][]int{
			{0},
			{1},
			{2},
		},
	},
	{
		src: [][]string{
			{"h", "j", "k", "l"},
			{"a", "b", "c"},
			{"1", "2"},
			{"-"},
		},
		dst: [][]int{
			{0, 0, 0, 0},
			{0, 0, 1, 0},
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 2, 0, 0},
			{0, 2, 1, 0},
			{1, 0, 0, 0},
			{1, 0, 1, 0},
			{1, 1, 0, 0},
			{1, 1, 1, 0},
			{1, 2, 0, 0},
			{1, 2, 1, 0},
			{2, 0, 0, 0},
			{2, 0, 1, 0},
			{2, 1, 0, 0},
			{2, 1, 1, 0},
			{2, 2, 0, 0},
			{2, 2, 1, 0},
			{3, 0, 0, 0},
			{3, 0, 1, 0},
			{3, 1, 0, 0},
			{3, 1, 1, 0},
			{3, 2, 0, 0},
			{3, 2, 1, 0},
		},
	},
}

func TestProduct(t *testing.T) {
	for _, test := range productTests {
		expect := test.dst
		actual := make([][]int, 0)
		for indexes := range Product(test.src) {
			actual = append(actual, indexes)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Product(%q) returns %d, want %d",
				test.src, actual, expect)
		}
	}
}
