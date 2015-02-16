package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type Source struct {
	Id        int       `json:"id,omitempty"`
	CreatorId int       `json:"creator_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type SourceListOptions struct {
	CreatorId int    `url:"creator_id,omitempty"`
	Name      string `url:"name,omitempty"`

	ListOptions
}

type SourcesService interface {
	List(opt *SourceListOptions) ([]*Source, *Response, error)
	Get(id int) (*Source, *Response, error)
	Create(source *Source) (*Source, *Response, error)
	Edit(id int, source *Source) (*Source, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewSourcesService(client *Client) SourcesService {
	return &SourcesServiceOp{client}
}

type sourceRoot struct {
	Source *Source `json:"data"`
	Meta   *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type sourcesRoot struct {
	Items []*sourceRoot `json:"items"`
	Meta  *Meta         `json:"meta"`
}

func (r *sourcesRoot) Sources() []*Source {
	contacts := make([]*Source, len(r.Items))
	for i, root := range r.Items {
		contacts[i] = root.Source
	}
	return contacts
}

type SourcesServiceOp struct {
	client *Client
}

func (s *SourcesServiceOp) List(opt *SourceListOptions) ([]*Source, *Response, error) {
	u, err := addOptions("/v2/sources", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(sourcesRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Sources(), res, err
}

func (s *SourcesServiceOp) Get(id int) (*Source, *Response, error) {
	u := fmt.Sprintf("/v2/sources/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(sourceRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Source, res, err
}

func (s *SourcesServiceOp) Create(source *Source) (*Source, *Response, error) {
	u := "/v2/sources"
	envelope := &sourceRoot{Source: source}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(sourceRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Source, res, err
}

func (s *SourcesServiceOp) Edit(id int, source *Source) (*Source, *Response, error) {
	u := fmt.Sprintf("/v2/sources/%d", id)
	envelope := &sourceRoot{Source: source}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(sourceRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Source, res, err
}

func (s *SourcesServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/sources/%d", id)
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
