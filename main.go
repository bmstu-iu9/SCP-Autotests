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
func entryPointCountExceed(inString string) (bool, string, int, int) {
	s := inString
	if strings.Count(s, "$ENTRY") != 1 {
		return false, inString, -1, -1
	}
	ind := strings.Index(s, "$ENTRY")
	if inString[ind+6] == ' ' && inString[ind+7] == 'G' && (inString[ind+8] == 'o' || inString[ind+8] == 'O') {
		indBeg := strings.Index(inString[ind+8:], "{")+ind+9
		indEnd := strings.Index(inString[ind+8:], "}")+ind+8
		return true, inString[indBeg : indEnd], indBeg, indEnd
	}
	return false, inString, -1, -1
}

func cleanBytes(inBytes []byte) []byte {
	var (
		indBeg int
		indEnd int
	)

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

	inBytes = []byte(strings.ReplaceAll(string(inBytes), "  ", " "))
	return inBytes
}

// обработка файла
func checkFile(path string) ([]byte, error) {
	inBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	inBytes = cleanBytes(inBytes)
	var notExceed bool
	outString := string(inBytes)
	notExceed, outString, _, _ = entryPointCountExceed(outString) // поиск точек входа
	outString = deleteInsideSpaces(strings.TrimSpace(outString))
	if !notExceed {
		return nil, errors.New("Exceeding of entry points count\n")
	}
	return []byte(outString), nil
}
