package basecrm

import (
	"fmt"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // what status value can have ?
	Role      string    `json:"role"`   // what roles are available ?
	Confirmed bool      `json:"confirmed"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserListOptions struct {
	Name      string `url:"name,omitempty"`
	Email     string `url:"email,omitempty"`
	Role      string `url:"role,omitempty"`
	Status    string `url:"status,omitempty"`
	Confirmed bool   `url:"confirmed,omitempty"`

	ListOptions
}

type UsersService interface {
	List(opt *UserListOptions) ([]*User, *Response, error)
	Get(id int) (*User, *Response, error)
	Self() (*User, *Response, error)
}

func NewUsersService(client *Client) UsersService {
	return &UsersServiceOp{client}
}

type userRoot struct {
	User *User `json:"data"`
	Meta *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type usersRoot struct {
	Items []*userRoot `json:"items"`
	Meta  *Meta       `json:"meta"`
}

func (r *usersRoot) Users() []*User {
	users := make([]*User, len(r.Items))
	for i, root := range r.Items {
		users[i] = root.User
	}
	return users
}

type UsersServiceOp struct {
	client *Client
}

func (s *UsersServiceOp) List(opt *UserListOptions) ([]*User, *Response, error) {
	u, err := addOptions("/v2/users", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(usersRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Users(), res, err
}

func (s *UsersServiceOp) Get(id int) (*User, *Response, error) {
	u := fmt.Sprintf("/v2/users/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.User, res, err
}

func (s *UsersServiceOp) Self() (*User, *Response, error) {
	u := "/v2/users/self"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.User, res, err
}
