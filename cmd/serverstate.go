package cmd

import (
	"time"
)

type RawLogLine struct {
	Timestamp int64  `json:"timestamp"`
	Query     string `json:"query"`
}

type ServerState struct {
	StartTime  int64             `json:"-"`
	TypeValues map[string]string `json:"typevalues"`
	RawLog     []RawLogLine      `json:"log"`
}

func NewServerState() ServerState {
	starttime := time.Now().UnixMilli()
	return ServerState{starttime, make(map[string]string), []RawLogLine{}}
}

func (state *ServerState) ProcessQuery(query string) string {
	timestamp := time.Now().UnixMilli() - state.StartTime
	state.RawLog = append(state.RawLog, RawLogLine{timestamp, query})
	return "OK"
}
