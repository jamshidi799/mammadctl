package cmd

import (
	"fmt"
	"log"
	"mammadctl/internal/convert"
	"mammadctl/pkg/k8s"
	"os"

	"github.com/spf13/cobra"
)

var (
	namespace string
	cluster   string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert staff(for now just services)",
	Args:  cobra.ExactArgs(1),
	Run:   runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace")
	convertCmd.Flags().StringVarP(&cluster, "cluster", "c", "default", "cluster")
}

func runConvert(command *cobra.Command, args []string) {
	name := args[0]

	clientSet := k8s.BuildClientset()

	fileName := fmt.Sprintf("%s-%s.yaml", namespace, name)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	converter := convert.NewServiceConverter(clientSet, file)

	err = converter.Convert(name, namespace, cluster)
	if err != nil {
		log.Fatal(err)
	}
}
