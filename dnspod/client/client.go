package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const DefaultEndpoint = "https://dnsapi.cn/"
const DefaultLang = "cn"
const UserAgent = "terraform-provider-dnspod/1.0.0"

type Config struct {
	HttpClient *http.Client
	Logger     *log.Logger
	Endpoint   string
	LoginToken string
	Lang       string
}

type Client struct {
	httpClient *http.Client
	logger     *log.Logger
	endpoint   string
	loginToken string
	lang       string
}

type Response interface {
	ValidateResponse() error
}

type ResponseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type GeneralResponse struct {
	Status ResponseStatus `json:"status"`
}

func (resp *GeneralResponse) ValidateResponse() error {
	if resp.Status.Code != "1" {
		return &BadStatusCodeError{
			Code:    resp.Status.Code,
			Message: resp.Status.Message,
		}
	}

	return nil
}

func (c Config) Client() (*Client, error) {
	if c.LoginToken == "" {
		return nil, InvalidClientFieldError("LoginToken")
	}

	instance := &Client{
		httpClient: c.HttpClient,
		endpoint:   c.Endpoint,
		loginToken: c.LoginToken,
		logger:     c.Logger,
		lang:       c.Lang,
	}

	if instance.endpoint == "" {
		instance.endpoint = DefaultEndpoint
	}
	if !strings.HasSuffix(instance.endpoint, "/") {
		instance.endpoint = instance.endpoint + "/"
	}

	if instance.httpClient == nil {
		instance.httpClient = http.DefaultClient
	}

	if instance.lang == "" {
		instance.lang = DefaultLang
	}

	return instance, nil
}

// Get calls DNSPod API. It will generate signature and append it automatically.
func (c *Client) Post(action string, params url.Values) (resp *http.Response, err error) {
	if c.logger != nil {
		c.logger.Printf("[DEBUG] Request: %s?%s", action, params.Encode())
	}
	params.Set("login_token", c.loginToken)
	params.Set("format", "json")
	params.Set("lang", c.lang)

	req, err := http.NewRequest("POST", c.endpoint+action, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.httpClient.Do(req)
}

func (c *Client) Call(action string, req interface{}, v Response) error {
	params, err := BuildParams(req)
	if err != nil {
		return err
	}

	resp, err := c.Post(action, params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if c.logger != nil {
		c.logger.Printf("[DEBUG] Response: %s", string(bytes))
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return v.ValidateResponse()
}
