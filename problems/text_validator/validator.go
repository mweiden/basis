package text_validator

type Validator interface {
	Validate(next rune) bool
	Observe(next rune)
}
