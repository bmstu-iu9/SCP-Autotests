package main

import (
	"strings"
	"testing"
)

func TestCheckParameters(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/check_parameters_tests.json", 0)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests.unitTests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := checkParameters(test.TestCase[0]); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s",
						test.Expected[0], err.Error())
				}
			} else {
				for i, token := range res {
					if token != test.Expected[i] {
						t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected, res)
					}
				}
			}
		})
	}
}

func TestSplitToTokens(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/split_to_tokens_tests.json", 0)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests.unitTests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := splitToTokens(test.TestCase[0]); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s",
						test.Expected[0], err.Error())
				}
			} else {
				for i, token := range res {
					if token != test.Expected[i] {
						t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected, res)
					}
				}
			}
		})
	}
}

func TestCheckTokens(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/check_tokens_tests.json", 0)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests.unitTests {
		t.Run(strings.Join(test.TestCase, " "), func(t *testing.T) {
			if err = checkTokens(test.TestCase); err != nil && test.Expected[0] != err.Error() {
				t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s",
					test.Expected[0], err.Error())
			}
		})
	}
}

func TestCheckVariables(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/check_variables_tests.json", 0)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests.unitTests {
		t.Run(strings.Join(test.TestCase, " "), func(t *testing.T) {
			coincidences := make(map[string]bool)
			if err = checkVariables(test.TestCase, &coincidences); err != nil && test.Expected[0] != err.Error() {
				t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s",
					test.Expected[0], err.Error())
			}
		})
	}
}
