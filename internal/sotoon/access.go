package sotoon

import (
	"fmt"
	"mammadctl/configs"
	"strings"
)

const (
	cdnEditorRoleId                 = "38aba8b2-feb0-43b7-8f7d-8a0772d79f40"
	dnsViewerRoleId                 = "8ddd6ab8-9ed1-4eab-b586-7c286d170ddf"
	kaasClusterNamespaceAdminRoleId = "ba6f4b9c-ca8d-4d0a-91ba-631d8794b87b"
	skeClusterNamespaceAdminRoleId  = "b870d664-1bc6-42db-a128-8ee8df64cc40"
)

const (
	v1K8sClusterName = "bazaar-production"

	v2K8sZone        = "thr1"
	v2K8sTier        = "premium"
	V2K8sClusterName = "cafebazaar"
)

const S3PolicyArnPrefix = "arn:aws:iam:::user"

type policyDocument struct {
	Version   string      `json:"Version"`
	Statement []statement `json:"Statement"`
}

type statement struct {
	Sid       string      `json:"Sid,omitempty"`
	Effect    string      `json:"Effect"`
	Principal interface{} `json:"Principal"`
	Action    interface{} `json:"Action"`
	Resource  interface{} `json:"Resource"`
}

type principal struct {
	AWS []string `json:"AWS"`
}

type AccessService interface {
	AddCdnRole(userId, domain string) error
	AddK8sRole(userId string, datacenter string, namespace string) error
	AddS3Policy(userId string, bucket string) error
}

type accessServiceImpl struct {
	client accessClient
	config configs.Sotoon
}

func NewAccessService(c configs.Sotoon) AccessService {
	return &accessServiceImpl{
		client: newClient(c),
		config: c,
	}
}

func (a *accessServiceImpl) AddK8sRole(userId string, datacenter string, namespace string) error {
	switch datacenter {
	case "afra":
		addKaasClusterNamespaceAdminRoleRequest := addRoleRequest{
			roleId: kaasClusterNamespaceAdminRoleId,

			Users: []string{userId},
			Items: []addRoleItem{
				{
					Datacenter: datacenter,
					Cluster:    v1K8sClusterName,
					Namespace:  namespace,
				},
			},
		}
		if err := a.client.addRoleToUsers(addKaasClusterNamespaceAdminRoleRequest); err != nil {
			return err
		}

	case "neda":
		addSkeClusterNamespaceAdminRoleRequest := addRoleRequest{
			roleId: skeClusterNamespaceAdminRoleId,

			Users: []string{userId},
			Items: []addRoleItem{
				{
					Datacenter: datacenter,
					Zone:       v2K8sZone,
					Tier:       v2K8sTier,
					Cluster:    V2K8sClusterName,
					Namespace:  namespace,
				},
			},
		}
		if err := a.client.addRoleToUsers(addSkeClusterNamespaceAdminRoleRequest); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid datacenter %s. expected afra or neda", datacenter)
	}

	return nil
}

func (a *accessServiceImpl) AddCdnRole(userId string, domain string) error {
	addCdnEditorRoleRequest := addRoleRequest{
		roleId: cdnEditorRoleId,

		Users: []string{userId},
		Items: []addRoleItem{{Cdn: domain}},
	}

	if err := a.client.addRoleToUsers(addCdnEditorRoleRequest); err != nil {
		return err
	}

	domainzone := strings.Replace(domain, ".", "-", -1)
	dnssec := domainzone
	addDnsViewerRoleRequest := addRoleRequest{
		roleId: dnsViewerRoleId,

		Users: []string{userId},
		Items: []addRoleItem{{Domainzone: domainzone, Dnssec: dnssec}},
	}

	if err := a.client.addRoleToUsers(addDnsViewerRoleRequest); err != nil {
		return err
	}

	return nil
}

func (a *accessServiceImpl) AddS3Policy(userId string, bucket string) error {
	policy, err := a.client.getBucketPolicy(bucket)
	if err != nil {
		return err
	}

	isStarActionFound := false
	for i, s := range policy.Statement {
		if s.Action == "*" {
			isStarActionFound = true

			p, err := parsePrinciple(s.Principal)
			if err != nil {
				return err
			}

			newArn := fmt.Sprintf("%s/%s:%s", S3PolicyArnPrefix, a.config.WorkspaceId, userId)
			p.AWS = append(p.AWS, newArn)
			policy.Statement[i].Principal = p
			break
		}
	}

	if !isStarActionFound {
		return fmt.Errorf("action * not fount in bucket policy")
	}

	if policy.Version == "" {
		policy.Version = "2012-10-17"
	}

	err = a.client.setBucketPolicy(bucket, *policy)
	return err
}

func parsePrinciple(obj interface{}) (*principal, error) {
	var p principal

	stringToInterfaceMap, ok := obj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid principal, convertion to map failed: %+v", obj)
	}

	list, ok := stringToInterfaceMap["AWS"]
	if !ok {
		return nil, fmt.Errorf("invalid principal, AWS key not found in principal: %+v", obj)
	}

	arnList, ok := list.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid principal, arn convertion failed: %+v", obj)
	}

	for _, l := range arnList {
		arn, _ := l.(string)
		p.AWS = append(p.AWS, arn)
	}

	return &p, nil
}
