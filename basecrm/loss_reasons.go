package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type LossReason struct {
	Id        int       `json:"id,omitempty"`
	CreatorId int       `json:"creator_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type LossReasonListOptions struct {
	CreatorId int    `url:"creator_id,omitempty"`
	Name      string `url:"name,omitempty"`

	ListOptions
}

type LossReasonsService interface {
	List(opt *LossReasonListOptions) ([]*LossReason, *Response, error)
	Get(id int) (*LossReason, *Response, error)
	Create(lossReason *LossReason) (*LossReason, *Response, error)
	Edit(id int, lossReason *LossReason) (*LossReason, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewLossReasonsService(client *Client) LossReasonsService {
	return &LossReasonsServiceOp{client}
}

type lossReasonRoot struct {
	LossReason *LossReason `json:"data"`
	Meta       *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type lossReasonsRoot struct {
	Items []*lossReasonRoot `json:"items"`
	Meta  *Meta             `json:"meta"`
}

func (r *lossReasonsRoot) LossReasons() []*LossReason {
	contacts := make([]*LossReason, len(r.Items))
	for i, root := range r.Items {
		contacts[i] = root.LossReason
	}
	return contacts
}

type LossReasonsServiceOp struct {
	client *Client
}

func (s *LossReasonsServiceOp) List(opt *LossReasonListOptions) ([]*LossReason, *Response, error) {
	u, err := addOptions("/v2/loss_reasons", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(lossReasonsRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.LossReasons(), res, err
}

func (s *LossReasonsServiceOp) Get(id int) (*LossReason, *Response, error) {
	u := fmt.Sprintf("/v2/loss_reasons/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(lossReasonRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.LossReason, res, err
}

func (s *LossReasonsServiceOp) Create(lossReason *LossReason) (*LossReason, *Response, error) {
	u := "/v2/loss_reasons"
	envelope := &lossReasonRoot{LossReason: lossReason}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(lossReasonRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.LossReason, res, err
}

func (s *LossReasonsServiceOp) Edit(id int, lossReason *LossReason) (*LossReason, *Response, error) {
	u := fmt.Sprintf("/v2/loss_reasons/%d", id)
	envelope := &lossReasonRoot{LossReason: lossReason}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(lossReasonRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.LossReason, res, err
}

func (s *LossReasonsServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/loss_reasons/%d", id)
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
