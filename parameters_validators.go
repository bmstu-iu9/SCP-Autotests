package main

import (
	"errors"
	"strings"
)

/*
	Function gets the left side of the entry point of a Refal-program as a string
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

/*
	Function gets the input of checkParameters function
	and returns tokens or an error that there are not enough parentheses in the string
*/
func splitToTokens(params string) ([]string, error) {
	tokens := make([]string, 0)
	var currentToken string
	for _, c := range params {
		if c == '(' {
			if len(currentToken) != 0 {
				if currentToken[0] == '(' {
					return nil, errors.New("Nested parentheses occurred\n")
				} else {
					tokens = append(tokens, currentToken)
				}
			}
			currentToken = string(c)
		} else if c == ')' {
			if len(currentToken) == 0 || currentToken[0] != '(' {
				return nil, errors.New("Not enough parentheses\n")
			}
			currentToken = strings.TrimSpace(currentToken)
			tokens = append(tokens, currentToken+string(c))
			currentToken = ""
		} else if c == ' ' {
			if len(currentToken) != 0 {
				if currentToken[0] == '(' {
					if len(currentToken) != 1 && currentToken[len(currentToken)-1] != ' ' {
						currentToken += string(c)
					}
				} else {
					tokens = append(tokens, currentToken)
					currentToken = ""
				}
			}
		} else {
			currentToken += string(c)
		}
	}
	if len(currentToken) != 0 {
		if len(currentToken) > 2 && currentToken[0] == '(' && currentToken[len(currentToken)-1] != ')' {
			return nil, errors.New("Not enough parentheses\n")
		}
		tokens = append(tokens, currentToken)
	}
	return tokens, nil
}

/*
	Function gets tokens, checks if there are nested parenthesis in a token
	and returns an error if a token is in correct, which is validated by function checkVariables
*/
func checkTokens(tokens []string) error {
	var noParenthesesTokens []string
	coincidences := make(map[string]bool)
	for _, t := range tokens {
		if strings.Index(t, "(") == -1 && strings.Index(t, ")") == -1 {
			noParenthesesTokens = append(noParenthesesTokens, t)
		} else if err := checkVariables(strings.Split(t[1:len(t)-1], " "), &coincidences); err != nil {
			return err
		}
	}
	if err := checkVariables(noParenthesesTokens, &coincidences); err != nil {
		return err
	}
	return nil
}

/*
	Function gets a set of strings that are supposed to be variables.
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
