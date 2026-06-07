package decode

import (
	"log"
	"mammadctl/internal/decode"
	"mammadctl/pkg/k8s"
	"os"

	"github.com/spf13/cobra"
)

var namespace string

var SecretCmd = &cobra.Command{
	Use:   "secret -n [secret's namespace] [secret's name to decode]",
	Short: "Find the secret and write it to ./secret.yaml in stringData mode.",
	Args:  cobra.ExactArgs(1),
	Run:   runSecret,
}

func init() {
	SecretCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace")
}

func runSecret(command *cobra.Command, args []string) {
	name := args[0]

	clientSet := k8s.BuildClientset()

	file, err := os.OpenFile("secret.yaml", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	d := decode.NewSecretDecoder(clientSet, file)

	err = d.Decode(name, namespace)
	if err != nil {
		log.Fatal(err)
	}
}
