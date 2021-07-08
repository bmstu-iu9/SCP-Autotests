package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)


func entryPointCountExceed(a []byte) (bool, []byte){
	s := string(a)
	if strings.Count(s, "$ENTRY") != 1 {
		return false, a
	}
	ind := strings.Index(s, "$ENTRY")
	if a[ind + 6] == ' ' && a[ind + 7] == 'G' && (a[ind + 8] == 'o' || a[ind + 8] == 'O') {
		return true, a[strings.Index(string(a[ind + 8 :]), "{") + ind + 9 : strings.Index(string(a[ind + 8 :]), "}") + ind + 8]
	}
	return false, a
}

func checkFile(path string) []byte {
	a, err := ioutil.ReadFile(path)
	var (
		indBeg int
		indEnd int
		indEnd1 int
		indEnd2 int
	)

	if err != nil {
		fmt.Println("Error opening file")
		log.Fatal(err)
	}

	for strings.Count(string(a), "/*") > 0 {
		indBeg = strings.Index(string(a), "/*")
		indEnd = strings.Index(string(a[indBeg+2:]), "*/") + indBeg + 2
		a = append(a[:indBeg], a[indEnd+2:]...)
	}

	for strings.Count(string(a), "*") > 0 {
		indBeg = strings.Index(string(a), "*")
		indEnd1 = strings.Index(string(a[indBeg+1:]), "\n")
		indEnd2 = strings.Index(string(a[indBeg+1:]), "*")
		if indEnd1 > indEnd2 {
			indEnd = indEnd2 + indBeg + 1
		} else {
			indEnd = indEnd1 + indBeg + 1
		}
		a = append(a[:indBeg], a[indEnd+1:]...)
	}

	for strings.Count(string(a), "  ") > 0 {
		indBeg = strings.Index(string(a), "  ")
		a = append(a[:indBeg], a[indBeg + 1:]...)
	}

	flag, res := entryPointCountExceed(a)
	if flag {
		for _, v := range res {
			fmt.Printf("%c", v)
		}
	} else {
		fmt.Println("Exceeding of entry points count")
	}
	return res
}

func main(){
	checkFile("./tests/tests_auto/test_order.ref")
}
