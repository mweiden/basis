package main

import (
	"errors"
	"fmt"
)

func newParserError(intrfc interface{}, msg string) error {
	return errors.New(fmt.Sprintf("Parser %v error: %s", intrfc, msg))
}

type ParserResult struct {
	offset          int
	capturedStrings []string
}

type Parser interface {
	Parse(s string) (error, ParserResult)
}

type Captor struct {
	p Parser
}

func (c Captor) Parse(s string) (error, ParserResult) {
	err, result := c.p.Parse(s)
	result.capturedStrings = append(result.capturedStrings, s[:result.offset])
	return err, result
}

// CharParser
type CharParser struct {
	r rune
}

func (cp CharParser) Parse(s string) (error, ParserResult) {
	if rune(s[0]) == cp.r {
		return nil, ParserResult{offset: 1}
	} else {
		return newParserError(cp, "no match"), ParserResult{offset: 0}
	}
}

// CharRangeParser
type CharRangeParser struct {
	rStart rune
	rEnd   rune
}

func (crp CharRangeParser) Parse(s string) (error, ParserResult) {
	if crp.rStart <= rune(s[0]) && rune(s[0]) <= crp.rEnd {
		return nil, ParserResult{offset: 1}
	} else {
		return newParserError(crp, "no match"), ParserResult{offset: 0}
	}
}

// OptParser
type OptParser struct {
	s string
}

func (op OptParser) Parse(s string) (error, ParserResult) {
	if s[:len(op.s)] == op.s {
		return nil, ParserResult{offset: len(op.s)}
	} else {
		return nil, ParserResult{offset: 0}
	}
}

// RepeatParser
type RepeatParser struct {
	s string
}

func (rp RepeatParser) Parse(s string) (error, ParserResult) {
	offset := 0
	for s[offset:(offset+len(rp.s))] == rp.s {
		offset += len(rp.s)
	}
	return nil, ParserResult{offset: offset}
}

// RepeatRange
type RepeatRangeParser struct {
	s    string
	from int
	to   int
}

func (rrp RepeatRangeParser) Parse(s string) (error, ParserResult) {
	offset := 0
	matches := 0
	for s[offset:(offset+len(rrp.s))] == rrp.s {
		offset += len(rrp.s)
		matches += 1
	}
	if rrp.from <= matches && matches <= rrp.to {
		return nil, ParserResult{offset: offset}
	} else {
		return newParserError(rrp, "incorrect number of matches"), ParserResult{offset: 0}
	}
}

// SequenceParser
type SequenceParser struct {
	parsers []Parser
}

func (sp SequenceParser) Parse(s string) (error, ParserResult) {
	offset := 0
	for _, parser := range sp.parsers {
		if offset < len(s) {
			err, result := parser.Parse(s[offset:])
			if err != nil {
				return err, ParserResult{offset: offset}
			} else {
				offset += result.offset
			}
		} else {
			return newParserError(sp, "no string remaining"), ParserResult{offset: offset}
		}
	}
	return nil, ParserResult{offset: offset}
}

// AltParser
type AltParser struct {
	parsers []Parser
}

func (sp AltParser) Parse(s string) (error, ParserResult) {
	for _, parser := range sp.parsers {
		err, result := parser.Parse(s)
		if err == nil {
			return nil, result
		}
	}
	return newParserError(sp, "no parsers succeeded"), ParserResult{offset: 0}
}
