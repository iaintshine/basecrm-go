package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type AssociatedContact struct {
	CreatorId int    `json:"creator_id,omitempty"`
	ContactId int    `json:"contact_id,omitempty"`
	Role      string `json:"role,omitempty"`
}

type Deal struct {
	Id                 int                    `json:"id,omitempty"`
	CreatorId          int                    `json:"creator_id,omitempty"`
	OwnerId            int                    `json:"owner_id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Value              int                    `json:"value,omitempty"`
	Currency           string                 `json:"currency,omitempty"`
	Hot                bool                   `json:"hot,omitempty"`
	SourceId           int                    `json:"source_id,omitempty"`
	LossReasonId       int                    `json:"loss_reason_id,omitempty"`
	AssociatedContacts []*AssociatedContact   `json:"associated_contacts,omitempty"`
	DropboxEmail       string                 `json:"dropbox_email,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	CustomFields       map[string]interface{} `json:"custom_fields,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
}

type DealListOptions struct {
	Q    string `url:"q,omitempty"`
	Name string `url:"name,omitempty"`

	CreatorId int `url:"creator_id,omitempty"`
	OwnerId   int `url:"owner_id,omitempty"`
	ContactId int `url:"contact_id,omitempty"`

	SourceId     int `url:"source_id,omitempty"`
	LossReasonId int `url:"loss_reason_id,omitempty"`

	Hot bool `url:"hot,omitempty"`

	ListOptions
}

type DealsService interface {
	List(opt *DealListOptions) ([]*Deal, *Response, error)
	Get(id int) (*Deal, *Response, error)
	Create(deal *Deal) (*Deal, *Response, error)
	Edit(id int, deal *Deal) (*Deal, *Response, error)
	Delete(id int) (bool, *Response, error)
	UpsertContact(id int, contact *AssociatedContact) (bool, *Response, error)
	DeleteContact(id, contactId int) (bool, *Response, error)
}

func NewDealsService(client *Client) DealsService {
	return &DealsServiceOp{client}
}

type dealRoot struct {
	Deal *Deal `json:"data"`
	Meta *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type dealsRoot struct {
	Items []*dealRoot `json:"items"`
	Meta  *Meta       `json:"meta"`
}

func (r *dealsRoot) Deals() []*Deal {
	deals := make([]*Deal, len(r.Items))
	for i, root := range r.Items {
		deals[i] = root.Deal
	}
	return deals
}

type DealsServiceOp struct {
	client *Client
}

func (s *DealsServiceOp) List(opt *DealListOptions) ([]*Deal, *Response, error) {
	u, err := addOptions("/v2/deals", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dealsRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Deals(), res, err
}

func (s *DealsServiceOp) Get(id int) (*Deal, *Response, error) {
	u := fmt.Sprintf("/v2/deals/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dealRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Deal, res, err
}

func (s *DealsServiceOp) Create(deal *Deal) (*Deal, *Response, error) {
	u := "/v2/deals"
	envelope := &dealRoot{Deal: deal}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(dealRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Deal, res, err
}

func (s *DealsServiceOp) Edit(id int, deal *Deal) (*Deal, *Response, error) {
	u := fmt.Sprintf("/v2/deals/%d", id)
	envelope := &dealRoot{Deal: deal}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(dealRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Deal, res, err
}

func (s *DealsServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/deals/%d", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return false, nil, err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return false, res, err
	}

	return res.StatusCode == http.StatusNoContent, res, err
}

func (s *DealsServiceOp) UpsertContact(id int, contact *AssociatedContact) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/deals/%d/associated_contacts/%d?role=%s", id, contact.ContactId, contact.Role)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return false, nil, err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return false, res, err
	}

	return res.StatusCode == http.StatusNoContent, res, err
}

func (s *DealsServiceOp) DeleteContact(id, contactId int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/deals/%d/associated_contacts/%d", id, contactId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return false, nil, err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return false, res, err
	}

	return res.StatusCode == http.StatusNoContent, res, err
}
