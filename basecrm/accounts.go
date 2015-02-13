package basecrm

import (
	"time"
)

type Account struct {
	Id         int       `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Role       string    `json:"role,omitempty"`
	Plan       string    `json:"plan,omitempty"`
	Currency   string    `json:"currency,omitempty"`
	TimeFormat string    `json:"time_format,omitempty"`
	Timezone   string    `json:"timezone,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

type AccountsService interface {
	Self() (*Account, *Response, error)
}

func NewAccountsService(client *Client) AccountsService {
	return &AccountsServiceOp{client}
}

type AccountsServiceOp struct {
	client *Client
}

type accountRoot struct {
	Account *Account `json:"data"`
	Meta    *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

func (s *AccountsServiceOp) Self() (*Account, *Response, error) {
	u := "/v2/accounts/self"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(accountRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Account, res, err
}
