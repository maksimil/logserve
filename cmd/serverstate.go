package cmd

import (
	"fmt"
	"time"
)

type RawLogLine struct {
	Timestamp int64  `json:"timestamp"`
	Query     string `json:"query"`
}

type ServerState struct {
	StartTime int64             `json:"-"`
	KeyValues map[string]string `json:"keyvalues"`
	RawLog    []RawLogLine      `json:"log"`
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
	endquery := query[widx:]

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

		"KEY_SET": func(query string, state *ServerState) string {
			argmap, err := ParseQueryArgs(query)

			if err != nil {
				return err.Error()
			}

			key, haskey := argmap["key"]
			value, hasvalue := argmap["value"]

			if !haskey {
				return "key argument is necessary"
			}

			if !hasvalue {
				return "value argument is necessary"
			}

			if key == "" {
				return "key should not be an empty string"
			}

			state.KeyValues[key] = value
			return "OK"
		},

		"KEY_REMOVE": func(query string, state *ServerState) string {
			argmap, err := ParseQueryArgs(query)

			if err != nil {
				return err.Error()
			}

			key, haskey := argmap["key"]
			if !haskey {
				return "key argument is necessary"
			}

			if key == "" {
				return "key should not be an empty string"
			}

			delete(state.KeyValues, key)

			return "OK"
		},
	}
}
