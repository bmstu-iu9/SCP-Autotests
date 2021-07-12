package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type Test struct {
	TestCase []string `json:"TestCase"`
	Expected []string `json:"Expected"`
}

func getTestsFromFile(path string) ([]Test, error) {
	jsonStorage, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonStorage.Close()

	data, err := ioutil.ReadAll(jsonStorage)
	if err != nil {
		return nil, err
	}

	var tests []Test

	if err = json.Unmarshal(data, &tests); err != nil {
		return nil, err
	}

	return tests, nil
}

func TestCheckEntry(t *testing.T) {
	tests, err := getTestsFromFile("tests/entry_point_tests.json")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		t.Run(test.TestCase[0], func(t *testing.T) {
			res, err := checkFile(test.TestCase[0])
			if err != nil {
				t.Error(err)
			} else {
				if test.Expected[0] != strings.TrimSpace(string(res)) {
					t.Errorf("Incorrect result\nExpected: %s\nRecieved: %s", test.Expected[0], res)
				}
			}
		})
	}
}
