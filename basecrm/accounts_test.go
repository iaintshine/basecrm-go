package basecrm

import (
	"fmt"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func TestAccountsService(t *testing.T) { TestingT(t) }

type AccountsSuite struct {
}

var _ = Suite(&AccountsSuite{})

func (s *AccountsSuite) TestAccountsService_Self(c *C) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/accounts/self", func(w http.ResponseWriter, req *http.Request) {
		c.Assert(req.Method, Equals, "GET")

		w.Header().Add("Content-Type", "application/json")

		jsonBlob := `
    {
      "data": {
        "id": 1,
        "name": "Sales Co",
        "type": "full",
        "plan": "enterprise",
        "currency": "USD",
        "time_format": "12H",
        "timezone": "UTC-05:00",
        "phone": "202-555-0141",
        "created_at": "2014-09-28T16:32:56Z",
        "updated_at": "2014-09-28T16:32:56Z"
      },
      "meta": {
        "type": "account"
      }
    }
    `
		fmt.Fprintf(w, jsonBlob)
	})

	account, res, err := client.Accounts.Self()
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	c.Assert(account, NotNil)

	c.Assert(account.Id, Equals, 1)
}
