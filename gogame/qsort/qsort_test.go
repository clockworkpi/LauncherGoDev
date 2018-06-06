package qsort

import (
	"reflect"
	"testing"
)

var testData = []struct {
	input          []int
	expectedOutput []int
}{
	{[]int{}, []int{}},
	{[]int{42}, []int{42}},
	{[]int{42, 23}, []int{23, 42}},
	{[]int{23, 42, 32, 64, 12, 4}, []int{4, 12, 23, 32, 42, 64}},
}

func TestQuickSort(t *testing.T) {
	for _, testCase := range testData {
		actual := QuickSort(testCase.input)
		expected := testCase.expectedOutput

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%v != %v\n", actual, expected)
		}
	}
}
