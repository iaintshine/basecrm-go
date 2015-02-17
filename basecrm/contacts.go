package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type Contact struct {
	Id             int                    `json:"id,omitempty"`
	CreatorId      int                    `json:"creator_id,omitempty"`
	OwnerId        int                    `json:"owner_id,omitempty"`
	IsOrganization bool                   `json:"is_organization,omitempty"`
	Private        bool                   `json:"private"`
	Name           string                 `json:"name,omitempty"`
	FirstName      string                 `json:"first_name,omitempty"`
	LastName       string                 `json:"last_name,omitempty"`
	CustomerStatus string                 `json:"customer_status,omitempty"`
	ProspectStatus string                 `json:"prospect_status,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Industry       string                 `json:"industry,omitempty"`
	Website        string                 `json:"website,omitempty"`
	Email          string                 `json:"email,omitempty"`
	Phone          string                 `json:"phone,omitempty"`
	Mobile         string                 `json:"mobile,omitempty"`
	Fax            string                 `json:"fax,omitempty"`
	Twitter        string                 `json:"twitter,omitempty"`
	Facebook       string                 `json:"facebook,omitempty"`
	Linkedin       string                 `json:"linkedin,omitempty"`
	Skype          string                 `json:"skype,omitempty"`
	Address        *Address               `json:"address,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	CustomFields   map[string]interface{} `json:"custom_fields,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty"`
}

type ContactListOptions struct {
	Q      string `url:"q,omitempty"`
	Letter string `url:"letter,omitempty"`

	CreatorId int `url:"creator_id,omitempty"`
	OwnerId   int `url:"owner_id,omitempty"`

	IsOrganization bool `url:"is_organization,omitempty"`

	Name      string `url:"name,omitempty"`
	FirstName string `url:"first_name,omitempty"`
	LastName  string `url:"last_name,omitempty"`

	CustomerStatus string `url:"customer_status,omitempty"`
	ProspectStatus string `url:"prospect_status,omitempty"`

	City       string `url:"address[city],omitempty"`
	PostalCode string `url:"address[postal_code],omitempty"`
	Country    string `url:"address[country],omitempty"`

	ListOptions
}

type ContactsService interface {
	List(opt *ContactListOptions) ([]*Contact, *Response, error)
	Get(id int) (*Contact, *Response, error)
	Create(contact *Contact) (*Contact, *Response, error)
	Edit(id int, contact *Contact) (*Contact, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewContactsService(client *Client) ContactsService {
	return &ContactsServiceOp{client}
}

type contactRoot struct {
	Contact *Contact `json:"data"`
	Meta    *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type contactsRoot struct {
	Items []*contactRoot `json:"items"`
	Meta  *Meta          `json:"meta"`
}

func (r *contactsRoot) Contacts() []*Contact {
	contacts := make([]*Contact, len(r.Items))
	for i, root := range r.Items {
		contacts[i] = root.Contact
	}
	return contacts
}

type ContactsServiceOp struct {
	client *Client
}

func (s *ContactsServiceOp) List(opt *ContactListOptions) ([]*Contact, *Response, error) {
	u, err := addOptions("/v2/contacts", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(contactsRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Contacts(), res, err
}

func (s *ContactsServiceOp) Get(id int) (*Contact, *Response, error) {
	u := fmt.Sprintf("/v2/contacts/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(contactRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Contact, res, err
}

func (s *ContactsServiceOp) Create(contact *Contact) (*Contact, *Response, error) {
	u := "/v2/contacts"
	envelope := &contactRoot{Contact: contact}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(contactRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Contact, res, err
}

func (s *ContactsServiceOp) Edit(id int, contact *Contact) (*Contact, *Response, error) {
	u := fmt.Sprintf("/v2/contacts/%d", id)
	envelope := &contactRoot{Contact: contact}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(contactRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Contact, res, err
}

func (s *ContactsServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/contacts/%d", id)
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
