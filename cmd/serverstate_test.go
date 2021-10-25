package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	goblin "github.com/franela/goblin"
)

func TestServerState(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("ServerState", func() {
		var serverstate ServerState
		g.BeforeEach(func() {
			serverstate = NewServerState()
		})
		g.It(".ProcessQuery on LOG  quieries", func() {
			serverstate.ProcessQuery("LOG hello")

			g.Assert(len(serverstate.RawLog)).Eql(1)
			g.Assert(serverstate.RawLog[0].Query).Eql("LOG hello")
		})

		g.It("convert to json", func() {
			serverstate.ProcessQuery("LOG hello")
			serverstate.RawLog[0].Timestamp = 0

			var buff bytes.Buffer
			json.NewEncoder(&buff).Encode(serverstate)

			data := buff.String()

			g.Assert(data).Eql("{\"keyvalues\":{},\"log\":[{\"timestamp\":0,\"query\":\"LOG hello\"}]}\n")
		})

		g.It(".LogsSince()", func() {
			serverstate.ProcessQuery("LOG hello")
			serverstate.RawLog[0].Timestamp = 0

			serverstate.ProcessQuery("LOG hi")
			serverstate.RawLog[1].Timestamp = 1

			g.Assert(serverstate.LogsSince(1)).Eql([]RawLogLine{{1, "LOG hi"}})
			g.Assert(serverstate.LogsSince(0)).Eql(serverstate.RawLog)
		})
	})

}

func TestQueryParsing(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("QueryParsing", func() {
		QUERIES := []struct {
			query    string
			expected map[string]string
		}{
			{
				"key=Hello value=Hi",
				map[string]string{
					"key":   "Hello",
					"value": "Hi",
				},
			},
			{
				"Pre-value key=value=",
				map[string]string{
					"":      "Pre-value",
					"key":   "",
					"value": "",
				},
			},
			{
				"key=Hi  value=Hello",
				map[string]string{
					"key":   "Hi",
					"value": "Hello",
				},
			},
		}

		g.It("ParseQuieryArgs", func() {
			for _, test := range QUERIES {
				g.Assert(ParseQueryArgs(test.query)).Eql(test.expected)
			}
		})
	})
}
