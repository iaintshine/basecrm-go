package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestContactsService(t *testing.T) { TestingT(t) }

type ContactsSuite struct {
}

var _ = Suite(&ContactsSuite{})

func (s *ContactsSuite) TestContactsService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/contacts", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"q":                    "john",
			"letter":               "J",
			"creator_id":           "1",
			"owner_id":             "1",
			"is_organization":      "true",
			"name":                 "john",
			"first_name":           "john",
			"last_name":            "doe",
			"customer_status":      "none",
			"prospect_status":      "none",
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
            "type": "contact"
          }
      }, {
          "data": {
            "id": 2
          },
          "meta": {
            "type": "contact"
          }
      }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
          "self": "http://api.getbase.com/v2/contacts.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &ContactListOptions{
		"john",
		"J",
		1,
		1,
		true,
		"john",
		"john",
		"doe",
		"none",
		"none",
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
	contacts, res, err := client.Contacts.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(contacts, NotNil)

	c.Assert(len(contacts), Equals, 2)
	c.Assert(contacts[0].Id, Equals, 1)
	c.Assert(contacts[1].Id, Equals, 2)
}

func (s *ContactsSuite) TestContactsService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/contacts/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1
      },
      "meta": {
          "type": "contact"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	contact, res, err := client.Contacts.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(contact, NotNil)

	c.Assert(contact.Id, Equals, 1)
}

func (s *ContactsSuite) TestContactsService_Create(c *C) {
	setup()
	defer teardown()

	input := &Contact{
		IsOrganization: false,
		FirstName:      "Mark",
		LastName:       "Johnson",
	}

	expected := &Contact{
		Id:             1,
		IsOrganization: false,
		FirstName:      "Mark",
		LastName:       "Johnson",
	}

	mux.HandleFunc("/v2/contacts", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(contactRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Contact, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "is_organization": false,
        "first_name": "Mark",
        "last_name": "Johnson"
      },
      "meta": {
        "type": "contact"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	contact, res, err := client.Contacts.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(contact, NotNil)

	c.Assert(contact, DeepEquals, expected)
}

func (s *ContactsSuite) TestContactsService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Contact{
		Description: "I know him via Tom",
	}

	expected := &Contact{
		Id:             1,
		IsOrganization: false,
		FirstName:      "Mark",
		LastName:       "Johnson",
		Description:    "I know him via Tom",
	}

	mux.HandleFunc("/v2/contacts/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(contactRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Contact, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "is_organization": false,
        "first_name": "Mark",
        "last_name": "Johnson",
        "description": "I know him via Tom"
      },
      "meta": {
        "type": "contact"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	contact, res, err := client.Contacts.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(contact, NotNil)

	c.Assert(contact, DeepEquals, expected)
}

func (s *ContactsSuite) TestContactsService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/contacts/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Contacts.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
