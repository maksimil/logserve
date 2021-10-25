package cmd

import (
	"errors"
	"testing"

	goblin "github.com/franela/goblin"
)

func TestQueryParsing(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("QueryParsing", func() {
		QUERIES := []struct {
			query    string
			expected (map[string]string)
			err      error
		}{
			{
				"",
				map[string]string{},
				nil,
			},
			{
				"Raw long value",
				map[string]string{
					"": "Raw long value",
				},
				nil,
			},
			{
				"key=Hello value=Hi",
				map[string]string{
					"key":   "Hello",
					"value": "Hi",
				},
				nil,
			},
			{
				"Pre-value key=value=",
				map[string]string{
					"":      "Pre-value",
					"key":   "",
					"value": "",
				},
				nil,
			},
			{
				"key=Hi  value=Hello",
				map[string]string{
					"key":   "Hi",
					"value": "Hello",
				},
				nil,
			},
			{
				"Long pre value key=hello   secondkey=hi hello bye lastone= llastone=",
				map[string]string{
					"":          "Long pre value",
					"key":       "hello",
					"secondkey": "hi hello bye",
					"lastone":   "",
					"llastone":  "",
				},
				nil,
			},
			{
				"Errored query = value",
				nil,
				errors.New(EMPTY_KEY_ERROR),
			},
		}

		for i := range QUERIES {
			test := QUERIES[i]
			g.It(test.query, func() {
				argmap, err := ParseQueryArgs(test.query)
				g.Assert(argmap).Eql(test.expected)
				g.Assert(err).Eql(test.err)
			})
		}
	})
}
