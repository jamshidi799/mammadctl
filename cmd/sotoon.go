package cmd

import (
	"log"
	"mammadctl/cmd/sotoon"
	"mammadctl/configs"
	"os"

	"github.com/spf13/cobra"
)

var sotoonCmd = &cobra.Command{
	Use:   "sotoon",
	Short: "Some utilities for working with Sotoon api",
}

func init() {
	sotoonCmd.AddCommand(sotoon.UserCmd, sotoon.AccessCmd)
	sotoonCmd.PersistentFlags().StringVar(&configs.SotoonConfig.ApiBaseUrl, "url", "https://api.sotoon.ir", "Sotoon base url")
	sotoonCmd.PersistentFlags().StringVar(&configs.SotoonConfig.WorkspaceId, "workspace-id", "fee4bbf5-342d-4243-b4aa-f0bee508c39a", "Your workspace id")

	configs.SotoonConfig.Token = os.Getenv("BEPA_TOKEN")
	if configs.SotoonConfig.Token == "" {
		log.Fatalln("BEPA_TOKEN environment variable not set")
	}
}
