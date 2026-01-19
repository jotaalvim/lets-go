package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	//mux.Handle("GET  /{$}"      , app.sessionManager.LoadAndSave(app.preventCSRF(app.authenticate(http.HandlerFunc(app.home)))))
	//mux.Handle("GET  /view/{id}", app.sessionManager.LoadAndSave(app.preventCSRF(app.authenticate(http.HandlerFunc(app.view)))))

	//mux.Handle("GET  /user/signup", app.sessionManager.LoadAndSave(app.preventCSRF(app.authenticate(http.HandlerFunc(app.userSignup)))))
	mux.Handle("POST /user/signup", app.sessionManager.LoadAndSave(app.preventCSRF(app.authenticate(http.HandlerFunc(app.userSignupPost)))))
	//mux.Handle("GET  /user/login", app.sessionManager.LoadAndSave(app.preventCSRF(app.authenticate(http.HandlerFunc(app.userLogin)))))
	mux.Handle("POST /user/login", app.sessionManager.LoadAndSave(app.preventCSRF(app.authenticate(http.HandlerFunc(app.userLoginPost)))))

	//mux.Handle("GET  /create"     , app.sessionManager.LoadAndSave(http.HandlerFunc(app.create)))
	//mux.Handle("POST /create"     , app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.createPost))))
	mux.Handle("POST /user/logout", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.userLogoutPost))))

	// onde ponho o meu 404 render?
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))

}
