package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"post-search-available", "/search-ava", "POST", []postData{
		{key: "start", value: "2013-01-01"},
		{key: "end", value: "2013-01-04"},
	}, http.StatusOK},
	{"post-search-available-json", "/search-ava-json", "POST", []postData{
		{key: "first_name", value: "test"},
		{key: "last_name", value: "test"},
		{key: "email", value: "test@test.com"},
		{key: "phone", value: "+2348886645658"},
	}, http.StatusOK},
}

func TestHandler(t *testing.T) {
	routes := getRoute()
	testServer := httptest.NewTLSServer(routes)

	defer testServer.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			res, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for  %s, exxpected %d, got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			res, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for  %s, exxpected %d, got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		}
	}
}
