package configs

var SotoonConfig Sotoon

type Sotoon struct {
	ApiBaseUrl  string
	Token       string
	WorkspaceId string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
}
