package phone_numbers

import (
	"reflect"
	"testing"
)

func TestWith_PossibleWordsForPhoneNumber(t *testing.T) {
	t.Parallel()

	dictionary := Dictionary{[]string{"BEE", "ADD", "WING", "STING"}}

	result := With(dictionary.ToPhoneNumberMap()).PossibleWordsForPhoneNumber(
		[]int{1, 8, 0, 0, 2, 3, 3, 9, 4, 6, 4},
	)
	expected := []string{"1-800-BEE-WING", "1-800-ADD-WING"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("%v != %v", result, expected)
	}
}
