package anagrams

import (
	"testing"
	"reflect"
)

func TestAnagrams_makeKey(t *testing.T) {
	t.Parallel()
	expected := "a1m1t2"
	result := makeKey("matt")
	if result != expected {
		t.Errorf("%v != %v", result, expected)
	}
}

func TestAnagrams(t *testing.T) {
	t.Parallel()
	expected := map[string][]string{
		"a1c1t1": []string{"cat", "tac"},
		"d1g1o1": []string{"dog"},
		"e1o1p1r1s1t1": []string{"presto", "repost"},
		"a1e1f1o1r1s2t1w1": []string{"softwares", "swears oft"},
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
