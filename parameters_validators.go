package main

import (
	"errors"
	"strings"
)

/* 	Function gets the left side of the entry point of a Refal-program as a string
	and returns so called tokens or an error if at least one token is incorrect.
	Tokens are substrings, which consist of variables e, s, t dimension or an parenthetical set of such variables
 */
func checkParameters(params string) ([]string, error) {
	tokens, err := splitToTokens(params)
	if err != nil {
		return nil, err
	}
	if err = checkTokens(tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

/* 	Function gets the input of checkParameters function
	and returns tokens or an error that there are not enough parentheses in the string
 */
func splitToTokens(params string) ([]string, error) {
	tokens := strings.Split(params, " ")
	for i := 0; i < len(tokens); i++ {
		if len(tokens[i]) == 0 {
			tokens = append(tokens[ : i], tokens[i + 1 :]...)
			i--
		} else if tokens[i][0] == '(' && tokens[i][len(tokens[i]) - 1] != ')' {
			j := i + 1
			for ;j < len(tokens) && tokens[j][len(tokens[j]) - 1] != ')'; j++ {	}
			if j == len(tokens) {
				return nil, errors.New("Not enough parentheses\n")
			}
			tokens[i] = strings.Join(tokens[i : j + 1], " ")
			tokens = append(tokens[ : i + 1], tokens[j + 1 : ]...)
		} else if tokens[i][0] != '(' {
			j := i + 1
			for ;j < len(tokens) && tokens[j][0] != '('; j++ {	}
			tokens[i] = strings.Join(tokens[i : j], " ")
			tokens = append(tokens[ : i + 1], tokens[j : ]...)
		}
	}

	return tokens, nil
}

/* 	Function gets tokens, checks if there are nested parenthesis in a token
	and returns an error if a token is in correct, which is validated by function checkVariables
*/
func checkTokens(tokens []string) error {
	var noParenthesesTokens []string
	coincidences := make(map[string]bool)
	for _, t := range tokens {
		if strings.Index(t, "(") == -1 && strings.Index(t, ")") == -1 {
			noParenthesesTokens = append(noParenthesesTokens, strings.Split(t, " ")...)
		} else if strings.LastIndex(t, "(") != 0 || strings.Index(t, ")") != len(t) - 1 {
			return errors.New("Nested parentheses occurred\n")
		} else if err := checkVariables(strings.Split(t[ 1: len(t) - 1], " "), &coincidences); err != nil {
			return err
		}
	}
	if err := checkVariables(noParenthesesTokens, &coincidences); err != nil {
		return err
	}
	return nil
}

/*	Function gets a set of strings that are supposed to be variables.
	This set is generated in function checkTokens by splitting a token with a single space.
	The set is checked whether there are any coincidences in t- and s-variables and more than one e-variable.
	An error is returned.
 */
func checkVariables(vars []string, coincidences *map[string]bool) error {
	for _, v := range vars {
		if strings.HasPrefix(v, "e.") {
			if (*coincidences)["e"] {
				return errors.New("More than one e-variable occurred\n")
			} else {
				(*coincidences)["e"] = true
			}
		} else if !strings.HasPrefix(v, "s.") && !strings.HasPrefix(v, "t.") {
			return errors.New("Incorrect variable type\n")
		}
		if (*coincidences)[v] {
			return errors.New("Coincidences in variables occurred\n")
		} else {
			(*coincidences)[v] = true
		}
	}
	(*coincidences)["e"] = false
	return nil
}
