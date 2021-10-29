package cmd

import (
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
			r := serverstate.ProcessQuery("LOG hello")

			g.Assert(r).Eql("OK")
			g.Assert(len(serverstate.RawLog)).Eql(1)
			g.Assert(serverstate.RawLog[0].Query).Eql("LOG hello")
		})

		g.It(".LogsSince()", func() {
			serverstate.ProcessQuery("LOG hello")
			serverstate.RawLog[0].Timestamp = 0

			serverstate.ProcessQuery("LOG hi")
			serverstate.RawLog[1].Timestamp = 1

			g.Assert(serverstate.LogsSince(1)).Eql(
				[]RawLogLine{{1, "LOG hi", LogAttributes{Group: GROUP_DEFAULT}}})
			g.Assert(serverstate.LogsSince(0)).Eql(serverstate.RawLog)
		})

		g.It("Try invalid command", func() {
			r := serverstate.ProcessQuery("L hello")

			g.Assert(len(serverstate.RawLog)).Eql(1)
			g.Assert(r).Eql("Server didn't find command L")
		})

		g.It("Set key value", func() {
			r := serverstate.ProcessQuery("KEY_SET key=hi value=hello")

			g.Assert(r).Eql("OK")
			g.Assert(serverstate.KeyValues).Eql(map[string]string{"hi": "hello"})
		})

		g.It("Remove key value", func() {
			r := serverstate.ProcessQuery("KEY_SET key=hi value=hello")

			g.Assert(r).Eql("OK")
			g.Assert(serverstate.KeyValues).Eql(map[string]string{"hi": "hello"})

			r = serverstate.ProcessQuery("KEY_REMOVE key=hi")

			g.Assert(r).Eql("OK")
			g.Assert(serverstate.KeyValues).Eql(map[string]string{})
		})

		var keyerrors = [][2]string{
			{"KEY_SET value=hello", "key argument is necessary"},
			{"KEY_SET key=h", "value argument is necessary"},
			{"KEY_SET key=value=", "key should not be an empty string"},
			{"KEY_REMOVE", "key argument is necessary"},
			{"KEY_REMOVE key=", "key should not be an empty string"},
		}

		g.It("KEY_SET and KEY_REMOVE errors", func() {
			for i, keyerror := range keyerrors {
				r := serverstate.ProcessQuery(keyerror[0])
				g.Assert(r).Eql(keyerror[1])
				g.Assert(serverstate.RawLog[i].Attributes.Group).Eql(GROUP_INVALID)
			}

			g.Assert(len(serverstate.KeyValues)).Eql(0)
		})

		g.It("Set group", func() {
			serverstate.ProcessQuery("LOG hello group=red")
			serverstate.ProcessQuery("KEY_SET key=hello value=3 group=red")
			serverstate.ProcessQuery("KEY_REMOVE key=hello group=red")

			for _, v := range serverstate.RawLog {
				g.Assert(v.Attributes.Group).Eql("red")
			}
		})
	})

}
