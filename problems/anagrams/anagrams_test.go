package anagrams

import (
	"testing"
	"reflect"
)

func TestAnagrams_makeKey(t *testing.T) {
	t.Parallel()
	expected := "abde4himnor2t3w2"
	result := makeKey("matthew robert weiden")
	if result != expected {
		t.Errorf("%v != %v", result, expected)
	}
}

func TestAnagrams(t *testing.T) {
	t.Parallel()
	expected := map[string][]string{
		"act": []string{"cat", "tac"},
		"dgo": []string{"dog"},
		"eoprst": []string{"presto", "repost"},
		"aefors2tw": []string{"softwares", "swears oft"},
	}
	corpus := []string{
		"cat",
		"tac",
		"dog",
		"presto",
		"repost",
		"softwares",
		"swears oft",
	}
	result := Anagrams(corpus)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("%v != %v", result, expected)
	}
}
