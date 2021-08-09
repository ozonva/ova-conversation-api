package utils

import (
	"reflect"
	"testing"
)

func TestMakeSliceOfSlices(t *testing.T) {
	type testType struct {
		source []int
		num    int
		result [][]int
	}

	var tests = []testType{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 4, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 2, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 100, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}},
	}

	for _, value := range tests {
		actual, _ := MakeSliceOfSlices(value.source, value.num)
		if !reflect.DeepEqual(value.result, actual) {
			t.Errorf("Expected: %v, actual: %v", value.result, actual)
		}
	}
}

func TestReplaceKeyValue(t *testing.T) {
	type testType struct {
		source map[int]string
		result map[string]int
	}

	var tests = []testType{
		{map[int]string{1: "one", 2: "two", 3: "three"}, map[string]int{"one": 1, "two": 2, "three": 3}},
		{map[int]string{1: "one", 2: "two", 3: "two"}, nil},
	}

	for _, value := range tests {
		actual, _ := ReplaceKeyValue(value.source)
		if !reflect.DeepEqual(value.result, actual) {
			t.Errorf("Expected: %v, actual: %v", value.result, actual)
		}
	}
}

func TestFilterSlice(t *testing.T) {
	type testType struct {
		source []int
		filter []int
		result []int
	}

	var tests = []testType{
		{[]int{1, 2, 3, 4, 5}, []int{2, 4}, []int{1, 3, 5}},
		{[]int{1, 2, 3, 4, 5}, []int{}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}, []int{}},
		{[]int{1, 2, 3, 4, 5}, []int{6, 7}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 1, 1}, []int{1}, []int{}},
	}

	for _, value := range tests {
		actual := FilterSlice(value.source, value.filter)
		if !reflect.DeepEqual(value.result, actual) {
			t.Errorf("Expected: %v, actual: %v", value.result, actual)
		}
	}
}
