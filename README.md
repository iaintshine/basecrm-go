# basecrm-go

basecrm-go is a Go client library for accessing the BaseCRM API v2 beta.

You can view BaseCRM API docs here: [https://developers.getbase.com/](https://developers.getbase.com/)

## Usage

```go
import "github.com/iaintshine/basecrm-go/basecrm"
```

Create a new BaseCRM client

```go
client := basecrm.NewClient(httpClient)
```

Now use the exposed services to access different parts of the BaseCRM API.

## Authentication

The basecrm-go library does not handle authentication. Instead, you pass an instance of `http.Client`
that is able to handle authentication to resource servers using Bearer authentication schema.
The recommended way to do this is to use [goauth2](https://code.google.com/p/goauth2/), newer 
[golang.org/x/oauth2](https://godoc.org/golang.org/x/oauth2) or any other library that can provide
an instance of `http.Client`. The easiest way to interact with BaseCRM is to use Personal Access Tokens
(PATs), which can be generated in the Accounts Settings page.

Assuming you have set `BASECRM_TOKEN` environment variable to a Personal Access Token.

```go
t := oauth.Transport{
  Token: &oauth.Token{
    AccessToken: os.Getenv("BASECRM_TOKEN"),
  },
}

client := basecrm.NewClient(t.Client())

me, _, _ := client.Users.Self()
```

Full running examples can be found under [examples](https://github.com/iaintshine/basecrm-go/tree/master/examples/) directory.    

## Examples

To create a new Contact:

```go
createRequest := &Contact{
  IsOrganization: false,
  FirstName:      "Mark",
  LastName:       "Johnson",
}

newContact, _, err := client.Contacts.Create(createRequest)

if err != nil {
  fmt.Printf("Something bad happened during the request to the contacts service %v", err)
  return  
}

// do something with the new contact
```
