package cmd

import (
	"errors"
	"strings"
)

const (
	EMPTY_KEY_ERROR = "\"\" is not a valid key"
)

func ParseQueryArgs(query string) (map[string]string, error) {
	query = strings.Trim(query, " ")

	// check if the query is not empty
	if query == "" {
		return map[string]string{}, nil
	}

	// find all the key statements
	keyranges := [][2]int{}
	for i := 0; i < len(query); i++ {
		if query[i] == '=' {
			j := i - 1
			for j >= 0 && query[j] != ' ' && query[j] != '=' {
				j -= 1
			}
			if j == i-1 {
				return nil, errors.New(EMPTY_KEY_ERROR)
			}
			keyranges = append(keyranges, [2]int{j + 1, i})
		}
	}

	keyranges = append(keyranges, [2]int{len(query), -1})

	argmap := map[string]string{}

	// pre-named args
	nokey := strings.Trim(query[:keyranges[0][0]], " ")
	if nokey != "" {
		argmap[""] = nokey
	}

	// to last named arg
	for i := 0; i < len(keyranges)-1; i++ {
		ks := keyranges[i][0]
		ke := keyranges[i][1]
		k := keyranges[i+1][0]

		key := strings.Trim(query[ks:ke], " ")

		if key == "" {
			return nil, errors.New(EMPTY_KEY_ERROR)
		}

		value := strings.Trim(query[ke+1:k], " ")
		argmap[key] = value
	}

	return argmap, nil
}
