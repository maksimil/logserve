package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		json.NewEncoder(rw).Encode(state)
	} else {
		log.Printf("/data/json was called using %s", r.Method)
	}
}

func LogHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rqbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		query := string(rqbody)
		log.Println(query)
		response := state.ProcessQuery(query)
		fmt.Fprint(rw, response)
	} else {
		log.Printf("/log was called using %s", r.Method)
	}
}

func RunServe(cmd *cobra.Command, args []string) {
	state = NewServerState()

	// getting the data in json form by GET /data/json
	http.HandleFunc("/data/json", RecordHandler)

	// adding log information
	http.HandleFunc("/log", LogHandler)

	fmt.Printf("Listening on http://localhost:%d\n", port)
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
