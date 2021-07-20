package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type MainTest struct {
	FilePath []string  `json:"FilePath"`
	TestData []string  `json:"TestData"`
}

func getMainTests(path string) ([]MainTest, error) {
	jsonStorage, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonStorage.Close()

	data, err := ioutil.ReadAll(jsonStorage)
	if err != nil {
		return nil, err
	}

	var mainTests []MainTest

	if err = json.Unmarshal(data, &mainTests); err != nil {
		return nil, err
	}

	return mainTests, nil
}

func getOutput(path string, data string) string {
	var (
		s              string
		indBeg, indEnd int
		out			   bytes.Buffer
		cmd            *exec.Cmd
		//cmd2            *exec.Cmd
		ok 			   bool
	)

	inBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", path)
		return ""
	}
	newStr, err := getUnparametrizedEntryPoint(path, data)
	if err != nil {
		fmt.Println("Error while unparametrize")
		return ""
	}

	inBytes = cleanBytes(inBytes)
	s = string(inBytes)
	ok, _, indBeg, indEnd = entryPointCountExceed(s)
	if !ok {
		fmt.Println("Exceeding of entry points' count")
		return ""
	}

	inBytesEnd := inBytes[indEnd:]
	inBytes = append(inBytes[:indBeg], newStr...)
	inBytes = append(inBytes, inBytesEnd...)
	f, _ := os.Create("new1.ref")
	f.Write(inBytes)
	f.Close()

	cmd = exec.Command("./run_refal.sh")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Remove("new1.ref")
	return out.String()
}

func main() {
	tests, err := getMainTests("tests/main_tests.json")
	if err != nil {
		errors.New("Unable to open test files\n")
	}
	allOk := true

	for _, test := range tests {
		path := test.FilePath[0]
		sOld := getOutput(path, test.TestData[0])
		path = strings.ReplaceAll(path, "tests/", "")
		path = fmt.Sprintf("tests/RSDs/rsd_%s", path)
		sNew := getOutput(path, test.TestData[0])

		diff := strings.Compare(sOld, sNew)
		if diff != 0 {
			fmt.Printf("FAIL: %s\n\tDefault program output : %s\n\tResidual program output : %s\n", path, sOld, sNew)
			allOk = false
		} else {
			fmt.Printf("OK: %s\n\tDefault program output : %s\n\tResidual program output : %s\n", path, sOld, sNew)
			//fmt.Printf("OK\n")
		}
	}
	if allOk {
		fmt.Printf("--------------------------------\n\tAll tests passed!\n--------------------------------\n")
	}
}
