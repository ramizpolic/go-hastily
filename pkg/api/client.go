package api

// This interface is used to facilitate and extend HTTP Client methods
// when talking to backend APIs.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	cfg "github.com/fhivemind/go-hastily/config"
	"github.com/fhivemind/go-hastily/pkg/auth"
	"github.com/fhivemind/go-hastily/pkg/common"
	. "github.com/fhivemind/go-hastily/pkg/global"
)

// config loads environment configuration
var envCfg = cfg.LoadConfig()

// Client defines wrapper structure of http client.
type Client struct {
	Auth     *auth.Credentials
	Endpoint string
	Model    string
	Instance *http.Client
}

// Response generalizes http request results.
type Response struct {
	Success    bool
	Message    string
	StatusCode int
}

// Request generalizes http request form.
type Request struct {
	URI         string
	Id          string
	Path        string
	Query       map[string]string
	Body        interface{}
	requestType string
}

// ResponseList holds values of statuses for a specific
// key-value mapping.
type ResponseList struct {
	Data map[string]*Response
}

// ApiClient exposes methods designed to
// interact with backend for http client.
type ApiClient interface {
	// shared
	CheckConnection() Response
	DefaultResponse(string, error) Response
	// http requests
	Get(Request, interface{}) Response
	Put(Request, interface{}) Response
	Post(Request, interface{}) Response
	Delete(Request, interface{}) Response
}

// NewClient creates a new http ApiClient
// to interact with API.
func NewClient(model string) *Client {

	// make default
	client := Client{
		Auth:     &auth.Credentials{},
		Endpoint: envCfg.ApiEndpoint,
		Instance: &http.Client{},
		Model:    model,
	}

	// check if auth needed
	if envCfg.LoginEndpoint != "" {

		// load credentials
		creds, err := auth.LoadCredentials()
		HandleError(err)

		// update client
		client.Auth = creds

		// verify client
		resp := client.CheckConnection()
		if !resp.Success {
			HandleError(errors.New(resp.Message))
		}
	}

	return &client
}

// CheckConnection verifies the validity of endpoints and credentials for backend APIs.
func (client *Client) CheckConnection() Response {

	// request form
	request := Request{
		URI: envCfg.VerifyEndpoint,
	}

	// user model
	type User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	// check current credentials
	var user []*User
	resp := client.request(request, &user)
	if !resp.Success || len(user) == 0 || user[0].Email == "" {
		return client.DefaultResponse("Invalid or expired credentials. Please login again.", errors.New(resp.Message))
	}

	return client.DefaultResponse("", nil)
}

// Get controls GET requests on backend APIs.
func (client *Client) Get(request Request, object interface{}) Response {
	request.requestType = http.MethodGet
	return client.request(request, object)
}

// Put controls PUT requests on backend APIs.
func (client *Client) Put(request Request, object interface{}) Response {
	request.requestType = http.MethodPut
	return client.request(request, object)
}

// Post controls POST requests on backend APIs.
func (client *Client) Post(request Request, object interface{}) Response {
	request.requestType = http.MethodPost
	return client.request(request, object)
}

// Delete controls DELETE requests on backend APIs.
func (client *Client) Delete(request Request, object interface{}) Response {
	request.requestType = http.MethodDelete
	return client.request(request, object)
}

// DefaultResponse returns Response object based on error.
func (client *Client) DefaultResponse(str string, err error) Response {
	if err != nil {
		msg := fmt.Sprintf("%+v", err)
		if str != "" {
			msg = str
		}

		return Response{
			Success:    false,
			StatusCode: 0,
			Message:    msg,
		}
	}

	return Response{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    str,
	}
}

// reguest is private generic function of http REST API methods.
func (client *Client) request(request Request, object interface{}) Response {

	// request params
	var reqBody io.Reader
	if request.Body != nil {
		json, err := json.Marshal(request.Body)
		if err != nil {
			return client.DefaultResponse("", err)
		}
		reqBody = bytes.NewBuffer(json)
	}

	// create request
	endpoint := client.getEndpointForRequest(request)
	req, err := http.NewRequest(request.requestType, endpoint, reqBody)
	if err != nil {
		return client.DefaultResponse("", err)
	}

	// set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", client.Auth.AccessToken))

	// send request
	resp, err := client.Instance.Do(req)
	if err != nil {
		return client.DefaultResponse("", err)
	}
	defer resp.Body.Close()

	// check if not 200
	if resp.StatusCode != http.StatusOK {
		return Response{
			Success:    false,
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Non-OK HTTP Status code"),
		}
	}

	// save result?
	if object != nil {
		err = json.NewDecoder(resp.Body).Decode(object)
		if err != nil {
			client.DefaultResponse("", err)
		}
	}

	return client.DefaultResponse("", nil)
}

// getEndpointForRequest shared function to create API URI for given
// client model and uri query parameters.
// e.g. {http://facebook.com} / {v2/users} / {1} ? {arg1=val1} & {arg2=val2}
func (client *Client) getEndpointForRequest(request Request) string {

	// base url
	apiPath := request.URI
	if apiPath == "" {
		apiPath = client.getEndpointForId(request.Id)
		if request.Path != "" {
			apiPath = client.getEndpointForPath(request.Path)
		}
	}
	u, _ := url.Parse(apiPath)

	// add query args
	queryString := u.Query()
	for k, v := range request.Query {
		queryString.Set(k, v)
	}
	u.RawQuery = queryString.Encode()

	return u.String()
}

// getEndpointForPath shared function to create API URI for given
// optional path. e.g. {http://facebook.com} / {v2/users}
func (client *Client) getEndpointForPath(path string) string {
	return fmt.Sprintf("%s/%s", client.Endpoint, path)
}

// getEndpointForId shared function to create API URI for given
// client model and requested Id.
// e.g. {http://facebook.com} / {v2/users} / {1}
func (client *Client) getEndpointForId(id string) string {
	return fmt.Sprintf("%s/%s/%s", client.Endpoint, client.Model, id)
}

// NewResponseList initializes a new list.
func NewResponseList() *ResponseList {
	return &ResponseList{
		Data: make(map[string]*Response),
	}
}

// Get returns a key from response list.
func (list *ResponseList) Get(key string) (*Response, bool) {
	val, ok := list.Data[key]
	return val, ok
}

// HasKey checks if key is inside response list.
func (list *ResponseList) HasKey(key string) bool {
	_, ok := list.Data[key]
	return ok
}

// Insert adds a new key, value pair to response list.
func (list *ResponseList) Insert(key string, resp *Response) {
	list.Data[key] = resp
}

// Size returns the length of list.
func (list *ResponseList) Size() int {
	return len(list.Data)
}

// Successes counts number of valid responses inside list.
func (list *ResponseList) Successes() int {
	success := 0
	for _, v := range list.Data {
		if v.Success {
			success++
		}
	}
	return success
}

// ToGeneric converts list to map of generic objects.
func (list *ResponseList) ToGeneric() map[string]*common.Generic {
	data := make(map[string]*common.Generic)
	for key, value := range list.Data {
		data[key] = common.ObjectToGeneric(value)
	}
	return data
}
