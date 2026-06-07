package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mammadctl",
	Short: "mammadctl is a very useless cli",
	Long:  `in fact, mammadctl is the most useless cli.`,
}

func init() {
	rootCmd.AddCommand(decodeCmd, sotoonCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
