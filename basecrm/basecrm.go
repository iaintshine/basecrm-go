package basecrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.1.0"
	defaultBaseURL = "https://api.getbase.com"
	userAgent      = "basecrm-go/" + libraryVersion

	defaultMediaType = "application/json"
)

// A client manages communication with the API
type Client struct {
	// HTTP client to communicate with the API
	client *http.Client

	// Base URL for the API requests. Defaults to the public Base API
	// at https://api.getbase.com/ but can be set to sandbox environment.
	BaseURL *url.URL

	// User Agent for the API Client
	UserAgent string

	// Services used to communicating with the API.
	Accounts    AccountsService
	Users       UsersService
	Contacts    ContactsService
	Sources     SourcesService
	LossReasons LossReasonsService
	Leads       LeadsService
	Deals       DealsService
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`

	// A comma-separated list of IDs to be returned in the request.
	Ids []int `url:"ids,comma,omitempty"`

	// A comma-separated list of fields to sort by.
	SortBy []string `url:"sort_by,comma,omitempty"`
}

// NewClient returns a new instance of the Base API v2 client.
// If no client is provided, default http client is used instead.
//
// To use API methods which require authentication, provide an http.Client
// that will perform the authentication for you.
// Use the goauth2 library.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.Accounts = NewAccountsService(c)
	c.Users = NewUsersService(c)
	c.Contacts = NewContactsService(c)
	c.Sources = NewSourcesService(c)
	c.LossReasons = NewLossReasonsService(c)
	c.Leads = NewLeadsService(c)
	c.Deals = NewDealsService(c)

	return c
}

// Response is a Base API response. This wraps the standard http.Response
// returned from the Base and provides convenient access to thinks like pagination
// and meta data.
type Response struct {
	*http.Response

	Meta *Meta
}

// An ErrorResponse reports one or more errors caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// More details on individual errors
	Errors *ErrorsEnvelope
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors.String())
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if buf != nil {
		req.Header.Add("Content-Type", defaultMediaType)
	}

	req.Header.Add("Accept", defaultMediaType)

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	if err = checkResponse(resp); err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return response, err
			}
		}
	}

	return response, err
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qv, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qv.Encode()
	return u.String(), nil
}

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)

	if err == nil && data != nil {
		errorResponse.Errors = &ErrorsEnvelope{}
		json.Unmarshal(data, errorResponse.Errors)
	}
	return errorResponse
}
