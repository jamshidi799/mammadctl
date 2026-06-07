package user

import (
	"fmt"
	"log"
	"mammadctl/configs"
	"mammadctl/internal/sotoon"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Run:   runList,
}

func runList(_ *cobra.Command, args []string) {
	us := sotoon.NewUserService(configs.SotoonConfig)
	users, err := us.List()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		fmt.Printf("type:%s, name: %s, id: %s\n", user.Type, user.Name, user.Id)
	}
}
