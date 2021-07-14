package main

import (
	"errors"
	"strings"
)

/*
	Function gets a path to a Refal_program with a parametrized entry point and some data
	that should be matched to parameters. An unparametrized entry point or an error is returned.
*/
func getUnparametrizedEntryPoint(path, data string) ([]byte, error) {
	entryPoint, err := checkFile(path)
	if err != nil {
		return nil, err
	}
	params := strings.TrimSpace(string(entryPoint)[:strings.Index(string(entryPoint), "=")])
	tokens, err := checkParameters(params)
	if err != nil {
		return nil, err
	}
	values, err := mapDataToParameters(tokens, data)
	if err != nil {
		return nil, err
	}
	newEntryPoint := strings.TrimSpace(string(entryPoint)[strings.Index(string(entryPoint), "="):])
	return []byte(substituteParameters(values, newEntryPoint)), nil
}

/*
	Function gets tokens that were created in checkParameters and data that should be mapped to the tokens.
	The data is split by parentheses and each data unit is matched to a token with function matchData in a loop.
	It returns a map (string variable -> string value), where variables are e, s, t dimension and value is a constant.
*/
func mapDataToParameters(tokens []string, data string) (map[string]string, error) {
	values := make(map[string]string)
	splitData, err := splitData(data)
	if err != nil {
		return nil, err
	}
	if len(tokens) != len(splitData) {
		return nil, errors.New("Mapping is not possible\n")
	} else {
		for i, t := range tokens {
			if strings.HasPrefix(t, "(") {
				if !strings.HasPrefix(splitData[i], "(") {
					return nil, errors.New("Mapping is not possible\n")
				} else if err = matchData(t[1:len(t)-1], splitData[i][1:len(splitData[i])-1], &values); err != nil {
					return nil, err
				}
			} else if err = matchData(t, splitData[i], &values); err != nil {
				return nil, err
			}
		}
	}
	return values, nil
}

/*
	Function gets the data that should substitute the parameters in entry point of a Refal-program and
	split it by parentheses. It returns split data.
*/
func splitData(data string) ([]string, error) {
	splitData := make([]string, 0)
	if len(data) == 0 {
		return []string{""}, nil
	}
	for i := 0; i < len(data); i++ {
		if data[i] == '(' {
			j, err := parse(data, i+1, ')')
			if err != nil {
				return nil, err
			}
			splitData = append(splitData, strings.TrimSpace(data[i:j+1]))
			i = j
		} else if data[i] != ' ' {
			j, err := parse(data, i, '(')
			if err != nil {
				return nil, err
			}
			splitData = append(splitData, strings.TrimSpace(data[i:j]))
			i = j - 1
		}
	}
	return splitData, nil
}

/*
	Function gets a string, an integer from which character in order start to parse the string and
	one of two characters: ')' or '('. It recognises a string enclosed with single quotes or a macrodigit and
	returns an integer where a string or a set of macrodigits ends.
*/
func parse(data string, i int, c byte) (int, error) {
	for ; i < len(data) && data[i] != c; i++ {
		if data[i] == '\'' {
			if j, err := parseString(data, i+1); err != nil {
				return -1, err
			} else {
				i = j
			}
		} else if data[i] >= '0' && data[i] <= '9' {
			if j, err := parseMacrodigit(data, i+1); err != nil {
				return -1, err
			} else {
				i = j - 1
			}
		} else if data[i] != ' ' {
			return -1, errors.New("Only strings enclosed with single quotes are supported\n")
		}
	}
	if c == ')' && i == len(data) {
		return -1, errors.New("Not enough parentheses\n")
	}
	return i, nil
}

/*
	Function is a helper to function parse. It recognises a string enclosed by single quotes and
	returns	an integer where the string ends.
*/
func parseString(data string, i int) (int, error) {
	for ; i < len(data) && data[i] != '\''; i++ {
		if data[i] == '\\' {
			i++
		}
	}
	if i == len(data) {
		return -1, errors.New("Incorrect string definition\n")
	} else {
		return i, nil
	}
}

