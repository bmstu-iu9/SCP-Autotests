package main

import (
	"testing"
)

func TestCheckEntry(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/entry_point_tests.json", 0)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests.unitTests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			var refalProgram string
			if res, _, _, err := checkFile(test.TestCase[0], &refalProgram); err != nil {
				t.Error(err)
			} else if test.Expected[0] != string(res) {
				t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected[0], res)
			}
		})
	}
}
