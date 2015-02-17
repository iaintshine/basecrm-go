package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestDealsService(t *testing.T) { TestingT(t) }

type DealsSuite struct {
}

var _ = Suite(&DealsSuite{})

func (s *DealsSuite) TestDealsService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/deals", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"q":              "website",
			"name":           "Website redesign",
			"creator_id":     "1",
			"owner_id":       "1",
			"contact_id":     "1",
			"source_id":      "1",
			"loss_reason_id": "1",
			"hot":            "true",
			"page":           "1",
			"per_page":       "25",
			"ids":            "1,2,3",
			"sort_by":        "name:desc,created_at:asc",
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
            "type": "deal"
          }
      }, {
          "data": {
            "id": 2
          },
          "meta": {
            "type": "deal"
          }
      }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
          "self": "http://api.getbase.com/v2/deals.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &DealListOptions{
		"website",
		"Website redesign",
		1,
		1,
		1,
		1,
		1,
		true,
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	deals, res, err := client.Deals.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deals, NotNil)

	c.Assert(len(deals), Equals, 2)
	c.Assert(deals[0].Id, Equals, 1)
	c.Assert(deals[1].Id, Equals, 2)
}

func (s *DealsSuite) TestDealsService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/deals/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1
      },
      "meta": {
          "type": "deal"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	deal, res, err := client.Deals.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deal, NotNil)

	c.Assert(deal.Id, Equals, 1)
}

func (s *DealsSuite) TestDealsService_Create(c *C) {
	setup()
	defer teardown()

	input := &Deal{
		Name: "Website redesign",
		AssociatedContacts: []*AssociatedContact{
			&AssociatedContact{
				ContactId: 1,
				Role:      "primary",
			},
		},
	}

	expected := &Deal{
		Id:   1,
		Name: "Website redesign",
		AssociatedContacts: []*AssociatedContact{
			&AssociatedContact{
				CreatorId: 1,
				ContactId: 1,
				Role:      "primary",
			},
		},
	}

	mux.HandleFunc("/v2/deals", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(dealRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Deal, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Website redesign",
        "associated_contacts": [{
          "creator_id": 1,
          "contact_id": 1,
          "role": "primary"
        }]
      },
      "meta": {
        "type": "deal"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	deal, res, err := client.Deals.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deal, NotNil)

	c.Assert(deal.AssociatedContacts[0], DeepEquals, expected.AssociatedContacts[0])
	c.Assert(deal, DeepEquals, expected)
}

func (s *DealsSuite) TestDealsService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Deal{
		Name: "Website Redesign for Company X",
	}

	expected := &Deal{
		Id:   1,
		Name: "Website Redesign for Company X",
	}

	mux.HandleFunc("/v2/deals/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(dealRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Deal, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Website Redesign for Company X"
      },
      "meta": {
        "type": "deal"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	deal, res, err := client.Deals.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deal, NotNil)
	c.Assert(deal, DeepEquals, expected)
}

func (s *DealsSuite) TestDealsService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/deals/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Deals.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
