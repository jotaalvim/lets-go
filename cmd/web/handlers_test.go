package main

import (
	"net/http"
	"net/url"
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

type signupTest struct {
	name         string
	userName     string
	userEmail    string
	userPassword string
	wantStatus   int
	wantFormTag  string
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
			wantBody:   "aaaaaaaaaa",
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

const (
	validName     = "Bob"
	validPassword = "pass56789"
	validEmail    = "bob@gmail.com"
	formTag       = "<form action='/user/signup' method='POST' novalidate>"
)

func TestUserSignUp(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []signupTest{

		{
			name:         "Valid Submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			wantStatus:   http.StatusSeeOther, // redirection
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			wantStatus:   http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			wantStatus:   http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			wantStatus:   http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			wantStatus:   http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "ola tudo bem",
			userPassword: validPassword,
			wantStatus:   http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "1234567",
			wantStatus:   http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, "/user/signup")

			form := url.Values{}
			//t.Logf("form: %v", tt)
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)

			res = ts.postForm(t, "/user/signup", form)

			assert.Equal(t, res.status, tt.wantStatus)
			assert.StringContains(t, tt.wantFormTag, res.body)
			//token := extractCRSRFToken (t, res.body )
			//	t.Logf("CSRF Token: %s",token)

		})
	}

}
