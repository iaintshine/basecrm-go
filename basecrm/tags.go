package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type Tag struct {
	Id           int          `json:"id,omitempty"`
	CreatorId    int          `json:"creator_id,omitempty"`
	ResourceType ResourceType `json:"resource_type,omitempty"`
	Name         string       `json:"name,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
}

type TagListOptions struct {
	Name string `url:"name,omitempty"`

	CreatorId int `url:"creator_id,omitempty"`

	ResourceType ResourceType `url:"resource_type,omitempty"`

	ListOptions
}

type TagsService interface {
	List(opt *TagListOptions) ([]*Tag, *Response, error)
	Get(id int) (*Tag, *Response, error)
	Create(tag *Tag) (*Tag, *Response, error)
	Edit(id int, tag *Tag) (*Tag, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewTagsService(client *Client) TagsService {
	return &TagsServiceOp{client}
}

type tagRoot struct {
	Tag  *Tag `json:"data"`
	Meta *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type tagsRoot struct {
	Items []*tagRoot `json:"items"`
	Meta  *Meta      `json:"meta"`
}

func (r *tagsRoot) Tags() []*Tag {
	tags := make([]*Tag, len(r.Items))
	for i, root := range r.Items {
		tags[i] = root.Tag
	}
	return tags
}

type TagsServiceOp struct {
	client *Client
}

func (s *TagsServiceOp) List(opt *TagListOptions) ([]*Tag, *Response, error) {
	u, err := addOptions("/v2/tags", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagsRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Tags(), res, err
}

func (s *TagsServiceOp) Get(id int) (*Tag, *Response, error) {
	u := fmt.Sprintf("/v2/tags/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Tag, res, err
}

func (s *TagsServiceOp) Create(tag *Tag) (*Tag, *Response, error) {
	u := "/v2/tags"
	envelope := &tagRoot{Tag: tag}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Tag, res, err
}

func (s *TagsServiceOp) Edit(id int, tag *Tag) (*Tag, *Response, error) {
	u := fmt.Sprintf("/v2/tags/%d", id)
	envelope := &tagRoot{Tag: tag}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Tag, res, err
}

func (s *TagsServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/tags/%d", id)
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
