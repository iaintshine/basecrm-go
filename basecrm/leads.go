package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type Lead struct {
	Id               int                    `json:"id,omitempty"`
	CreatorId        int                    `json:"creator_id,omitempty"`
	OwnerId          int                    `json:"owner_id,omitempty"`
	FirstName        string                 `json:"first_name,omitempty"`
	LastName         string                 `json:"last_name,omitempty"`
	OrganizationName string                 `json:"organization_name,omitempty"`
	Status           string                 `json:"status,omitempty"`
	Title            string                 `json:"title,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Industry         string                 `json:"industry,omitempty"`
	Website          string                 `json:"website,omitempty"`
	Email            string                 `json:"email,omitempty"`
	Phone            string                 `json:"phone,omitempty"`
	Mobile           string                 `json:"mobile,omitempty"`
	Fax              string                 `json:"fax,omitempty"`
	Twitter          string                 `json:"twitter,omitempty"`
	Facebook         string                 `json:"facebook,omitempty"`
	Linkedin         string                 `json:"linkedin,omitempty"`
	Skype            string                 `json:"skype,omitempty"`
	Address          *Address               `json:"address,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	CustomFields     map[string]interface{} `json:"custom_fields,omitempty"`
	UpdatedAt        time.Time              `json:"updated_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at,omitempty"`
}

type LeadListOptions struct {
	Q      string `url:"q,omitempty"`
	Letter string `url:"letter,omitempty"`

	CreatorId int `url:"creator_id,omitempty"`
	OwnerId   int `url:"owner_id,omitempty"`

	FirstName        string `url:"first_name,omitempty"`
	LastName         string `url:"last_name,omitempty"`
	OrganizationName string `url:"organization_name,omitempty"`

	Status string `url:"status,omitempty"`

	City       string `url:"address[city],omitempty"`
	PostalCode string `url:"address[postal_code],omitempty"`
	Country    string `url:"address[country],omitempty"`

	ListOptions
}

type LeadConversionOptions struct {
	Includes       []string
	DealName       string
	TaskContent    string
	TaskDueDate    time.Time
	OrganizationId int
	OwnerId        int
	Strict         bool
}

type LeadsService interface {
	List(opt *LeadListOptions) ([]*Lead, *Response, error)
	Get(id int) (*Lead, *Response, error)
	Create(lead *Lead) (*Lead, *Response, error)
	Edit(id int, lead *Lead) (*Lead, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewLeadsService(client *Client) LeadsService {
	return &LeadsServiceOp{client}
}

type leadRoot struct {
	Lead *Lead `json:"data"`
	Meta *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type leadsRoot struct {
	Items []*leadRoot `json:"items"`
	Meta  *Meta       `json:"meta"`
}

func (r *leadsRoot) Leads() []*Lead {
	leads := make([]*Lead, len(r.Items))
	for i, root := range r.Items {
		leads[i] = root.Lead
	}
	return leads
}

type LeadsServiceOp struct {
	client *Client
}

func (s *LeadsServiceOp) List(opt *LeadListOptions) ([]*Lead, *Response, error) {
	u, err := addOptions("/v2/leads", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(leadsRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Leads(), res, err
}

func (s *LeadsServiceOp) Get(id int) (*Lead, *Response, error) {
	u := fmt.Sprintf("/v2/leads/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(leadRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Lead, res, err
}

func (s *LeadsServiceOp) Create(lead *Lead) (*Lead, *Response, error) {
	u := "/v2/leads"
	envelope := &leadRoot{Lead: lead}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(leadRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Lead, res, err
}

func (s *LeadsServiceOp) Edit(id int, lead *Lead) (*Lead, *Response, error) {
	u := fmt.Sprintf("/v2/leads/%d", id)
	envelope := &leadRoot{Lead: lead}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(leadRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Lead, res, err
}

func (s *LeadsServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/leads/%d", id)
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
