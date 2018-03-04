package main

import (
	"reflect"
	"testing"
)

type Expect struct {
	err    error
	result ParserResult
}

func compare(t *testing.T, exp Expect, err error, result ParserResult) {
	if !reflect.DeepEqual(err, exp.err) || !reflect.DeepEqual(result, exp.result) {
		t.Errorf("Expected (%v, %v), got (%v, %v)", exp.err, exp.result, err, result)
	}
}

func TestCharParser_Parse(t *testing.T) {
	t.Parallel()
	parser := CharParser{'H'}

	// negative result
	err, offset := parser.Parse("Guten Tag")
	exp := Expect{newParserError(parser, "no match"), ParserResult{offset: 0}}
	compare(t, exp, err, offset)

	// positive result
	err, offset = parser.Parse("Hello")
	exp = Expect{nil, ParserResult{offset: 1}}
	compare(t, exp, err, offset)
}

func TestSequenceParser_Parse(t *testing.T) {
	t.Parallel()
	parsers := []Parser{
		CharParser{'H'},
		CharParser{'e'},
		CharParser{'l'},
		CharParser{'l'},
		CharParser{'o'},
	}
	seqParser := SequenceParser{parsers}

	// string too short
	err, offset := seqParser.Parse("Hell")
	exp := Expect{newParserError(seqParser, "no string remaining"), ParserResult{offset: 4}}
	compare(t, exp, err, offset)

	// no match
	err, offset = seqParser.Parse("Guten Tag!")
	exp = Expect{newParserError(CharParser{'H'}, "no match"), ParserResult{offset: 0}}
	compare(t, exp, err, offset)

	// match
	err, offset = seqParser.Parse("Hello")
	exp = Expect{nil, ParserResult{offset: 5}}
	compare(t, exp, err, offset)
}

func TestAltParser_Parse(t *testing.T) {
	t.Parallel()
	parsers := []Parser{
		CharParser{'H'},
		CharParser{'G'},
	}
	altParser := AltParser{parsers}

	// no match
	err, offset := altParser.Parse("Servus!")
	exp := Expect{newParserError(altParser, "no parsers succeeded"), ParserResult{offset: 0}}
	compare(t, exp, err, offset)

	// match
	err, offset = altParser.Parse("Hello")
	exp = Expect{nil, ParserResult{offset: 1}}
	compare(t, exp, err, offset)
}

func TestOptParser_Parse(t *testing.T) {
	t.Parallel()
	altParser := OptParser{"Hi "}

	// match
	err, offset := altParser.Parse("Hi Jason!")
	exp := Expect{nil, ParserResult{offset: 3}}
	compare(t, exp, err, offset)

	// no match
	err, offset = altParser.Parse("Jason!")
	exp = Expect{nil, ParserResult{offset: 0}}
	compare(t, exp, err, offset)
}

func TestRepeatParser_Parse(t *testing.T) {
	t.Parallel()
	parser := RepeatParser{"Hi "}

	// no match
	err, offset := parser.Parse("Jason!")
	exp := Expect{nil, ParserResult{offset: 0}}
	compare(t, exp, err, offset)

	// match
	err, offset = parser.Parse("Hi Jason!")
	exp = Expect{nil, ParserResult{offset: 3}}
	compare(t, exp, err, offset)

	// 3 matches
	err, offset = parser.Parse("Hi Hi Hi Jason!")
	exp = Expect{nil, ParserResult{offset: 9}}
	compare(t, exp, err, offset)
}

func TestRepeatRangeParser_Parse(t *testing.T) {
	t.Parallel()
	parser := RepeatRangeParser{"Hi ", 2, 3}

	// no match
	err, offset := parser.Parse("Hi Jason!")
	exp := Expect{newParserError(parser, "incorrect number of matches"), ParserResult{offset: 0}}
	compare(t, exp, err, offset)

	// 2 matches
	err, offset = parser.Parse("Hi Hi Jason!")
	exp = Expect{nil, ParserResult{offset: 6}}
	compare(t, exp, err, offset)

	// 3 matches
	err, offset = parser.Parse("Hi Hi Hi Jason!")
	exp = Expect{nil, ParserResult{offset: 9}}
	compare(t, exp, err, offset)

	// 4 matches, no match
	err, offset = parser.Parse("Hi Hi Hi Hi Jason!")
	exp = Expect{newParserError(parser, "incorrect number of matches"), ParserResult{offset: 0}}
	compare(t, exp, err, offset)
}

func TestCaptor_Parse(t *testing.T) {
	t.Parallel()
	parser := Captor{RepeatRangeParser{"Hi ", 2, 3}}

	// 3 matches
	err, offset := parser.Parse("Hi Hi Hi Jason!")
	exp := Expect{nil, ParserResult{offset: 9, capturedStrings: []string{"Hi Hi Hi "}}}
	compare(t, exp, err, offset)
}
