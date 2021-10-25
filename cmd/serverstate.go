package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type RawLogLine struct {
	Timestamp int64  `json:"timestamp"`
	Query     string `json:"query"`
}

type ServerState struct {
	StartTime  int64             `json:"-"`
	TypeValues map[string]string `json:"keyvalues"`
	RawLog     []RawLogLine      `json:"log"`
}

func NewServerState() ServerState {
	starttime := time.Now().UnixMilli()
	return ServerState{starttime, make(map[string]string), []RawLogLine{}}
}

func (state *ServerState) ProcessQuery(query string) string {
	timestamp := time.Now().UnixMilli() - state.StartTime

	widx := 0
	for widx < len(query) && query[widx] != ' ' {
		widx += 1
	}

	ty := query[:widx]
	endquery := query[widx+1:]

	function, includes := TYPES_QUERIES[ty]

	if includes {
		state.RawLog = append(state.RawLog, RawLogLine{timestamp, query})
		return function(endquery, state)
	} else {
		return fmt.Sprintf("Server didn't find command %s", ty)
	}
}

func (state *ServerState) LogsSince(timestamp int64) []RawLogLine {
	i := len(state.RawLog) - 1
	for i >= 0 && state.RawLog[i].Timestamp >= timestamp {
		i -= 1
	}
	return state.RawLog[i+1:]
}

var TYPES_QUERIES map[string]func(query string, state *ServerState) string

func init() {
	TYPES_QUERIES = map[string]func(query string, state *ServerState) string{
		"LOG": func(query string, state *ServerState) string {
			return "OK"
		},
	}
}

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
