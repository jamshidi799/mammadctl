package sotoon

import (
	"mammadctl/cmd/sotoon/user"

	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Adding an access to a user",
}

func init() {
	UserCmd.AddCommand(user.ListCmd)
}
