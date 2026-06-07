package sotoon

import (
	"errors"
	"reflect"
	"testing"
)

type fakeUserClient struct {
	users        []userResponse
	serviceUsers []serviceUserResponse
	userErr      error
	serviceErr   error
}

func (f *fakeUserClient) getUsers() ([]userResponse, error) {
	return f.users, f.userErr
}

func (f *fakeUserClient) getServiceUsers() ([]serviceUserResponse, error) {
	return f.serviceUsers, f.serviceErr
}

func TestUserServiceImpl_List(t *testing.T) {
	fake := &fakeUserClient{
		users: []userResponse{
			{"test@example.com", "u1"},
		},
		serviceUsers: []serviceUserResponse{
			{"test-service-user", "u2"},
		},
		userErr:    nil,
		serviceErr: nil,
	}

	svc := userServiceImpl{client: fake}
	got, err := svc.List()

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []User{
		{"user", "test@example.com", "u1"},
		{"service-user", "test-service-user", "u2"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() got = %#v, want %#v", got, want)
	}
}

func TestUserServiceImpl_List_GetUsersError(t *testing.T) {
	fake := &fakeUserClient{
		userErr: errors.New("some error"),
	}

	svc := userServiceImpl{client: fake}
	got, err := svc.List()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}
