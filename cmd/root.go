package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logserve",
	Short: "Runs a local logserve webserver",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi", port)
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
