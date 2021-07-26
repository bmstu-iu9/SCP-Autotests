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
	FilePath []string `json:"FilePath"`
	TestData []string `json:"TestData"`
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

func createSCP(scpVersion string) error { //refalVersion can be added
	cmd := exec.Command("./create_scp.sh", fmt.Sprintf("MSCPAver%s", scpVersion))
	if err := cmd.Run(); err != nil {
		return errors.New("Error while compiling and executing a supercompiler refal-program\n")
	}
	return nil
}

func deleteSCP(scpVersion string) error {
	os.Remove(fmt.Sprintf("MSCPAver%s/mscp-a", scpVersion))
	return nil
}

func createResidual(path string, scpVersion string) (string, error) {
	cmd := exec.Command("./create_rsd.sh", fmt.Sprintf("MSCPAver%s", scpVersion), fmt.Sprintf("../tests/%s", path), fmt.Sprintf("../tests/rsd_%s", path))
	if err := cmd.Run(); err != nil {
		return "", errors.New("Error while compiling the refal program\n")
	}
	path = fmt.Sprintf("rsd_%s", path)
	return path, nil
}

func getOutput(path string, data string, refalVersion string) (string, error) {
	var (
		refalProgram string
		out          bytes.Buffer
	)

	entryPoint, indBeg, indEnd, err := checkFile(path, &refalProgram)
	if err != nil {
		return "", err
	}

	newEntryPoint, err := getUnparametrizedEntryPoint(entryPoint, data)
	if err != nil {
		return "", err
	}

	newEntryPoint = "\n = <Prout " + newEntryPoint[2:len(newEntryPoint)-1] + ">;\n"

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

func runTests(tests []MainTest, refalVersion string, SCPVersion string) {
	failCount := 0
	createSCP(SCPVersion)
	for _, test := range tests {
		path := test.FilePath[0]
		fmt.Printf("RUN  %s\n", path)

		defaultProgramOutput, err1 := getOutput(fmt.Sprintf("tests/%s", path), test.TestData[0], refalVersion)

		path, _ = createResidual(path, SCPVersion)
		residualProgramOutput, err2 := getOutput(fmt.Sprintf("tests/%s", path), test.TestData[0], refalVersion)

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
			fmt.Println("\n")
			failCount++
		} else {
			fmt.Printf("OK:\n\tDefault program output :\t%s\tResidual program output :\t%s\n\n", defaultProgramOutput, residualProgramOutput)
		}
	}
	deleteSCP(SCPVersion)
	if failCount == 0 {
		fmt.Printf("--------------------------------\n\tAll tests passed!\n--------------------------------\n")
	} else {
		fmt.Println("Tests passed: ", len(tests)-failCount)
		fmt.Println("Tests failed: ", failCount)
	}
}

func main() {

	SCPVersion := flag.String("scp", "1", "supercompiler version")
	refalVersion := flag.String("v", "default", "refal version")
	file := flag.String("path", "tests/main_tests.json", "path to tests")
	flag.Parse()

	tests, err := getMainTests(*file)
	if err != nil {
		log.Fatal(err)
	}

	runTests(tests, *refalVersion, *SCPVersion)
}
