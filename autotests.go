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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MainTest struct {
	FilePath    string        `json:"FilePath"`
	TestData    string        `json:"TestData"`
	PerfectTime time.Duration `json:"PerfectTime"`
}

var (
	SCPVersion   *string
	refalVersion *string
	file         *string
	rsd          *string
	markTime     *string
)

func getMainTests() ([]MainTest, error) {
	jsonStorage, err := os.Open(*file)
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

func createSCP() error { //refalVersion can be added
	cmd := exec.Command("./scripts/create_scp.sh", fmt.Sprintf("MSCPAver%s", *SCPVersion))
	if err := cmd.Run(); err != nil {
		return errors.New("Error while compiling and executing a supercompiler refal-program\n")
	}
	return nil
}

func deleteAll() error {
	err := os.Remove(fmt.Sprintf("MSCPAver%s/mscp-a", *SCPVersion))
	err = os.Remove("info.txt")
	if *rsd == "no" {
		err = os.RemoveAll("tests/rsd/")
	}
	if err != nil {
		return err
	}
	return nil
}

func createResidual(filename string, expectedTime time.Duration, i int) (string, time.Duration, error) {
	start := time.Now()
	path := fmt.Sprintf("rsd_%s_%s_%d.ref", *SCPVersion, filename, i)
	cmd := exec.Command("./scripts/create_rsd.sh", fmt.Sprintf("MSCPAver%s", *SCPVersion),
		fmt.Sprintf("../tests/%s.ref", filename), fmt.Sprintf("../tests/rsd/%s", path),
		strconv.FormatInt(expectedTime.Milliseconds()/1000, 10))
	if err := cmd.Run(); err != nil {
		return "", time.Duration(0), errors.New("Error while compiling the refal program\n")
	}
	end := time.Now()
	duration := end.Sub(start)
	return path, duration, nil
}

func calcTime(perfectTime time.Duration) time.Duration {
	var timeOut time.Duration
	if perfectTime.Minutes() > 5 {
		timeOut = time.Duration(perfectTime.Minutes() * 1.2)
	} else if perfectTime.Minutes() > 1 {
		timeOut = time.Duration(perfectTime.Minutes() * 1.5)
	} else {
		timeOut = perfectTime*2 + (2 * time.Second)
	}
	return timeOut
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
	return nil
}

func printTestResult(defaultProgramOutput, residualProgramOutput string, err1,
	err2 error, failCount *int32, result *string) {
	if strings.Compare(defaultProgramOutput, residualProgramOutput) != 0 || err1 != nil || err2 != nil {
		*result += fmt.Sprint("FAIL:\n\tDefault program output :\t")
		atomic.AddInt32(failCount, 1)
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

func test(wg *sync.WaitGroup, tests *[]MainTest, i int, test MainTest, failCount *int32) {
	defer wg.Done()
	var expectedTime time.Duration

	path := test.FilePath
	filename := path[:len(path)-4] + fmt.Sprintf("_%d", i)
	result := fmt.Sprintf("RUN  %s\t\t TEST: %s\n", path, test.TestData)

	defaultProgramOutput, err1 := getOutput(fmt.Sprintf("tests/%s", path), filename, test.TestData)

	if test.PerfectTime != 0 {
		expectedTime = calcTime(test.PerfectTime)
	} else {
		expectedTime = time.Minute * 15
	}
	path, perfectTime, err := createResidual(path[:len(path)-4], expectedTime, i)
	(*tests)[i].PerfectTime = perfectTime

	if err != nil {
		printTestResult(defaultProgramOutput, "", err1, err, failCount, &result)
	} else {
		filename = path[:len(path)-4]
		residualProgramOutput, err2 := getOutput(fmt.Sprintf("tests/rsd/%s", path), filename, test.TestData)
		printTestResult(defaultProgramOutput, residualProgramOutput, err1, err2, failCount, &result)
	}
	fmt.Println(result)
}

func runTests(tests []MainTest) error {
	failCount := int32(0)

	fmt.Println("Building supercompiler...")
	if err := createSCP(); err != nil {
		return err
	}

	fmt.Println("Starting autotests...")
	var wg sync.WaitGroup
	for i := range tests {
		wg.Add(1)
		go test(&wg, &tests, i, tests[i], &failCount)
	}
	wg.Wait()

	data, _ := json.Marshal(tests)
	_ = ioutil.WriteFile("tests/main_tests.json", data, 0644)

	if failCount == 0 {
		fmt.Printf("--------------------------------\n\tAll tests passed!\n--------------------------------\n")
	} else {
		fmt.Printf("MSCPAver%s:\n", *SCPVersion)
		fmt.Println("Tests passed: ", len(tests)-int(failCount))
		fmt.Println("Tests failed: ", failCount)
	}

	if err := deleteAll(); err != nil {
		return err
	}

	return nil
}

func parseCommandLineFlags() {
	SCPVersion = flag.String("scp", "1", "supercompiler version")
	refalVersion = flag.String("v", "default", "refal version")
	file = flag.String("path", "tests/main_tests.json", "path to tests")
	rsd = flag.String("rsd", "no", "delete rsd files after")
	markTime = flag.String("time", "no", "mark the time")
	flag.Parse()
}

func main() {
	parseCommandLineFlags()

	tests, err := getMainTests()
	if err != nil {
		log.Fatal(err)
	}

	if err = runTests(tests); err != nil {
		log.Fatal(err)
	}
}
