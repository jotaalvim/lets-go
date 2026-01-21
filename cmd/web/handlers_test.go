package main

import (
	"net/http"
	"testing"

	"modulo.porreiro/internal/assert"
)

func TestPing(t *testing.T) {
	t.Parallel()
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	res := ts.get(t, "/ping")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "ok")
}

type viewTest struct {
	name       string
	urlPath    string
	wantBody   string
	wantStatus int
}

func TestView(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []viewTest{

		{
			name:       "teste com id 1",
			urlPath:    "/view/1",
			wantStatus: http.StatusOK,
			wantBody:   "aaaaaaaa",
		},
		{
			name:       "id inválido",
			urlPath:    "/view/10",
			wantStatus: http.StatusNotFound,
		},

		{
			name:       "id negativo",
			urlPath:    "/view/-1",
			wantStatus: http.StatusNotFound,
		},

		{
			name:       "id negativo",
			urlPath:    "/view/1.23",
			wantStatus: http.StatusNotFound,
		},

		{
			name:       "id negativo",
			urlPath:    "/view/io",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// remove cookies from previous sessions
			ts.resetClientCookieJar(t)

			res := ts.get(t, tt.urlPath)

			assert.Equal(t, res.status, tt.wantStatus)
			assert.StringContains(t, tt.wantBody, res.body)
		})
	}
}
