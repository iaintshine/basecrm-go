package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestTagsService(t *testing.T) { TestingT(t) }

type TagsSuite struct {
}

var _ = Suite(&TagsSuite{})

func (s *TagsSuite) TestTagsService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/tags", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"name":          "publisher",
			"creator_id":    "1",
			"resource_type": "lead",
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
            "type": "tag"
          }
      }, {
          "data": {
            "id": 2
          },
          "meta": {
            "type": "tag"
          }
      }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
          "self": "http://api.getbase.com/v2/tags.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &TagListOptions{
		"publisher",
		1,
		LeadResource,
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	tags, res, err := client.Tags.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(tags, NotNil)

	c.Assert(len(tags), Equals, 2)
	c.Assert(tags[0].Id, Equals, 1)
	c.Assert(tags[1].Id, Equals, 2)
}

func (s *TagsSuite) TestTagsService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/tags/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1
      },
      "meta": {
          "type": "tag"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	tag, res, err := client.Tags.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(tag, NotNil)

	c.Assert(tag.Id, Equals, 1)
}

func (s *TagsSuite) TestTagsService_Create(c *C) {
	setup()
	defer teardown()

	input := &Tag{
		ResourceType: LeadResource,
		Name:         "publisher",
	}

	expected := &Tag{
		Id:           1,
		ResourceType: LeadResource,
		Name:         "publisher",
	}

	mux.HandleFunc("/v2/tags", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(tagRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Tag, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "resource_type": "lead",
        "name": "publisher"
      },
      "meta": {
        "type": "tag"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	tag, res, err := client.Tags.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(tag, NotNil)

	c.Assert(tag, DeepEquals, expected)
}

func (s *TagsSuite) TestTagsService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Tag{
		Name: "super important",
	}

	expected := &Tag{
		Id:   1,
		Name: "super important",
	}

	mux.HandleFunc("/v2/tags/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(tagRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Tag, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "super important"
      },
      "meta": {
        "type": "tag"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	tag, res, err := client.Tags.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(tag, NotNil)
	c.Assert(tag, DeepEquals, expected)
}

func (s *TagsSuite) TestTagsService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/tags/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Tags.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
