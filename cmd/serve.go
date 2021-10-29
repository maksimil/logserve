package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	port int
)

var (
	state ServerState
)

const ERR_METHOD = "%s was called using %s"

//go:embed _gen/build.html
var WEBPAGEHTML string

func WebPageHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rw.Write([]byte(WEBPAGEHTML))
	} else {
		log.Printf(ERR_METHOD, "/data", r.Method)
	}
}

func RecordHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(rw).Encode(state)
	} else {
		log.Printf(ERR_METHOD, "/data/json", r.Method)
	}
}

func LogHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rqbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		query := string(rqbody)
		response := state.ProcessQuery(query)

		log := state.RawLog[len(state.RawLog)-1]
		fmt.Printf("[%d] %s\n", log.Timestamp, log.Query)

		fmt.Fprint(rw, response)
	} else {
		log.Printf(ERR_METHOD, "/log", r.Method)
	}
}

func SinceHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ts := r.URL.Query().Get("t")
		if ts == "" {
			ts = "0"
		}
		timestamp, err := strconv.Atoi(ts)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Failed to convert t into an int with error: %s", err), 400)
			return
		}
		data := state.LogsSince(int64(timestamp))
		json.NewEncoder(rw).Encode(data)
	} else {
		log.Printf(ERR_METHOD, "/data/since", r.Method)
	}
}

func KeyvaluesHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(rw).Encode(state.KeyValues)
	} else {
		log.Printf(ERR_METHOD, "/data/keyvalues", r.Method)
	}
}

func RunServe(cmd *cobra.Command, args []string) {
	state = NewServerState()

	// getting the static web page
	http.HandleFunc("/data", WebPageHandler)

	// getting the data in json form
	http.HandleFunc("/data/json", RecordHandler)

	// getting list of logs since a time
	http.HandleFunc("/data/since", SinceHandler)

	// getting keyvalue data
	http.HandleFunc("/data/keyvalues", KeyvaluesHandler)

	// adding log information
	http.HandleFunc("/log", LogHandler)

	fmt.Printf("Listening on http://localhost:%d\n", port)
	fmt.Printf("Access data on http://localhost:%d/data\n", port)
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
