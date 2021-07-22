package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

// удаление пробельных символов с концов строки
func deleteSpaces(inString *string) {
	*inString = strings.Join(strings.Fields(*inString), " ")
}

// поиск точек входа и проверка на единственность
func entryPointSearch(refalProgram string) (int, int, error) {
	if strings.Count(refalProgram, "$ENTRY") != 1 {
		return -1, -1, errors.New("Exceeded entry points amount\n")
	}
	ind := strings.Index(refalProgram, "$ENTRY")
	if refalProgram[ind+6] == ' ' && refalProgram[ind+7] == 'G' && (refalProgram[ind+8] == 'o' || refalProgram[ind+8] == 'O') {
		indBeg := strings.Index(refalProgram[ind+8:], "{") + ind + 9
		indEnd := strings.Index(refalProgram[ind+8:], "}") + ind + 8
		return indBeg, indEnd, nil
	}
	return -1, -1, errors.New("Invalid entry point\n")
}

func deleteComments(refalProgram *string) {
	var (
		indBeg int
		indEnd int
	)

	for strings.Count(*refalProgram, "/*") > 0 { // удаление многострочных комментариев
		indBeg = strings.Index(*refalProgram, "/*")
		indEnd = strings.Index((*refalProgram)[indBeg+2:], "*/") + indBeg + 2
		*refalProgram = (*refalProgram)[:indBeg] + (*refalProgram)[indEnd+2:]
	}

	for strings.Count(*refalProgram, "*") > 0 { // удаление однострочных комментариев
		indBeg = strings.Index(*refalProgram, "*")
		indEnd = strings.Index((*refalProgram)[indBeg+1:], "\n") + indBeg + 1
		*refalProgram = (*refalProgram)[:indBeg] + (*refalProgram)[indEnd+1:]
	}

	*refalProgram = strings.ReplaceAll(*refalProgram, "\r", "")
}

// обработка файла
func checkFile(path string, refalProgram *string) (string, int, int, error) {
	inBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", -1, -1, err
	}
	*refalProgram = string(inBytes)

	deleteComments(refalProgram)

	indBeg, indEnd, err := entryPointSearch(*refalProgram) // поиск точек входа
	if err != nil {
		return "", -1, -1, err
	}

	entryPoint := (*refalProgram)[indBeg:indEnd]
	deleteSpaces(&entryPoint)

	return entryPoint, indBeg, indEnd, nil
}
