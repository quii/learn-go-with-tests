package main

import "testing"

type CallSpy []string

func (c *CallSpy) Fn(input string) {
	*c = append(*c, input)
}

func TestWalk(t *testing.T) {

	expected := "Chris"

	x := struct {
		Name string
	}{expected}

	var fnSpy CallSpy

	walk(x, fnSpy.Fn)

	if len(fnSpy) != 1 {
		t.Errorf("wrong number of calls to CallSpy, got %d want %d", len(fnSpy), 1)
	}

}
