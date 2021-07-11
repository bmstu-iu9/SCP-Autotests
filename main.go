package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

// удаление пробельных символов с концов строки
func deleteSideSpaces(inBytes []byte) (outBytes []byte) {
	var (
		indBeg int
		indEnd int
	)

	indBeg = 0
	for inBytes[indBeg] == ' ' || inBytes[indBeg] == '\t' || inBytes[indBeg] == '\n' || inBytes[indBeg] == '\r' {
		indBeg += 1
	}

	indEnd = len(inBytes) - 1
	for inBytes[indEnd] == ' ' || inBytes[indEnd] == '\t' || inBytes[indEnd] == '\n' || inBytes[indBeg] == '\r' {
		indEnd -= 1
	}
	outBytes = inBytes[indBeg : indEnd]

	return outBytes
}

// поиск точек входа и проверка на единственность
func entryPointCountExceed(inBytes []byte) (bool, []byte){
	s := string(inBytes)
	if strings.Count(s, "$ENTRY") != 1 {
		return false, inBytes
	}
	ind := strings.Index(s, "$ENTRY")
	if inBytes[ind + 6] == ' ' && inBytes[ind + 7] == 'G' && (inBytes[ind + 8] == 'o' || inBytes[ind + 8] == 'O') {
		return true, inBytes[strings.Index(string(inBytes[ind + 8 :]), "{") + ind + 9 : strings.Index(string(inBytes[ind + 8 :]), "}") + ind + 9]
	}
	return false, inBytes
}

// обработка файла
func checkFile(path string) ([]byte, error) {
	inBytes, err := ioutil.ReadFile(path)
	var (
		indBeg int
		indEnd int
	)

	if err != nil {
		return nil, err
	}

	for strings.Count(string(inBytes), "/*") > 0 {	// удаление многострочных комментариев
		indBeg = strings.Index(string(inBytes), "/*")
		indEnd = strings.Index(string(inBytes[indBeg+2:]), "*/") + indBeg + 2
		inBytes = append(inBytes[:indBeg], inBytes[indEnd+2:]...)
	}

	for strings.Count(string(inBytes), "*") > 0 {		// удаление однострочных комментариев
		indBeg = strings.Index(string(inBytes), "*")
		indEnd = strings.Index(string(inBytes[indBeg+1:]), "\n") + indBeg + 1
		inBytes = append(inBytes[:indBeg], inBytes[indEnd+1:]...)
	}

	for strings.Count(string(inBytes), "  ") > 0 {	// удаление двойных пробелов
		indBeg = strings.Index(string(inBytes), "  ")
		inBytes = append(inBytes[:indBeg], inBytes[indBeg + 1:]...)
	}

	notExceed, outBytes := entryPointCountExceed(inBytes)	// поиск точек входа
	outBytes = deleteSideSpaces(outBytes)
	if !notExceed {
		return nil, errors.New("Exceeding of entry points count\n")
	}
	return outBytes, nil
}

func main(){
	checkFile("./tests/test_....ref")	// подставьте вместо многоточия нужный файл
}
