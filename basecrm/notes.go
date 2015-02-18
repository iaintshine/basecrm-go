package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type Note struct {
	Id           int          `json:"id,omitempty"`
	CreatorId    int          `json:"creator_id,omitempty"`
	ResourceType ResourceType `json:"resource_type,omitempty"`
	ResourceId   int          `json:"resource_id,omitempty"`
	Content      string       `json:"content,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
}

type NoteListOptions struct {
	Q string `url:"q,omitempty"`

	CreatorId int `url:"creator_id,omitempty"`

	ResourceType ResourceType `url:"resource_type,omitempty"`
	ResourceId   int          `url:"resource_id,omitempty"`

	ListOptions
}

type NotesService interface {
	List(opt *NoteListOptions) ([]*Note, *Response, error)
	Get(id int) (*Note, *Response, error)
	Create(note *Note) (*Note, *Response, error)
	Edit(id int, note *Note) (*Note, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewNotesService(client *Client) NotesService {
	return &NotesServiceOp{client}
}

type noteRoot struct {
	Note *Note `json:"data"`
	Meta *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type notesRoot struct {
	Items []*noteRoot `json:"items"`
	Meta  *Meta       `json:"meta"`
}

func (r *notesRoot) Notes() []*Note {
	notes := make([]*Note, len(r.Items))
	for i, root := range r.Items {
		notes[i] = root.Note
	}
	return notes
}

type NotesServiceOp struct {
	client *Client
}

func (s *NotesServiceOp) List(opt *NoteListOptions) ([]*Note, *Response, error) {
	u, err := addOptions("/v2/notes", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(notesRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Notes(), res, err
}

func (s *NotesServiceOp) Get(id int) (*Note, *Response, error) {
	u := fmt.Sprintf("/v2/notes/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(noteRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Note, res, err
}

func (s *NotesServiceOp) Create(note *Note) (*Note, *Response, error) {
	u := "/v2/notes"
	envelope := &noteRoot{Note: note}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(noteRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Note, res, err
}

func (s *NotesServiceOp) Edit(id int, note *Note) (*Note, *Response, error) {
	u := fmt.Sprintf("/v2/notes/%d", id)
	envelope := &noteRoot{Note: note}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(noteRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Note, res, err
}

func (s *NotesServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/notes/%d", id)
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
