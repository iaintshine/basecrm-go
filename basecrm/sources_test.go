package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestSourcesService(t *testing.T) { TestingT(t) }

type SourcesSuite struct {
}

var _ = Suite(&SourcesSuite{})

func (s *SourcesSuite) TestSourcesService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/sources", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"creator_id": "1",
			"name":       "Word of mouth",
			"page":       "1",
			"per_page":   "25",
			"ids":        "1,2,3",
			"sort_by":    "name:desc,created_at:asc",
		}
		c.Assert(req, HasHttpMethod, "GET")
		c.Assert(req, HasHttpHeader, "Accept", "application/json")
		c.Assert(req, HasQueryParams, expected)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
     {
      "items": [{
          "data": {
              "id": 1,
              "creator_id": 1,
              "name": "Our website",
              "created_at": "2014-08-27T16:32:56Z",
              "updated_at": "2014-08-27T16:32:56Z"
          },
          "meta": {
              "type": "source"
          }
        },
        {
          "data": {
                "id": 2,
                "creator_id": 1,
                "name": "Word of mouth",
                "created_at": "2014-08-27T16:32:57Z",
                "updated_at": "2014-08-27T16:32:57Z"
            },
            "meta": {
                "type": "source"
            }
        }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
            "self": "http://api.getbase.com/v2/sources.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &SourceListOptions{
		1,
		"Word of mouth",
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	sources, res, err := client.Sources.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(sources, NotNil)

	c.Assert(len(sources), Equals, 2)
	c.Assert(sources[0].Id, Equals, 1)
	c.Assert(sources[1].Id, Equals, 2)
}

func (s *SourcesSuite) TestSourcesService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/sources/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
          "id": 5,
          "creator_id": 1,
          "name":  "Tom referral",
          "created_at": "2014-08-27T16:33:00Z",
          "updated_at": "2014-08-27T16:33:01Z"
        },
        "meta": {
          "type": "source"
        }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	source, res, err := client.Sources.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(source, NotNil)

	c.Assert(source.Id, Equals, 5)
}

func (s *SourcesSuite) TestSourcesService_Create(c *C) {
	setup()
	defer teardown()

	input := &Source{
		Name: "Tom",
	}

	expected := &Source{
		Id:   1,
		Name: "Tom",
	}

	mux.HandleFunc("/v2/sources", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(sourceRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Source, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Tom"
      },
      "meta": {
        "type": "source"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	source, res, err := client.Sources.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(source, NotNil)

	c.Assert(source, DeepEquals, expected)
}

func (s *SourcesSuite) TestSourcesService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Source{
		Name: "Tom referral",
	}

	expected := &Source{
		Id:   1,
		Name: "Tom referral",
	}

	mux.HandleFunc("/v2/sources/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(sourceRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Source, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Tom referral"
      },
      "meta": {
        "type": "source"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	source, res, err := client.Sources.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(source, NotNil)

	c.Assert(source, DeepEquals, expected)
}

func (s *SourcesSuite) TestSourcesService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/sources/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Sources.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
