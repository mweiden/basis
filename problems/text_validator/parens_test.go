package text_validator

import (
	"testing"
)

func TestParensValidator_Validate(t *testing.T) {
	t.Parallel()
	var validator ParensValidator

	result := validator.Validate('a')
	expected := true
	if result != expected {
		t.Errorf("%v != %v", expected, result)
	}

	result = validator.Validate('}')
	expected = false
	if result != expected {
		t.Errorf("%v != %v", expected, result)
	}

	validator.Observe('{')

	result = validator.Validate('a')
	expected = true
	if result != expected {
		t.Errorf("%v != %v", expected, result)
	}

	result = validator.Validate('(')
	expected = true
	if result != expected {
		t.Errorf("%v != %v", expected, result)
	}

	result = validator.Validate(')')
	expected = false
	if result != expected {
		t.Errorf("%v != %v", expected, result)
	}

	result = validator.Validate('}')
	expected = true
	if result != expected {
		t.Errorf("%v != %v", expected, result)
	}
}
