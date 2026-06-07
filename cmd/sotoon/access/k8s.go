package access

import (
	"mammadctl/configs"
	"mammadctl/internal/sotoon"

	"github.com/spf13/cobra"
)

var K8sCmd = &cobra.Command{
	Use:   "k8s [user's uuid] [datacenter] [namespace]",
	Short: "Add namespaced-scope access to the specified user",
	Args:  cobra.ExactArgs(3),
	RunE:  runK8s,
}

func runK8s(cmd *cobra.Command, args []string) error {
	userId, datacenter, namespace := args[0], args[1], args[2]
	as := sotoon.NewAccessService(configs.SotoonConfig)
	return as.AddK8sRole(userId, datacenter, namespace)
}
