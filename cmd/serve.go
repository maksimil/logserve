package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	port int
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs logging web api on localhost",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/record", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Printf("You requested: %s\n", r.URL.Path)
		})

		fmt.Printf("Listening on http://localhost:%d/record\n", port)
		if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8000, "--port <PORT NUMBER>")
}
