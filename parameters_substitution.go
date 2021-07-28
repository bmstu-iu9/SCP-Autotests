package main

import (
	"errors"
	"strings"
)

/*
	Function gets a parametrized entry point of a Refal_program and some data
	that should be matched to parameters. An unparametrized entry point or an error is returned.
*/
func getUnparametrizedEntryPoint(entryPoint, data string) (string, error) {
	params := strings.TrimSpace(entryPoint[:strings.Index(entryPoint, "=")])
	tokens, err := checkParameters(params)
	if err != nil {
		return "", err
	}
	values, err := mapDataToParameters(tokens, data)
	if err != nil {
		return "", err
	}
	newEntryPoint := strings.TrimSpace(entryPoint[strings.Index(entryPoint, "="):])
	return substituteParameters(values, newEntryPoint), nil
}

/*
	Function gets tokens that were created in checkParameters and data that should be mapped to the tokens.
	The data is split by parentheses and each data unit is matched to a token with function matchData in a loop.
	It returns a map (string variable -> string value), where variables are e, s, t dimension and value is a constant.
*/
func mapDataToParameters(tokens []string, data string) (map[string]string, error) {
	values := make(map[string]string)
	data = strings.ReplaceAll(data, "\\(", "(")
	data = strings.ReplaceAll(data, "\\)", ")")
	if err := matchData(tokens, &data, &values); err != nil {
		return nil, err
	}
	return values, nil
}

/*
	Function gets a token with no parentheses, data with no parentheses out of quotes and
	a pointer to map that was described in comment to mapDataToParameters function.
	The function in two loops assigns values to variables by helper matchDataHelper.
	The first loop goes from the starting point to the right till e-variable.
	The second loop - from the ending point to the left till e-variable (that is why the data string is reversed).
*/
func matchData(tokens []string, data *string, values *map[string]string) error {
	i := 0
	for ; i < len(tokens) && tokens[i][0] != 'e'; i++ {
		if err := matchDataHelper(tokens[i], data, values, false); err != nil {
			return err
		}
	}
	if len(tokens) != i {
		*data = reverse(*data)
		for j := len(tokens) - 1; j != i; j-- {
			if err := matchDataHelper(tokens[j], data, values, true); err != nil {
				return err
			}
		}
		(*values)[tokens[i]] = reverse(*data)
	} else if len(*data) != 0 {
		return errors.New("Recognition impossible\n")
	}
	return nil
}

/*
	Function is a helper to function matchData.
	If a token is enclosed with parentheses, function parse finds the substring in data
	that is also enclosed with structural parentheses. So, matchDataHelper is called recursively
	in order to work with a variable of e, s, t dimension.
	A value is set to each variable by the function extractValue.
*/
func matchDataHelper(token string, data *string, values *map[string]string, reverseInd bool) error {
	*data = strings.TrimSpace(*data)
	if len(*data) == 0 {
		return errors.New("Recognition impossible\n")
	}
	if token[0] == '(' {
		if (*data)[0] != '(' {
			return errors.New("Recognition impossible\n")
		}
		i, err := parse(*data, 1, ')')
		if err != nil {
			return err
		}
		var inParenthesesData string
		if reverseInd {
			inParenthesesData = reverse((*data)[1:i])
		} else {
			inParenthesesData = (*data)[1:i]
		}
		*data = (*data)[i+1:]
		if err = matchData(strings.Split(token[1:len(token)-1], " "), &inParenthesesData, values); err != nil {
			return err
		}
	} else if reverseInd {
		if err := extractValue(token, data, values, true); err != nil {
			return err
		}
	} else {
		if err := extractValue(token, data, values, false); err != nil {
			return err
		}
	}
	return nil
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
		} else if data[i] == '(' {
			if j, err := parse(data, i+1, ')'); err != nil {
				return -1, err
			} else {
				i = j
			}
		} else if data[i] != ' ' {
			return -1, errors.New("Invalid data\n")
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
		return -1, errors.New("Invalid data\n")
	} else {
		return i, nil
	}
}

/*
	Function is a helper to function parse. It recognises a macrodigit and
	returns	an integer where the macrodigit ends.
*/
func parseMacrodigit(data string, i int) (int, error) {
	for ; i < len(data) && data[i] >= '0' && data[i] <= '9'; i++ {
	}
	if i != len(data) && data[i] != ' ' && data[i] != ')' && data[i] != '(' && data[i] != '\'' {
		return -1, errors.New("Invalid data\n")
	} else {
		return i, nil
	}
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
	var j int
	if (*data)[0] == '(' {
		j, err = parse(*data, 1, ')')
		if err != nil {
			return
		}
		res = (*data)[:j+1]
		*data = (*data)[j+1:]
		return
	}
	res, err = symbol(data)
	return
}

/*
	Function recognises a symbol, which can be a character or a macrogigit.
	The macrodigit is found by function parseMacrodigit.
*/
func symbol(data *string) (res string, err error) {
	var i int
	if (*data)[0] == '\'' {
		if (*data)[1] != '\'' {
			res = "'" + string((*data)[1]) + "'"
			*data = "'" + (*data)[2:]
			if strings.HasPrefix(*data, "''") {
				*data = (*data)[2:]
			}
		} else {
			err = errors.New("No symbol found\n")
		}
	} else if (*data)[0] == '(' || (*data)[0] == ')' {
		err = errors.New("Invalid data\n")
	} else {
		i, err = parseMacrodigit(*data, 0)
		if err != nil {
			return
		}
		res = (*data)[:i]
		if i != len(*data) {
			if (*data)[i] == ' ' {
				*data = (*data)[i+1:]
			} else {
				*data = (*data)[i:]
			}
		} else {
			*data = ""
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
