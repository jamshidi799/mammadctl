package sotoon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mammadctl/configs"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

type userClient interface {
	getUsers() ([]userResponse, error)
	getServiceUsers() ([]serviceUserResponse, error)
}

type accessClient interface {
	addRoleToUsers(request addRoleRequest) error
	getBucketPolicy(bucket string) (*policyDocument, error)
	setBucketPolicy(bucket string, policy policyDocument) error
}

type client struct {
	httpClient  *http.Client
	apiBaseUrl  string
	workspaceId string
	token       string

	s3Client *minio.Client
}

func newClient(config configs.Sotoon) *client {
	c := &client{
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		apiBaseUrl:  config.ApiBaseUrl,
		workspaceId: config.WorkspaceId,
		token:       config.Token,
	}

	if config.S3Endpoint == "" {
		return c
	}

	s3Client, err := minio.New(config.S3Endpoint, config.S3AccessKey, config.S3SecretKey, false)
	if err != nil {
		log.Fatal(err)
	}

	c.s3Client = s3Client
	return c
}

func (c *client) getUsers() ([]userResponse, error) {
	url := fmt.Sprintf("%s/iam/v1/api/v1/detailed/workspace/%s/user", c.apiBaseUrl, c.workspaceId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	users, err := parse[[]userResponse](resp.Body)
	resp.Body.Close()

	return users, err
}

func (c *client) getServiceUsers() ([]serviceUserResponse, error) {
	url := fmt.Sprintf("%s/iam/v1/api/v1/detailed/workspace/%s/service-user", c.apiBaseUrl, c.workspaceId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	serviceUsers, err := parse[[]serviceUserResponse](resp.Body)
	resp.Body.Close()

	return serviceUsers, err
}

func (c *client) addRoleToUsers(request addRoleRequest) error {
	url := fmt.Sprintf("%s/iam/v1/api/v1/workspace/%s/role/%s/bulk-add-users/", c.apiBaseUrl, c.workspaceId, request.roleId)
	body, _ := json.Marshal(request)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}
	return nil
}

func (c *client) getBucketPolicy(bucket string) (*policyDocument, error) {
	resp, err := c.s3Client.GetBucketPolicy(bucket)
	if err != nil {
		return nil, err
	}

	policy, err := parse[policyDocument](strings.NewReader(resp))
	return &policy, err
}

func (c *client) setBucketPolicy(bucket string, policy policyDocument) error {
	j, err := json.MarshalIndent(policy, "", "    ")
	if err != nil {
		return err
	}

	err = c.s3Client.SetBucketPolicy(bucket, string(j))
	if err != nil {
		return err
	}

	return nil
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func parse[T any](body io.Reader) (t T, err error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(bytes, &t)

	return t, err
}

type userResponse struct {
	Email string `json:"email"`
	Uuid  string `json:"uuid"`
}

type serviceUserResponse struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

type addRoleRequest struct {
	roleId string

	Users []string      `json:"users"`
	Items []addRoleItem `json:"items"`
}

type addRoleItem struct {
	Cdn string `json:"cdn,omitempty"`

	Domainzone string `json:"domainzone,omitempty"`
	Dnssec     string `json:"dnssec,omitempty"`

	Datacenter string `json:"datacenter,omitempty"`
	Cluster    string `json:"cluster,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Tier       string `json:"tier,omitempty"`
	Zone       string `json:"zone,omitempty"`
}
