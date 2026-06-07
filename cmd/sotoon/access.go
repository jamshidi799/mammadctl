package sotoon

import (
	"mammadctl/cmd/sotoon/access"

	"github.com/spf13/cobra"
)

var AccessCmd = &cobra.Command{
	Use:   "access",
	Short: "",
}

func init() {
	AccessCmd.AddCommand(access.CdnCmd, access.K8sCmd, access.S3Cmd)
}
