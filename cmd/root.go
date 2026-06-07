package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mohammadctl",
	Short: "mohammadctl is a very useless cli",
	Long:  `in fact, mohammadctl is the most useless cli.`,
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
