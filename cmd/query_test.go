package cmd

import (
	"errors"
	"fmt"
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

func TestQueryPopulation(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("QueryPopulating", func() {
		g.It("Populates values", func() {
			var key1, value, thing, thing2 string
			err := PopulateQueryArgs(map[string]ArgAdress{"key1": Adress(&key1), "value": Adress(&value), "thing": Adress(&thing), "thing2": Adress(&thing2)}, "key1=hi value=hello thing=thing2=")

			g.Assert(err).Eql(nil)
			g.Assert(key1).Eql("hi")
			g.Assert(value).Eql("hello")
			g.Assert(thing).Eql("")
			g.Assert(thing2).Eql("")
		})

		g.It("Errors out on not having an arg", func() {
			var key string
			err := PopulateQueryArgs(map[string]ArgAdress{"key": Adress(&key)}, "value=")

			g.Assert(err).Eql(fmt.Errorf(NECESSARY_ARG, "key"))
		})

		g.It("Errors out on having an unnecessary arg", func() {
			var key string
			err := PopulateQueryArgs(map[string]ArgAdress{"key": Adress(&key)}, "value=key=v")

			g.Assert(err).Eql(fmt.Errorf(UNNECESSARY_ARG, "value"))
		})

		g.It("Does not require an optional arg", func() {
			var key, optional string
			optional = "value"
			err := PopulateQueryArgs(map[string]ArgAdress{"key": Adress(&key), "optional": Optional(&optional)}, "key=va")

			g.Assert(err).Eql(nil)
			g.Assert(key).Eql("va")
			g.Assert(optional).Eql("value")

			optional = "value"
			err = PopulateQueryArgs(map[string]ArgAdress{"key": Adress(&key), "optional": Optional(&optional)}, "key=va optional=vl")

			g.Assert(err).Eql(nil)
			g.Assert(key).Eql("va")
			g.Assert(optional).Eql("vl")
		})
	})
}
