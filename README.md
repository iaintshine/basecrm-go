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
