package main

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	assert := func(s [][]string, expected []string) {
		got := Intersection(s...)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected: %q, but got %q", expected, got)
		}
	}

	assert([][]string{{"hoge", "fuga"}, {"hoge", "poyo"}}, []string{"hoge"})
	assert([][]string{{"hoge", "fuga"}, {"poyopoyo", "poyo"}}, []string{})
	assert([][]string{{"hoge", "fuga"}, {"hoge", "poyo"}, {"hoge", "nya"}}, []string{"hoge"})
}
