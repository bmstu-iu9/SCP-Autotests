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
	FilePath string `json:"FilePath"`
	TestData string `json:"TestData"`
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
	cmd := exec.Command("./scripts/create_scp.sh", fmt.Sprintf("MSCPAver%s", scpVersion))
	if err := cmd.Run(); err != nil {
		return errors.New("Error while compiling and executing a supercompiler refal-program\n")
	}
	return nil
}

func deleteSCP(scpVersion string) error {
	err := os.Remove(fmt.Sprintf("MSCPAver%s/mscp-a", scpVersion))
	if err != nil {
		return err
	}
	return nil
}

func createResidual(path string, scpVersion string) (string, error) {
	cmd := exec.Command("./scripts/create_rsd.sh", fmt.Sprintf("MSCPAver%s", scpVersion), fmt.Sprintf("../tests/%s", path), fmt.Sprintf("../tests/rsd_%s_%s", scpVersion, path))
	if err := cmd.Run(); err != nil {
		return "", errors.New("Error while compiling the refal program\n")
	}
	path = fmt.Sprintf("rsd_%s_%s", scpVersion, path)
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
	cmd := exec.Command("./scripts/run_refal.sh", refalVersion)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return errors.New("Error while compiling and executing a refal-program\n")
	}
	return nil
}

func customPrintln(s string) {
	fmt.Println(strings.TrimSpace(s))
}

func printTestResult(defaultProgramOutput, residualProgramOutput string, err1, err2 error, failCount *int) {
	if strings.Compare(defaultProgramOutput, residualProgramOutput) != 0 || err1 != nil || err2 != nil {
		fmt.Print("FAIL:\n\tDefault program output :\t")
		*failCount++
	} else {
		fmt.Print("OK:\n\tDefault program output :\t")
	}
	if err1 != nil {
		customPrintln(err1.Error())
	} else {
		customPrintln(defaultProgramOutput)
	}
	fmt.Print("\tResidual program output :\t")
	if err2 != nil {
		customPrintln(err2.Error())
	} else {
		customPrintln(residualProgramOutput)
	}
	fmt.Println()
}

func runTests(tests []MainTest, refalVersion string, SCPVersion string) error {
	failCount := 0
	if err := createSCP(SCPVersion); err != nil {
		return err
	}
	for _, test := range tests {
		path := test.FilePath
		fmt.Printf("RUN  %s\t\t TEST: %s\n", path, test.TestData)

		defaultProgramOutput, err1 := getOutput(fmt.Sprintf("tests/%s", path), test.TestData, refalVersion)

		path, _ = createResidual(path, SCPVersion)
		residualProgramOutput, err2 := getOutput(fmt.Sprintf("tests/%s", path), test.TestData, refalVersion)

		printTestResult(defaultProgramOutput, residualProgramOutput, err1, err2, &failCount)
	}
	if err := deleteSCP(SCPVersion); err != nil {
		return err
	}
	if failCount == 0 {
		fmt.Printf("--------------------------------\n\tAll tests passed!\n--------------------------------\n")
	} else {
		fmt.Println("Tests passed: ", len(tests)-failCount)
		fmt.Println("Tests failed: ", failCount)
	}
	return nil
}

func parseCommandLineFlags() (SCPVersion, refalVersion, file *string) {
	SCPVersion = flag.String("scp", "1", "supercompiler version")
	refalVersion = flag.String("v", "default", "refal version")
	file = flag.String("path", "tests/main_tests.json", "path to tests")
	flag.Parse()
	return
}

func main() {
	SCPVersion, refalVersion, file := parseCommandLineFlags()

	tests, err := getMainTests(*file)
	if err != nil {
		log.Fatal(err)
	}

	if err = runTests(tests, *refalVersion, *SCPVersion); err != nil {
		log.Fatal(err)
	}
}
