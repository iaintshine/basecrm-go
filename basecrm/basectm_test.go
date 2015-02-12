package basecrm

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	. "gopkg.in/check.v1"
)

var (
	mux *http.ServeMux

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

var HasHttpMethod = &hasHttpMethodChecker{
	&CheckerInfo{Name: "HasHttpMethod", Params: []string{"request", "method"}},
}

var HasHttpStatus = &hasHttpStatusChecker{
	&CheckerInfo{Name: "HasHttpStatus", Params: []string{"response", "status"}},
}

var HasHttpHeader = &hasHttpHeaderChecker{
	&CheckerInfo{Name: "HasHttpHeader", Params: []string{"request|response", "header_name", "header_value"}},
}

var HasQueryParams = &hasQueryParamsChecker{
	&CheckerInfo{Name: "HasQueryParams", Params: []string{"request", "query_params"}},
}

type hasHttpMethodChecker struct {
	*CheckerInfo
}

func (checker *hasHttpMethodChecker) Check(params []interface{}, names []string) (result bool, error string) {
	req, ok := params[0].(*http.Request)
	if !ok {
		return false, "request must be an instance of the http.Request type"
	}
	method, ok := params[1].(string)
	if !ok {
		return false, "method must be a string"
	}
	return req.Method == method, ""
}

type hasHttpStatusChecker struct {
	*CheckerInfo
}

func (checker *hasHttpStatusChecker) Check(params []interface{}, names []string) (result bool, error string) {
	res, ok := params[0].(*http.Response)
	if !ok {
		return false, "response must be an instance of the http.Response type"
	}
	status, ok := params[1].(int)
	if !ok {
		return false, "status must be an int"
	}
	return res.StatusCode == status, ""
}

type hasHttpHeaderChecker struct {
	*CheckerInfo
}

func (checker *hasHttpHeaderChecker) Check(params []interface{}, names []string) (result bool, error string) {
	if len(params) != 3 {
		return false, "invalid number of parameters"
	}

	var headers http.Header
	req, ok := params[0].(*http.Request)
	if ok {
		headers = req.Header
	} else {
		res, ok := params[0].(*http.Response)
		if !ok {
			return false, "first parameter must be either an instance of the http.Response or http.Request"
		}
		headers = res.Header
	}
	header, ok := params[1].(string)
	if !ok {
		return false, "header_key must be a string"
	}
	value, ok := params[2].(string)
	if !ok {
		return false, "header_value must be a string"
	}

	return headers.Get(header) == value, ""
}

type hasQueryParamsChecker struct {
	*CheckerInfo
}

func (checker *hasQueryParamsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	req, ok := params[0].(*http.Request)
	if !ok {
		return false, "request must be an instance of the http.Request type"
	}
	values, ok := params[1].(map[string]string)
	if !ok {
		return false, "request_params must be a map[string]string"
	}

	expected := url.Values{}
	for k, v := range values {
		expected.Add(k, v)
	}

	req.ParseForm()
	got := req.Form

	return reflect.DeepEqual(got, expected), ""
}
