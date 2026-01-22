package main

import (
	"context"
	"fmt"
	"net/http"
	//"github.com/justinas/nosurf"
)

func commonHeaders(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)

}

func (app *application) logRequest(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// check if a panic happend
			pv := recover()
			if pv != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// this is so that require authentication routes are not stores in users browser cache
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *application) preventCSRF(next http.Handler) http.Handler {
	//cop := http.NewCrossOriginProtection()

	////cop.AddTrustedOrigin("https://...")
	//cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusBadRequest)
	//	_, err := w.Write([]byte("CSRF check failed"))
	//	if err != nil {
	//		app.logger.Error(err.Error())
	//	}
	//}))

	//return cop.Handler(next)
	return http.Handler(next)
	//csrfHandler := nosurf.New(next)
	//csrfHandler.SetBaseCookie(http.Cookie{
	//    HttpOnly: true,
	//    Path    : "/",
	//    Secure  : true,
	//})
	//return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Retrieve the authenticatedUserID, GetInt returns 0 if does not exists
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		var exists bool

		exists, err := app.users.Exists(id)

		if err != nil {
			app.serverError(w, r, err)
			return
		}
		// an user exists in our database
		if exists {
			// create a copy of the context with the authkey = true
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
