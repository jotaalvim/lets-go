package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	app := &application{
		logger: slog.New(slog.DiscardHandler),
	}
	return app
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// initialize a new cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(t)
	}

	// add the jar to the server, so that stores the
	// coockies in the new jar
	ts.Client().Jar = jar

	// prevent the test server client from following redirects by
	// setting a custom CheckRedirect, it tells the client to
	// immediatly return the received response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) resetClientCookieJar(t *testing.T) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar
}

type testResponse struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
	body    string
}

func (ts *testServer) get(t *testing.T, urlPath string) testResponse {

	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		//Marca o teste como FAIL
		t.Fatal(err)
	}

	//envia e executa um request to the server
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}

}
