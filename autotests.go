package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

func getOutput(path string, data string, refalVersion string) (string, error) {
	var (
		refalProgram   string
		out			   bytes.Buffer
	)

	entryPoint, indBeg, indEnd, err := checkFile(path, &refalProgram)
	if err != nil {
		return "", err
	}

	newEntryPoint, err := getUnparametrizedEntryPoint(entryPoint, data)
	if err != nil {
		return "", err
	}

	newEntryPoint = "\n = <Prout " + newEntryPoint[2:len(newEntryPoint) - 1] + ">;\n"

	refalProgram = refalProgram[:indBeg] + newEntryPoint + refalProgram[indEnd:]

	if err = createFile(refalProgram); err != nil {
		return "", err
	}

	if err = executeRefalProgram(&out, refalVersion); err != nil {
		return "", err
	}

	return out.String(), nil
}

func createFile(refalProgram string) error {
	f, err := os.Create("executable.ref")
	if err != nil {
		return err
	}

	_, err = f.WriteString(refalProgram)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func executeRefalProgram(out *bytes.Buffer, refalVersion string) error {
	cmd := exec.Command("./run_refal.sh", refalVersion)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return errors.New("Error while compiling and executing a refal-program\n")
	}
	return nil
}

func runTests(tests []MainTest, refalVersion string) {
	count := 0
	for _, test := range tests {
		path := test.FilePath[0]
		fmt.Printf("RUN  %s\n", path)

		defaultProgramOutput, err1 := getOutput(path, test.TestData[0], refalVersion)
		path = strings.ReplaceAll(path, "tests/", "")
		path = fmt.Sprintf("tests/RSDs/rsd_%s", path)
		residualProgramOutput, err2 := getOutput(path, test.TestData[0], refalVersion)

		if strings.Compare(defaultProgramOutput, residualProgramOutput) != 0 {
			fmt.Print("FAIL:\n\tDefault program output :\t")
			if err1 != nil {
				fmt.Print(err1.Error())
			} else {
				fmt.Print(defaultProgramOutput)
			}
			fmt.Print("\tResidual program output :\t")
			if err2 != nil {
				fmt.Print(err2.Error())
			} else {
				fmt.Print(residualProgramOutput)
			}
			count++
		} else {
			fmt.Printf("OK:\n\tDefault program output :\t%s\tResidual program output :\t%s", defaultProgramOutput, residualProgramOutput)
		}
	}
	if count == 0 {
		fmt.Printf("--------------------------------\n\tAll tests passed!\n--------------------------------\n")
	} else {
		fmt.Println("Tests passed: ", len(tests) - count)
		fmt.Println("Tests failed: ", count)
	}
}

func main() {

	refalVersion := flag.String("v", "default", "refal version")
	file := flag.String("path", "tests/main_tests.json", "path to tests")
	flag.Parse()

	tests, err := getMainTests(*file)
	if err != nil {
		log.Fatal(err)
	}

	runTests(tests, *refalVersion)
}
