package main

import (
	"strings"
	"testing"
)

func Test_Split(t *testing.T) {
	s := "123"
	res := strings.Split(s, " ")
	t.Logf("len: %d\n", len(res))
}
