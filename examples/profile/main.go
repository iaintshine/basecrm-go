package main

import (
	"fmt"
	"net/url"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/iaintshine/basecrm-go/basecrm"

	. "github.com/iaintshine/basecrm-go/examples/support"
)

const (
	BaseCrmAccessTokenEnv = "BASECRM_TOKEN"
	EnvVariableMissingMsg = `BaseCRM access token has been not found.
  Please set os environment variable 'BASECRM_TOKEN' to one of your
  personal access tokens (PATs). You can register one in the settings.
  `
)

func main() {
	basecrmToken := os.Getenv(BaseCrmAccessTokenEnv)
	if basecrmToken == "" {
		fmt.Printf(EnvVariableMissingMsg)
		os.Exit(1)
	}

	t := oauth.Transport{
		Token: &oauth.Token{
			AccessToken: basecrmToken,
		},
	}

	client := basecrm.NewClient(t.Client())

	company, _, _ := client.Accounts.Self()
	me, _, _ := client.Users.Self()
	team, _, _ := client.Users.List(nil)

	PrintWhoami(me)
	PrintCompany(company)
	PrintTeam(team)
	Flush()
}
