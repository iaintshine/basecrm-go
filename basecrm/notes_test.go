package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestNotesService(t *testing.T) { TestingT(t) }

type NotesSuite struct {
}

var _ = Suite(&NotesSuite{})

func (s *NotesSuite) TestNotesService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/notes", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"q":             "important",
			"creator_id":    "1",
			"resource_type": "lead",
			"resource_id":   "1",
			"page":          "1",
			"per_page":      "25",
			"ids":           "1,2,3",
			"sort_by":       "name:desc,created_at:asc",
		}
		c.Assert(req, HasHttpMethod, "GET")
		c.Assert(req, HasHttpHeader, "Accept", "application/json")
		c.Assert(req, HasQueryParams, expected)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
     {
      "items": [{
          "data": {
            "id": 1
          },
          "meta": {
            "type": "note"
          }
      }, {
          "data": {
            "id": 2
          },
          "meta": {
            "type": "note"
          }
      }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
          "self": "http://api.getbase.com/v2/notes.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &NoteListOptions{
		"important",
		1,
		LeadResource,
		1,
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	notes, res, err := client.Notes.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(notes, NotNil)

	c.Assert(len(notes), Equals, 2)
	c.Assert(notes[0].Id, Equals, 1)
	c.Assert(notes[1].Id, Equals, 2)
}

func (s *NotesSuite) TestNotesService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/notes/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1
      },
      "meta": {
          "type": "note"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	note, res, err := client.Notes.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(note, NotNil)

	c.Assert(note.Id, Equals, 1)
}

func (s *NotesSuite) TestNotesService_Create(c *C) {
	setup()
	defer teardown()

	input := &Note{
		ResourceType: LeadResource,
		ResourceId:   1,
		Content:      "Highly important.",
	}

	expected := &Note{
		Id:           1,
		ResourceType: LeadResource,
		ResourceId:   1,
		Content:      "Highly important.",
	}

	mux.HandleFunc("/v2/notes", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(noteRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Note, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "resource_type": "lead",
        "resource_id": 1,
        "content": "Highly important."
      },
      "meta": {
        "type": "note"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	note, res, err := client.Notes.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(note, NotNil)

	c.Assert(note, DeepEquals, expected)
}

func (s *NotesSuite) TestNotesService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Note{
		Content: "Highly important. Assign to Tom.",
	}

	expected := &Note{
		Id:      1,
		Content: "Highly important. Assign to Tom.",
	}

	mux.HandleFunc("/v2/notes/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(noteRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Note, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "content": "Highly important. Assign to Tom."
      },
      "meta": {
        "type": "note"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	note, res, err := client.Notes.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(note, NotNil)
	c.Assert(note, DeepEquals, expected)
}

func (s *NotesSuite) TestNotesService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/notes/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Notes.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
