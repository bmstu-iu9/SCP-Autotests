package main

import (
	"testing"
)

func TestGetUnparametrizedEntryPoint(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/get_unparametrized_entry_point_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := getUnparametrizedEntryPoint(test.TestCase[0], test.TestCase[1]); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if test.Expected[0] != res {
				t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected[0], res)
			}
		})
	}
}

func TestMatchData(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/match_data_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			values := make(map[string]string)
			if err := matchData(test.TestCase[1:], &test.TestCase[0], &values); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else {
				for i := 0; i < len(test.Expected); i += 2 {
					if values[test.Expected[i]] != test.Expected[i + 1] {
						t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected[i] + "->" + test.Expected[i+1], test.Expected[i] + "->" + values[test.Expected[i]])
					}
				}
			}
		})
	}
}

func TestMatchDataHelper(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/match_data_helper_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			value := make(map[string]string)
			if err = matchDataHelper(test.TestCase[0], &test.TestCase[1], &value, false); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else {
				if test.TestCase[1] != test.Expected[0] {
					t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s %s", test.Expected, test.TestCase[1], value)
				}
				for i := 1; i < len(test.Expected); i += 2 {
					if value[test.Expected[i]] != test.Expected[i + 1] {
						t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected[i] + "->" + test.Expected[i+1], test.Expected[i] + "->" + value[test.Expected[i]])
					}
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/parse_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := parse(test.TestCase[0], 0, test.TestCase[1][0]); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if res != int(test.Expected[0][0] - '0') {
				t.Errorf("Incorrect result\nExpected: %s\nRecieved: %d", test.Expected, res)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/parse_string_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := parseString(test.TestCase[0], 0); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if res != int(test.Expected[0][0] - '0') {
				t.Errorf("Incorrect result\nExpected: %s\nRecieved: %d", test.Expected, res)
			}
		})
	}
}

func TestParseMacrodigit(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/parse_macrodigit_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := parseMacrodigit(test.TestCase[0], 0); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if res != int(test.Expected[0][0] - '0') {
				t.Errorf("Incorrect result\nExpected: %s\nRecieved: %d", test.Expected, res)
			}
		})
	}
}

func TestExtractValue(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/extract_value_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			value := make(map[string]string)
			if err = extractValue(test.TestCase[0], &test.TestCase[1], &value, false); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if value[test.TestCase[0]] != test.Expected[0] || test.TestCase[1] != test.Expected[1] {
				t.Errorf("Incorrect result\nExpected: %s %s\nRecieved: %s %s", test.Expected[0], test.Expected[1], value[test.TestCase[0]], test.TestCase[1])
			}
		})
	}
}

func TestTerm(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/term_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := term(&test.TestCase[0]); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if res != test.Expected[0] || test.TestCase[0] != test.Expected[1] {
				t.Errorf("Incorrect result\nExpected: %s %s\nRecieved: %s %s", test.Expected[0], test.Expected[1], res, test.TestCase[0])
			}
		})
	}
}

func TestSymbol(t *testing.T) {
	tests, err := getTestsFromFile("tests/unit_tests/symbol_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			if res, err := symbol(&test.TestCase[0]); err != nil {
				if test.Expected[0] != err.Error() {
					t.Errorf("Incorrect error recognition\nExpected: %s\nRecieved: %s", test.Expected[0], err.Error())
				}
			} else if res != test.Expected[0] || test.TestCase[0] != test.Expected[1] {
				t.Errorf("Incorrect result\nExpected: %s %s\nRecieved: %s %s", test.Expected[0], test.Expected[1], res, test.TestCase[0])
			}
		})
	}
}