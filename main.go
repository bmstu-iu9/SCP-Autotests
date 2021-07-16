package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

// удаление пробельных символов с концов строки
func deleteInsideSpaces(inString string) string {
	inString = strings.ReplaceAll(inString, "\t", "")
	inString = strings.ReplaceAll(inString, "\n", "")
	inString = strings.ReplaceAll(inString, "\r", "")
	return inString
}

// поиск точек входа и проверка на единственность
func entryPointCountExceed(inString string) (bool, string) {
	s := inString
	if strings.Count(s, "$ENTRY") != 1 {
		return false, inString
	}
	ind := strings.Index(s, "$ENTRY")
	if inString[ind+6] == ' ' && inString[ind+7] == 'G' && (inString[ind+8] == 'o' || inString[ind+8] == 'O') {
		return true, inString[strings.Index(string(inString[ind+8:]), "{")+ind+9 : strings.Index(string(inString[ind+8:]), "}")+ind+8]
	}
	return false, inString
}

// обработка файла
func checkFile(path string) ([]byte, error) {
	inBytes, err := ioutil.ReadFile(path)
	var (
		indBeg    int
		indEnd    int
		outString string
	)

	if err != nil {
		return nil, err
	}

	for strings.Count(string(inBytes), "/*") > 0 { // удаление многострочных комментариев
		indBeg = strings.Index(string(inBytes), "/*")
		indEnd = strings.Index(string(inBytes[indBeg+2:]), "*/") + indBeg + 2
		inBytes = append(inBytes[:indBeg], inBytes[indEnd+2:]...)
	}

	for strings.Count(string(inBytes), "*") > 0 { // удаление однострочных комментариев
		indBeg = strings.Index(string(inBytes), "*")
		indEnd = strings.Index(string(inBytes[indBeg+1:]), "\n") + indBeg + 1
		inBytes = append(inBytes[:indBeg], inBytes[indEnd+1:]...)
	}

	var notExceed bool
	outString = strings.ReplaceAll(string(inBytes), "  ", " ")
	notExceed, outString = entryPointCountExceed(outString) // поиск точек входа
	outString = deleteInsideSpaces(strings.TrimSpace(outString))
	if !notExceed {
		return nil, errors.New("Exceeding of entry points count\n")
	}
	return []byte(outString), nil
}
