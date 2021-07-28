package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

/*
	Function gets a string
	and replaces all single or multiple white-space characters by only one space character
*/
func deleteSpaces(inString *string) {
	*inString = strings.Join(strings.Fields(*inString), " ")
}

/*
	Function gets a Refal-program as a string,
	finds the entry point in it, checks whether Go function definition is correct
	and returns bounds of the entry point in the Refal-program and an error, if something went wrong
*/
func entryPointSearch(refalProgram string) (int, int, error) {
	if strings.Count(refalProgram, "$ENTRY") != 1 {
		return -1, -1, errors.New("Exceeded entry points amount\n")
	}
	ind := strings.Index(refalProgram, "$ENTRY")
	if refalProgram[ind+6] == ' ' && refalProgram[ind+7] == 'G' &&
		(refalProgram[ind+8] == 'o' || refalProgram[ind+8] == 'O') {
		indBeg := strings.Index(refalProgram[ind+8:], "{") + ind + 9
		indEnd := strings.Index(refalProgram[ind+8:], "}") + ind + 8
		return indBeg, indEnd, nil
	}
	return -1, -1, errors.New("Invalid entry point\n")
}

/*
	Function gets a Refal-program as a string
	and deletes all multi-line and single-line comments in it
*/
func deleteComments(refalProgram *string) {
	var (
		indBeg int
		indEnd int
	)

	for strings.Count(*refalProgram, "/*") > 0 { // multi-line comment deletion
		indBeg = strings.Index(*refalProgram, "/*")
		indEnd = strings.Index((*refalProgram)[indBeg+2:], "*/") + indBeg + 2
		*refalProgram = (*refalProgram)[:indBeg] + (*refalProgram)[indEnd+2:]
	}

	for strings.Count(*refalProgram, "\n*") > 0 { // single-line comment deletion
		indBeg = strings.Index(*refalProgram, "\n*")
		indEnd = strings.Index((*refalProgram)[indBeg+2:], "\n") + indBeg + 2
		*refalProgram = (*refalProgram)[:indBeg+1] + (*refalProgram)[indEnd+1:]
	}

	*refalProgram = strings.ReplaceAll(*refalProgram, "\r", "")
}

/*
	Function gets a path to a file with a Refal-program and a pointer to a string where
	the Refal-program turns out to be after the function execution. It deletes comments in the Refal-program
	and returns entry point as a string and bounds where entry point is located in the Refal-program
*/
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
