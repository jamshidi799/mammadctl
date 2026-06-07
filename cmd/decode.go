package cmd

import (
	"mammadctl/cmd/decode"

	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode staff(for now just secrets)",
}

func init() {
	decodeCmd.AddCommand(decode.SecretCmd)
}
