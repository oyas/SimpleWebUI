package main

import (
	"../parser"
	"strings"
	"testing"
)

var testParse_strings = [][]string{
	{"# Hello", "# Hello"},
	{"##Input", "##Input"},
	{"#Input", "<label> <input type=\"text\" class=\"form-control\" name=\"\"></label>"},
	{"#Submit", "<button type=\"submit\" class=\"btn btn-default\">submit</button>"},
}

func TestParse(t *testing.T) {
	for _, testcase := range testParse_strings {
		input := testcase[0]
		expected := testcase[1]
		t.Logf("\"%s\" => \"%s\"", input, expected)
		actual := strings.Trim(parser.Parse(input), " \n")
		if actual != expected {
			t.Errorf("Fail. Output is \"%s\"", actual)
		}
	}
}
