package basecrm

import (
	"time"
)

type Account struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Plan       string    `json:"plan"`
	Currency   string    `json:"currency"`
	TimeFormat string    `json:"time_format"`
	Timezone   string    `json:"timezone"`
	Phone      string    `json:"phone"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
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
