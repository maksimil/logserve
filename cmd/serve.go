package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	port int
)

var (
	state ServerState
)

func RecordHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "" {
		json.NewEncoder(rw).Encode(state.rawlog)
	} else {
		log.Printf("/record was called using %s", r.Method)
	}
}

func RunServe(cmd *cobra.Command, args []string) {
	state = NewServerState()

	// getting the record by GET
	http.HandleFunc("/data/json", RecordHandler)

	fmt.Printf("Listening on http://localhost:%d/record\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		panic(err)
	}
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs logging web api on localhost",
	Run:   RunServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8000, "--port <PORT NUMBER>")
}

type ServerState struct {
	rawlog []string
}

func NewServerState() ServerState {
	return ServerState{[]string{"hi"}}
}
