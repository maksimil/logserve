package cmd

import (
	"fmt"
	"time"
)

const GROUP_DEFAULT = "default"
const GROUP_INVALID = "invalid_cmd"

type LogAttributes struct {
	Group string `json:"group"`
}

func InvalidAttributes() LogAttributes {
	return LogAttributes{
		Group: GROUP_INVALID,
	}
}

type RawLogLine struct {
	Timestamp  int64         `json:"timestamp"`
	Query      string        `json:"query"`
	Attributes LogAttributes `json:"attributes"`
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
		err, attrs := function(endquery, state)
		state.RawLog = append(state.RawLog,
			RawLogLine{timestamp, query, attrs})
		return err
	} else {
		state.RawLog = append(state.RawLog,
			RawLogLine{timestamp, query, InvalidAttributes()})
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

var TYPES_QUERIES map[string]func(query string, state *ServerState) (string, LogAttributes)

func init() {
	TYPES_QUERIES = map[string]func(query string, state *ServerState) (string, LogAttributes){
		"LOG": func(query string, state *ServerState) (string, LogAttributes) {
			var log, group string
			log = ""
			group = GROUP_DEFAULT
			PopulateQueryArgs(map[string]ArgAdress{
				"":      Optional(&log),
				"group": Optional(&group),
			}, query)

			return "OK", LogAttributes{Group: group}
		},

		"KEY_SET": func(query string, state *ServerState) (string, LogAttributes) {
			var key, value, group string
			group = GROUP_DEFAULT
			err := PopulateQueryArgs(map[string]ArgAdress{
				"key":   Adress(&key),
				"value": Adress(&value),
				"group": Optional(&group),
			}, query)

			if err != nil {
				return err.Error(), InvalidAttributes()
			}

			if key == "" {
				return "key should not be an empty string", InvalidAttributes()
			}

			state.KeyValues[key] = value

			return "OK", LogAttributes{Group: group}
		},

		"KEY_REMOVE": func(query string, state *ServerState) (string, LogAttributes) {
			var key, group string
			group = GROUP_DEFAULT
			err := PopulateQueryArgs(map[string]ArgAdress{
				"key":   Adress(&key),
				"group": Optional(&group),
			}, query)

			if err != nil {
				return err.Error(), InvalidAttributes()
			}

			if key == "" {
				return "key should not be an empty string", InvalidAttributes()
			}

			delete(state.KeyValues, key)

			return "OK", LogAttributes{Group: group}
		},
	}
}
