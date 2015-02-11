package basecrm

import (
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestUsersService(t *testing.T) { TestingT(t) }

type UsersSuite struct {
}

var _ = Suite(&UsersSuite{})

func (s *UsersSuite) TestUsersService_List_All(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
     {
        "items": [{
          "data": {
              "id": 1,
              "name": "Mark Johnson",
              "email": "mark@salesteam.com",
              "status": "active",
              "type": "user",
              "confirmed": true,
              "created_at": "2014-08-27T16:32:56Z",
              "updated_at": "2014-08-27T17:32:56Z"
          },
          "meta": {
              "type": "user"
          }
        },{
          "data": {
              "id": 2,
              "name": "John doe",
              "email": "john.doe@salesteam.com",
              "status": "active",
              "type": "developer",
              "confirmed": true,
              "created_at": "2014-08-27T16:32:56Z",
              "updated_at": "2014-08-27T17:32:56Z"
          },
          "meta": {
              "type": "user"
          }
        }],
        "meta": {
          "type": "collection",
          "count": 1,
          "links": {
              "self": "http://api.getbase.com/v2/users.json"
          }
        }
      }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	users, res, err := client.Users.List(nil)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(users, NotNil)

	c.Assert(len(users), Equals, 2)
	c.Assert(users[0].Id, Equals, 1)
	c.Assert(users[1].Id, Equals, 2)
}

func (s *UsersSuite) TestUsersService_Get(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/users/1", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
          "id": 1,
          "name": "Mark Johnson",
          "email": "mark@salesteam.com",
          "status": "active",
          "type": "user",
          "confirmed": true,
          "created_at": "2014-08-27T16:32:56Z",
          "updated_at": "2014-08-27T17:32:56Z"
      },
      "meta": {
          "type": "user"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	user, res, err := client.Users.Get(1)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(user, NotNil)

	c.Assert(user.Id, Equals, 1)
}

func (s *UsersSuite) TestUsersService_Self(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/users/self", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
          "id": 1,
          "name": "Mark Johnson",
          "email": "mark@salesteam.com",
          "status": "active",
          "type": "user",
          "confirmed": true,
          "created_at": "2014-08-27T16:32:56Z",
          "updated_at": "2014-08-27T17:32:56Z"
      },
      "meta": {
          "type": "user"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	user, res, err := client.Users.Self()
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(user, NotNil)

	c.Assert(user.Id, Equals, 1)
}
