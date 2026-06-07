package sotoon

import "mammadctl/configs"

type User struct {
	Type string
	Name string
	Id   string
}

type UserService interface {
	List() ([]User, error)
}

type userServiceImpl struct {
	client userClient
}

func NewUserService(c configs.Sotoon) UserService {
	return &userServiceImpl{client: newClient(c)}
}

func (u *userServiceImpl) List() ([]User, error) {
	users, err := u.client.getUsers()
	if err != nil {
		return nil, err
	}

	serviceUsers, err := u.client.getServiceUsers()
	if err != nil {
		return nil, err
	}

	list := make([]User, 0, len(users)+len(serviceUsers))
	for _, u := range users {
		user := User{
			Type: "user",
			Name: u.Email,
			Id:   u.Uuid,
		}
		list = append(list, user)
	}

	for _, su := range serviceUsers {
		serviceUser := User{
			Type: "service-user",
			Name: su.Name,
			Id:   su.Uuid,
		}
		list = append(list, serviceUser)
	}

	return list, nil
}
