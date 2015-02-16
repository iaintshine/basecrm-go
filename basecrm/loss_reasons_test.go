package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestLossReasonsService(t *testing.T) { TestingT(t) }

type LossReasonsSuite struct {
}

var _ = Suite(&LossReasonsSuite{})

func (s *LossReasonsSuite) TestLossReasonsService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/loss_reasons", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"creator_id": "1",
			"name":       "We were to expensive",
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
              "name": "We were too expensive",
              "created_at": "2014-08-27T16:32:56Z",
              "updated_at": "2014-08-27T16:32:56Z"
          },
          "meta": {
              "type": "loss_reason"
          }
        },
        {
          "data": {
                "id": 2,
                "creator_id": 1,
                "name": "Chosen a competitor",
                "created_at": "2014-08-27T16:32:57Z",
                "updated_at": "2014-08-27T16:32:57Z"
            },
            "meta": {
                "type": "loss_reason"
            }
        }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
            "self": "http://api.getbase.com/v2/loss_reasons.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &LossReasonListOptions{
		1,
		"We were to expensive",
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	lossReasons, res, err := client.LossReasons.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lossReasons, NotNil)

	c.Assert(len(lossReasons), Equals, 2)
	c.Assert(lossReasons[0].Id, Equals, 1)
	c.Assert(lossReasons[1].Id, Equals, 2)
}

func (s *LossReasonsSuite) TestLossReasonsService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/loss_reasons/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
          "id": 1,
          "creator_id": 1,
          "name":  "We were to expensive",
          "created_at": "2014-08-27T16:33:00Z",
          "updated_at": "2014-08-27T16:33:01Z"
        },
        "meta": {
          "type": "loss_reason"
        }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	lossReason, res, err := client.LossReasons.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lossReason, NotNil)

	c.Assert(lossReason.Id, Equals, 1)
}

func (s *LossReasonsSuite) TestLossReasonsService_Create(c *C) {
	setup()
	defer teardown()

	input := &LossReason{
		Name: "Lack of communication",
	}

	expected := &LossReason{
		Id:   1,
		Name: "Lack of communication",
	}

	mux.HandleFunc("/v2/loss_reasons", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(lossReasonRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.LossReason, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Lack of communication"
      },
      "meta": {
        "type": "loss_reason"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	lossReason, res, err := client.LossReasons.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lossReason, NotNil)

	c.Assert(lossReason, DeepEquals, expected)
}

func (s *LossReasonsSuite) TestLossReasonsService_Edit(c *C) {
	setup()
	defer teardown()

	input := &LossReason{
		Name: "Lack of communication with contact",
	}

	expected := &LossReason{
		Id:   1,
		Name: "Lack of communication with contact",
	}

	mux.HandleFunc("/v2/loss_reasons/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(lossReasonRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.LossReason, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Lack of communication with contact"
      },
      "meta": {
        "type": "loss_reason"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	lossReason, res, err := client.LossReasons.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(lossReason, NotNil)

	c.Assert(lossReason, DeepEquals, expected)
}

func (s *LossReasonsSuite) TestLossReasonsService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/loss_reasons/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.LossReasons.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
