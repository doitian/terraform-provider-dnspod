package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test InvalidClientFieldError
func TestClientConfig(t *testing.T) {
	c := Config{}
	_, err := c.Client()
	if err == nil || !strings.Contains(err.Error(), "LoginToken") {
		t.Error("Expect InvalidClientFieldError on LoginToken")
	}

	c = Config{LoginToken: "1,token"}
	_, err = c.Client()
	if err != nil {
		t.Error("Expect no error but got: ", err)
	}
}

func TestDomainInfo(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()

	c, err := Config{
		Endpoint:   hs.URL,
		LoginToken: "1,test",
	}.Client()
	if err != nil {
		t.Fatal("Error creating client ", err)
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/Domain.Info" {
			t.Error("Bad path!")
		}
		io.WriteString(rw, `{"status":{"code":"1"},"domain":{"id":"12345"}}`)
	}

	req := DomainInfoRequest{DomainId: "12345"}
	var resp DomainInfoResponse

	err = c.Call("Domain.Info", &req, &resp)
	if err != nil {
		t.Error("Got error sending item: ", err)
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/Domain.Info" {
			t.Error("Bad path!")
		}
		io.WriteString(rw, `{"status":{"code":"2","message":"test"}}`)
	}

	err = c.Call("Domain.Info", &req, &resp)
	if err == nil {
		t.Error("Expect BadStatusCodeError but got no exception")
	} else if err.Error() != "Bad StatusCode 2: test" {
		t.Error("Expect BadStatusCodeError but got: ", err)
	}
}
