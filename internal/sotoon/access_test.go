package sotoon

import (
	"errors"
	"reflect"
	"testing"
)

type fakeAccessClient struct {
	addRoleCalls []addRoleRequest
	addRoleErr   error

	getPolicyBucket string
	getPolicyResp   *policyDocument
	getPolicyErr    error

	setPolicyBucket string
	setPolicyInput  policyDocument
	setPolicyErr    error
}

func (f *fakeAccessClient) addRoleToUsers(request addRoleRequest) error {
	f.addRoleCalls = append(f.addRoleCalls, request)
	return f.addRoleErr
}

func (f *fakeAccessClient) getBucketPolicy(bucket string) (*policyDocument, error) {
	f.getPolicyBucket = bucket
	return f.getPolicyResp, f.getPolicyErr
}

func (f *fakeAccessClient) setBucketPolicy(bucket string, policy policyDocument) error {
	f.setPolicyBucket = bucket
	f.setPolicyInput = policy
	return f.setPolicyErr
}

func TestAccessServiceImpl_AddK8sRole(t *testing.T) {
	request := addRoleRequest{
		roleId: kaasClusterNamespaceAdminRoleId,
		Users:  []string{"id1"},
		Items: []addRoleItem{
			{
				Datacenter: "afra",
				Cluster:    v1K8sClusterName,
				Namespace:  "default",
			},
		},
	}

	f := &fakeAccessClient{}
	svc := accessServiceImpl{
		client: f,
	}

	err := svc.AddK8sRole(request.Users[0], request.Items[0].Datacenter, request.Items[0].Namespace)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(f.addRoleCalls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(f.addRoleCalls))
	}

	if !reflect.DeepEqual(f.addRoleCalls[0], request) {
		t.Fatalf("expected %v, got %v", request, f.addRoleCalls[0])
	}
}

func TestAccessServiceImpl_AddK8sRole_InvalidDatacenter(t *testing.T) {
	request := addRoleRequest{
		roleId: kaasClusterNamespaceAdminRoleId,
		Users:  []string{"id1"},
		Items: []addRoleItem{
			{
				Datacenter: "unknown",
				Cluster:    v1K8sClusterName,
				Namespace:  "default",
			},
		},
	}

	f := &fakeAccessClient{}
	svc := accessServiceImpl{client: f}

	err := svc.AddK8sRole(request.Users[0], request.Items[0].Datacenter, request.Items[0].Namespace)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if len(f.addRoleCalls) > 0 {
		t.Fatalf("expected 0 call, got %d", len(f.addRoleCalls))
	}
}

func TestAccessServiceImpl_AddK8sRole_AddRoleToUsersError(t *testing.T) {
	request := addRoleRequest{
		roleId: kaasClusterNamespaceAdminRoleId,
		Users:  []string{"id1"},
		Items: []addRoleItem{
			{
				Datacenter: "afra",
				Cluster:    v1K8sClusterName,
				Namespace:  "default",
			},
		},
	}

	f := &fakeAccessClient{
		addRoleErr: errors.New("AddRoleToUsersError"),
	}
	svc := accessServiceImpl{client: f}

	err := svc.AddK8sRole(request.Users[0], request.Items[0].Datacenter, request.Items[0].Namespace)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if len(f.addRoleCalls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(f.addRoleCalls))
	}

	if !reflect.DeepEqual(f.addRoleCalls[0], request) {
		t.Fatalf("expected %v, got %v", request, f.addRoleCalls[0])
	}
}
