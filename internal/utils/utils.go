package utils

import (
	"errors"
	"ova-conversation-api/internal/domain"
)

func MakeSliceOfSlices(source []int, num int) ([][]int, error) {
	if num <= 0 {
		return nil, errors.New("invalid value of the num parameter")
	}

	result := make([][]int, 0, len(source)/num+1)
	prev := 0

	for prev+num < len(source) {
		result = append(result, source[prev:prev+num])
		prev += num
	}
	result = append(result, source[prev:])

	return result, nil
}

func MakeSliceOfSlicesConversation(source []domain.Conversation, num int) ([][]domain.Conversation, error) {
	if num <= 0 {
		return nil, errors.New("invalid value of the num parameter")
	}

	result := make([][]domain.Conversation, 0, len(source)/num+1)
	prev := 0

	for prev+num < len(source) {
		result = append(result, source[prev:prev+num])
		prev += num
	}
	result = append(result, source[prev:])

	return result, nil
}

func ReplaceKeyValue(source map[int]string) (map[string]int, error) {
	result := make(map[string]int)

	for key, value := range source {
		if _, ok := result[value]; ok {
			return nil, errors.New("same keys in the result map")
		}
		result[value] = key
	}

	return result, nil
}

func MakeSliceToMapConversation(source []domain.Conversation) (map[uint64]domain.Conversation, error) {
	result := make(map[uint64]domain.Conversation)

	for _, value := range source {
		if _, ok := result[value.ID]; ok {
			return nil, errors.New("same keys in the result map")
		}

		result[value.ID] = value
	}

	return result, nil
}

func FilterSlice(source []int, filter []int) []int {
	filterSet := make(map[int]int)

	for _, value := range filter {
		filterSet[value] = 0
	}

	result := make([]int, 0)

	for _, value := range source {
		if _, ok := filterSet[value]; !ok {
			result = append(result, value)
		}
	}

	return result
}
