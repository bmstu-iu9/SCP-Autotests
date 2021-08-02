package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func runTests(tests []Autotest) error {
	failCount := 0
	completedCount := 0

	fmt.Println("Building supercompiler...")
	if err := buildSCP(); err != nil {
		return err
	}

	fmt.Println("Starting autotests...")
	var (
		wg sync.WaitGroup
		m sync.Mutex
	)
	c := make(chan string, len(tests))

	for i := range tests {
		fmt.Printf("Test #%d is active\n", i)
		wg.Add(1)
		go test(&wg, &m, &tests, i, &failCount, &completedCount, c)
	}
	fmt.Println("All tests are active. Waiting for results...")
	wg.Wait()
	fmt.Println("\n\tAll tests completed. Press ENTER to print results...")
	if _, err := fmt.Scanln(); err != nil {
		return err
	}
	for i := 0; i < len(tests); i++ {
		val, ok := <-c
		if ok == false {
			fmt.Println("Tests are broken. Aborting...")
			break
		} else {
			fmt.Println(val)
		}
	}

	fmt.Printf("\nUsed: MSCPAver%s\n", *SCPVersion)
	fmt.Println("Tests completed:", completedCount, "out of", len(tests))
	if failCount == 0 && completedCount == len(tests) {
		fmt.Printf("--------------------------------\n\tAll tests passed!\n--------------------------------\n")
	} else {
		fmt.Println("Tests passed: ", len(tests)-failCount)
		fmt.Println("Tests failed: ", failCount)
	}

	if err := deleteAll(); err != nil {
		return err
	}

	if *saveTime {
		data, err := json.MarshalIndent(tests, " ", "\t")
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile("tests/main_tests.json", data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func buildSCP() error {
	cmd := exec.Command("./scripts/prepare_directories.sh")
	if err := cmd.Run(); err != nil {
		return errors.New("Error while creating tests/rsd and tests/errors directories\n")
	}
	cmd = exec.Command("./scripts/build_scp.sh", fmt.Sprintf("MSCPAver%s", *SCPVersion))
	if err := cmd.Run(); err != nil {
		errBody, errFile := ioutil.ReadFile("err.txt")
		if errFile != nil {
			return errors.New("Error while compiling and executing a supercompiler refal-program\n" + errFile.Error())
		}
		return errors.New("Error while compiling and executing a supercompiler refal-program\n" + string(errBody))
	}
	errBody, errFile := ioutil.ReadFile("err.txt")
	if errFile != nil {
		return errors.New("Error while compiling and executing a supercompiler refal-program\n" + errFile.Error())
	}
	fmt.Println(string(errBody))
	return nil
}

func test(wg *sync.WaitGroup, m *sync.Mutex, tests *[]Autotest, i int, failCount *int, completedCount *int, c chan<- string) {
	defer wg.Done()
	path := (*tests)[i].FilePath
	filename := path[:len(path)-4] + fmt.Sprintf("_%d", i)
	result := fmt.Sprintf("RUN_TEST_%d  %s\t DATA: %s\n", i, path, (*tests)[i].TestData)

	defaultProgramOutput, err1 := getOutput(fmt.Sprintf("tests/%s", path), filename, (*tests)[i].TestData)

	expectedTime := calcTime((*tests)[i].PerfectTime)
	path, runtime, err := createResidual(path[:len(path)-4], expectedTime, i)
	(*tests)[i].PerfectTime = runtime

	if err != nil {
		printTestResult(m, defaultProgramOutput, "", err1, err, failCount, &result)
	} else {
		filename = path[:len(path)-4]
		residualProgramOutput, err2 := getOutput(fmt.Sprintf("tests/rsd/%s", path), filename, (*tests)[i].TestData)
		printTestResult(m, defaultProgramOutput, residualProgramOutput, err1, err2, failCount, &result)
	}
	m.Lock()
	*completedCount++
	fmt.Println("Test #", i, "completed. Tests remaining:", len(*tests) - *completedCount)
	m.Unlock()
	c <- result
}

func createResidual(filename string, expectedTime time.Duration, i int) (string, time.Duration, error) {
	start := time.Now()
	path := fmt.Sprintf("rsd_%s_%s_%d.ref", *SCPVersion, filename, i)
	cmd := exec.Command("./scripts/create_rsd.sh", fmt.Sprintf("MSCPAver%s", *SCPVersion),
		fmt.Sprintf("../tests/%s.ref", filename), fmt.Sprintf("../tests/rsd/%s", path),
		strconv.FormatInt(expectedTime.Milliseconds()/1000, 10))
	if err := cmd.Run(); err != nil {
		return "", time.Duration(0), errors.New("Error while creating the residual version of the refal program\n")
	}
	end := time.Now()
	duration := end.Sub(start)
	if duration > expectedTime {
		return "", time.Duration(0), errors.New("Time limit exceeded while creating the residual version of the refal program\n")
	}
	return path, duration, nil
}

func getOutput(path, filename, data string) (string, error) {
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

	if err = createFile(refalProgram, filename); err != nil {
		return "", err
	}

	if err = executeRefalProgram(&out, filename+"_executable"); err != nil {
		return "", err
	}

	return out.String(), nil
}

func createFile(refalProgram, filename string) error {
	f, err := os.Create(filename + "_executable.ref")
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

func executeRefalProgram(out *bytes.Buffer, path string) error {
	cmd := exec.Command("./scripts/run_refal.sh", path, *refalVersion)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return errors.New("Error while compiling and executing a refal-program\n")
	}

	data, errFile := ioutil.ReadFile("tests/errors/" + path + "_err.txt")
	if errFile != nil {
		return errors.New(errFile.Error())
	}
	errBody := strings.TrimSpace(string(data))
	if len(errBody) != 0 {
		return errors.New("Error while compiling and executing a refal-program\nFor more info: tests/errors/" + path + "_err.txt\n")
	} else {
		if err := os.Remove("tests/errors/" + path + "_err.txt"); err != nil {
			return err
		}
	}

	return nil
}

func printTestResult(m *sync.Mutex, defaultProgramOutput, residualProgramOutput string, err1,
	err2 error, failCount *int, result *string) {
	if strings.Compare(defaultProgramOutput, residualProgramOutput) != 0 || err1 != nil || err2 != nil {
		*result += fmt.Sprint("FAIL:\n\tDefault program output :\t")
		m.Lock()
		*failCount++
		m.Unlock()
	} else {
		*result += fmt.Sprint("OK:\n\tDefault program output :\t")
	}
	if err1 != nil {
		*result += fmt.Sprintln(strings.TrimSpace(err1.Error()))
	} else {
		*result += fmt.Sprintln(strings.TrimSpace(defaultProgramOutput))
	}
	*result += fmt.Sprint("\tResidual program output :\t")
	if err2 != nil {
		*result += fmt.Sprintln(strings.TrimSpace(err2.Error()))
	} else {
		*result += fmt.Sprintln(strings.TrimSpace(residualProgramOutput))
	}
}

func deleteAll() error {
	err := os.Remove(fmt.Sprintf("MSCPAver%s/mscp-a", *SCPVersion))
	err = os.Remove("info.txt")
	err = os.Remove("err.txt")
	if !*rsd {
		err = os.RemoveAll("tests/rsd/")
	}
	if err != nil {
		return err
	}
	return nil
}
