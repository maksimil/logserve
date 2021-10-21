package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logserve",
	Short: "Runs a local logserve webserver",
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

var (
	port int
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8000, "--port <PORT NUMBER>")
}
