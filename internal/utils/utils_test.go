package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"ova-conversation-api/internal/domain"
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

func TestMakeSliceOfSlicesConversation(t *testing.T) {
	c1 := domain.Conversation{ID: 1, UserID: 1, Text: "c1", Date: time.Now()}
	c2 := domain.Conversation{ID: 2, UserID: 2, Text: "c2", Date: time.Now()}
	c3 := domain.Conversation{ID: 3, UserID: 3, Text: "c3", Date: time.Now()}
	c4 := domain.Conversation{ID: 4, UserID: 4, Text: "c4", Date: time.Now()}
	c5 := domain.Conversation{ID: 5, UserID: 5, Text: "c5", Date: time.Now()}
	c6 := domain.Conversation{ID: 6, UserID: 6, Text: "c6", Date: time.Now()}

	type testType struct {
		source []domain.Conversation
		num    int
		result [][]domain.Conversation
	}

	var tests = []testType{
		{[]domain.Conversation{c1, c2, c3, c4, c5, c6}, 4, [][]domain.Conversation{{c1, c2, c3, c4}, {c5, c6}}},
		{[]domain.Conversation{c1, c2, c3, c4, c5, c6}, 2, [][]domain.Conversation{{c1, c2}, {c3, c4}, {c5, c6}}},
		{[]domain.Conversation{c1, c2, c3, c4, c5, c6}, 6, [][]domain.Conversation{{c1, c2, c3, c4, c5, c6}}},
		{[]domain.Conversation{c1, c2, c3, c4, c5, c6}, 100, [][]domain.Conversation{{c1, c2, c3, c4, c5, c6}}},
	}

	for _, value := range tests {
		actual, _ := MakeSliceOfSlicesConversation(value.source, value.num)
		require.Equal(t, value.result, actual)
	}
}

func TestMakeSliceToMapConversation(t *testing.T) {
	c1 := domain.Conversation{ID: 1, UserID: 1, Text: "c1", Date: time.Now()}
	c2 := domain.Conversation{ID: 2, UserID: 2, Text: "c2", Date: time.Now()}
	c3 := domain.Conversation{ID: 3, UserID: 3, Text: "c3", Date: time.Now()}

	type testType struct {
		source []domain.Conversation
		result map[uint64]domain.Conversation
	}

	var tests = []testType{
		{[]domain.Conversation{c1, c2, c3}, map[uint64]domain.Conversation{1: c1, 2: c2, 3: c3}},
		{[]domain.Conversation{c1, c2, c2}, nil},
	}

	for _, value := range tests {
		actual, _ := MakeSliceToMapConversation(value.source)
		require.Equal(t, value.result, actual)
	}
}
