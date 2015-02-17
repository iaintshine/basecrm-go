package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestLeadsService(t *testing.T) { TestingT(t) }

type LeadsSuite struct {
}

var _ = Suite(&LeadsSuite{})

func (s *LeadsSuite) TestLeadsService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/leads", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"q":                    "john",
			"letter":               "J",
			"creator_id":           "1",
			"owner_id":             "1",
			"first_name":           "john",
			"last_name":            "doe",
			"organization_name":    "Design Services Company",
			"status":               "new",
			"address[city]":        "Hyannis",
			"address[postal_code]": "02601",
			"address[country]":     "US",
			"page":                 "1",
			"per_page":             "25",
			"ids":                  "1,2,3",
			"sort_by":              "name:desc,created_at:asc",
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
            "type": "lead"
          }
      }, {
          "data": {
            "id": 2
          },
          "meta": {
            "type": "lead"
          }
      }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
          "self": "http://api.getbase.com/v2/leads.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &LeadListOptions{
		"john",
		"J",
		1,
		1,
		"john",
		"doe",
		"Design Services Company",
		"new",
		"Hyannis",
		"02601",
		"US",
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	leads, res, err := client.Leads.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(leads, NotNil)

	c.Assert(len(leads), Equals, 2)
	c.Assert(leads[0].Id, Equals, 1)
	c.Assert(leads[1].Id, Equals, 2)
}

func (s *LeadsSuite) TestLeadsService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/leads/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1
      },
      "meta": {
          "type": "lead"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	lead, res, err := client.Leads.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lead, NotNil)

	c.Assert(lead.Id, Equals, 1)
}

func (s *LeadsSuite) TestLeadsService_Create(c *C) {
	setup()
	defer teardown()

	input := &Lead{
		FirstName: "Mark",
		LastName:  "Johnson",
	}

	expected := &Lead{
		Id:        1,
		FirstName: "Mark",
		LastName:  "Johnson",
	}

	mux.HandleFunc("/v2/leads", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(leadRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Lead, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "first_name": "Mark",
        "last_name": "Johnson"
      },
      "meta": {
        "type": "lead"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	lead, res, err := client.Leads.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lead, NotNil)

	c.Assert(lead, DeepEquals, expected)
}

func (s *LeadsSuite) TestLeadsService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Lead{
		Description: "I know him via Tom",
	}

	expected := &Lead{
		Id:          1,
		FirstName:   "Mark",
		LastName:    "Johnson",
		Description: "I know him via Tom",
	}

	mux.HandleFunc("/v2/leads/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(leadRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Lead, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "first_name": "Mark",
        "last_name": "Johnson",
        "description": "I know him via Tom"
      },
      "meta": {
        "type": "lead"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	lead, res, err := client.Leads.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lead, NotNil)

	c.Assert(lead, DeepEquals, expected)
}

func (s *LeadsSuite) TestLeadsService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/leads/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Leads.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
