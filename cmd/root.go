package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "logserve",
	Short: "An app for logging data using a web api on localhost",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}