/*
	Function is a helper to function parse. It recognises a macrodigit and
	returns	an integer where the macrodigit ends.
*/
func parseMacrodigit(data string, i int) (int, error) {
	for ; i < len(data) && (data[i] >= '0' && data[i] <= '9') && data[i] != ' ' && data[i] != ')' && data[i] != '\''; i++ {
	}
	if i != len(data) && data[i] != ' ' && data[i] != ')' && data[i] != '\'' {
		return -1, errors.New("Incorrect macrodigit definition\n")
	} else {
		return i, nil
	}
}

/*
	Function gets a token with no parentheses, data with no parentheses out of quotes and
	a pointer to map that was described in comment to mapDataToParameters function.
	A token is split by a space in order to get variables of e, s, t dimension.
	A value is set to each variable in two loops by the function extractValue.
	The first loop goes from the starting point to the left till e-variable.
	The second loop - from the ending point to the right till e-variable (that is why the data string is reversed).
*/
func matchData(token string, data string, values *map[string]string) error {
	vars := strings.Split(token, " ")
	data = strings.ReplaceAll(data, "\\(", "(")
	data = strings.ReplaceAll(data, "\\)", ")")
	i := 0
	for ; i < len(vars) && vars[i][0] != 'e'; i++ {
		if err := extractValue(vars[i], &data, values, false); err != nil {
			return err
		}
	}
	data = reverse(data)
	for j := len(vars) - 1; j != i; j-- {
		if err := extractValue(vars[j], &data, values, true); err != nil {
			return err
		}
	}
	(*values)[vars[i]] = reverse(data)
	return nil
}

/*
	Function gets the type of variable, data and a map of values to variables.
	It recognises a term and symbol. The reverseInd shows whether the data was reversed.
*/
func extractValue(variable string, data *string, values *map[string]string, reverseInd bool) error {
	var (
		val string
		err error
	)
	*data = strings.TrimSpace(*data)
	if variable[0] == 't' {
		val, err = term(data)
	} else if variable[0] == 's' {
		val, err = symbol(data)
	}
	if err != nil {
		return err
	} else {
		if reverseInd {
			(*values)[variable] = reverse(val)
		} else {
			(*values)[variable] = val
		}
	}
	return nil
}

/*
	Function recognises a term, which can be an expression enclosed by parentheses or a symbol.
*/
func term(data *string) (res string, err error) {
	if strings.HasPrefix(*data, "'(") {
		i := 2
		for ; i < len(*data) && (*data)[i] != '\'' && (*data)[i] != ')'; i++ {
		}
		if i != len(*data) && (*data)[i] != '\'' {
			res = (*data)[:i+1] + "'"
			*data = (*data)[i+1:]
			return
		}
	}
	res, err = symbol(data)
	return
}

/*
	Function recognises a symbol, which can be a character or a macrogigit.
	The macrodigit is found by function parseMacrodigit.
*/
func symbol(data *string) (res string, err error) {
	if (*data)[0] == '\'' {
		if (*data)[1] != '\'' {
			res = "'" + string((*data)[1]) + "'"
			*data = "'" + (*data)[2:]
		} else {
			err = errors.New("No symbol found\n")
		}
	} else {
		i, _ := parseMacrodigit(*data, 0)
		res = (*data)[:i]
		if (*data)[i] == ' ' {
			*data = (*data)[i+1:]
		} else {
			*data = (*data)[i:]
		}
	}
	return
}

/*
	Function reverses a string. Parentheses are converted into
	each other to save the order of them.
*/
func reverse(str string) (result string) {
	for _, v := range str {
		if v == ')' {
			v = '('
		} else if v == '(' {
			v = ')'
		}
		result = string(v) + result
	}
	return
}

/*
	Function replaces the variables in the entry point of the Refal-program with the values
	that wee assigned to them in function MapDataToParameters.
*/
func substituteParameters(values map[string]string, entryPoint string) string {
	for key, value := range values {
		entryPoint = strings.ReplaceAll(entryPoint, key, value)
	}
	return entryPoint
}
