package access

import (
	"log"
	"mammadctl/configs"
	"mammadctl/internal/sotoon"
	"os"

	"github.com/spf13/cobra"
)

var nedaS3BaseUrl, afraS3BaseUrl string

var S3Cmd = &cobra.Command{
	Use:   "s3 [user's uuid] [datacenter] [bucket's name]",
	Short: "add full access to the bucket",
	Args:  cobra.ExactArgs(3),
	RunE:  runS3,
}

func init() {
	S3Cmd.PersistentFlags().StringVar(&nedaS3BaseUrl, "neda-s3-url", "s3.thr1.sotoon.ir", "Neda's s3 datacenter base url")
	S3Cmd.PersistentFlags().StringVar(&afraS3BaseUrl, "afra-s3-url", "s3.thr2.sotoon.ir", "Afra's s3 datacenter base url")
}

func runS3(cmd *cobra.Command, args []string) error {
	userId, datacenter, bucket := args[0], args[1], args[2]

	switch datacenter {
	case "neda":
		configs.SotoonConfig.S3Endpoint = nedaS3BaseUrl
	case "afra":
		configs.SotoonConfig.S3Endpoint = afraS3BaseUrl
	default:
		log.Fatalf("Unsupported data center %s", datacenter)
	}

	configs.SotoonConfig.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	if configs.SotoonConfig.S3AccessKey == "" {
		log.Fatalln("S3_ACCESS_KEY environment variable not set")
	}

	configs.SotoonConfig.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	if configs.SotoonConfig.S3SecretKey == "" {
		log.Fatalln("S3_SECRET_KEY environment variable not set")
	}

	as := sotoon.NewAccessService(configs.SotoonConfig)
	return as.AddS3Policy(userId, bucket)
}
