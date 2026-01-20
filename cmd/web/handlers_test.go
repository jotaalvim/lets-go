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
