package access

import (
	"mammadctl/configs"
	"mammadctl/internal/sotoon"

	"github.com/spf13/cobra"
)

var CdnCmd = &cobra.Command{
	Use:   "cdn [user's uuid] [domain]",
	Short: "Add cdn-editor and dns-viewer role to the specified user",
	Args:  cobra.ExactArgs(2),
	RunE:  runCdn,
}

func runCdn(cmd *cobra.Command, args []string) error {
	uuid, domain := args[0], args[1]
	as := sotoon.NewAccessService(configs.SotoonConfig)
	return as.AddCdnRole(uuid, domain)
}
