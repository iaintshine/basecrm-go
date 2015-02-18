package basecrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestTasksService(t *testing.T) { TestingT(t) }

type TasksSuite struct {
}

var _ = Suite(&TasksSuite{})

func (s *TasksSuite) TestTasksService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/tasks", func(w http.ResponseWriter, req *http.Request) {
		expected := map[string]string{
			"q":             "call",
			"creator_id":    "1",
			"owner_id":      "1",
			"resource_type": "lead",
			"resource_id":   "1",
			"completed":     "true",
			"overdue":       "true",
			"remind":        "true",
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
            "type": "task"
          }
      }, {
          "data": {
            "id": 2
          },
          "meta": {
            "type": "task"
          }
      }],
      "meta": {
        "type": "collection",
        "count": 2,
        "links": {
          "self": "http://api.getbase.com/v2/tasks.json"
        }
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	opt := &TaskListOptions{
		"call",
		1,
		1,
		LeadResource,
		1,
		true,
		true,
		true,
		ListOptions{
			Page:    1,
			PerPage: 25,
			Ids:     []int{1, 2, 3},
			SortBy:  []string{"name:desc", "created_at:asc"},
		},
	}
	tasks, res, err := client.Tasks.List(opt)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(tasks, NotNil)

	c.Assert(len(tasks), Equals, 2)
	c.Assert(tasks[0].Id, Equals, 1)
	c.Assert(tasks[1].Id, Equals, 2)
}

func (s *TasksSuite) TestTasksService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/tasks/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1
      },
      "meta": {
          "type": "task"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	task, res, err := client.Tasks.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(task, NotNil)

	c.Assert(task.Id, Equals, 1)
}

func (s *TasksSuite) TestTasksService_Create(c *C) {
	setup()
	defer teardown()

	input := &Task{
		ResourceType: LeadResource,
		ResourceId:   1,
		Content:      "Contact Tom.",
	}

	expected := &Task{
		Id:           1,
		ResourceType: LeadResource,
		ResourceId:   1,
		Content:      "Contact Tom.",
	}

	mux.HandleFunc("/v2/tasks", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "POST")

		root := new(taskRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Task, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "resource_type": "lead",
        "resource_id": 1,
        "content": "Contact Tom."
      },
      "meta": {
        "type": "task"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	task, res, err := client.Tasks.Create(input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(task, NotNil)

	c.Assert(task, DeepEquals, expected)
}

func (s *TasksSuite) TestTasksService_Edit(c *C) {
	setup()
	defer teardown()

	input := &Task{
		Content: "Contact Tom and Rachel.",
	}

	expected := &Task{
		Id:      1,
		Content: "Contact Tom and Rachel.",
	}

	mux.HandleFunc("/v2/tasks/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "PUT")

		root := new(taskRoot)
		json.NewDecoder(req.Body).Decode(root)
		c.Assert(root.Task, NotNil)

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "content": "Contact Tom and Rachel."
      },
      "meta": {
        "type": "task"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	task, res, err := client.Tasks.Edit(1, input)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(task, NotNil)
	c.Assert(task, DeepEquals, expected)
}

func (s *TasksSuite) TestTasksService_Delete(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/tasks/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	deleted, res, err := client.Tasks.Delete(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(deleted, Equals, true)
}
