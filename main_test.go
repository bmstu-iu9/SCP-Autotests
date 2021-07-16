package main

import (
	"strings"
	"testing"
)

func TestCheckEntry(t *testing.T) {
	tests, err := getTestsFromFile("tests/entry_point_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := checkFile(test.TestCase[0]); err != nil {
				t.Error(err)
			} else {
				if test.Expected[0] != strings.TrimSpace(string(res)) {
					t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected[0], res)
				}
			}
		})
	}
}
