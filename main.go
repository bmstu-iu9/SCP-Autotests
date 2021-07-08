package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func checkEntry(a []byte, i int) bool {
	if a[i] == '$' && a[i + 1] == 'E' && a[i + 2] == 'N' && a[i + 3] == 'T' && a[i + 4] == 'R' && a[i + 5] == 'Y' && a[i + 6] == ' '  && a[i + 7] == 'G' && (a[i + 8] == 'O' || a[i + 8] == 'o') {
		return true
	} else {
		return false
	}
}

func checkFile(path string) []byte {
	a, err := ioutil.ReadFile(path)
	noSpaces := a[:]
	noComm := a[:]
	cnt := 0
	i := 0
	var iBeg int
	var iEnd int

	if err != nil {
		fmt.Println("Error opening file")
		log.Fatal(err)
	}

	spaceBefore := false
	for _, v := range a {
		if v == ' ' || v == '\t'{
			if !spaceBefore && v != '\t' {
				noSpaces[cnt] = v
				cnt += 1
			}
			spaceBefore = true
		} else {
			noSpaces[cnt] = v
			cnt += 1
			spaceBefore = false
		}
	}
	noSpaces = noSpaces[ : cnt]

	cnt = 0
	for i < len(noSpaces) {
		if noSpaces[i] == '/' && noSpaces[i + 1] == '*' {
			i += 2
			for !(noSpaces[i] == '*' && noSpaces[i + 1] == '/') {
				i += 1
			}
			i += 1
		} else if noSpaces[i] == '*' {
			i += 1
			for !(noSpaces[i] == '*' || noSpaces[i] == '\n'){
				i += 1
			}
		} else {
			noComm[cnt] = noSpaces[i]
			cnt += 1
		}
		i += 1
	}
	noComm = noComm[ : cnt]

	for i, _ := range noComm {
		if checkEntry(noComm, i) {
			iBeg = i + 11
			iEnd = i + 11
			for noComm[iEnd] != '}' {
				iEnd += 1
			}
		}
	}
	res := noComm[iBeg + 1 : iEnd]
	for _, v := range res {
		fmt.Printf("%c", v)
	}
	return res
}

func main(){
	checkFile("./test.txt")
}
