package main

type Regex struct {
	pattern Parser
	fn      func(s string) bool
}
