package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	dynamic := aliceNew(app.sessionManager.LoadAndSave, app.authenticate)

	// mux.HandleFunc("/", app.home)
	mux.Handle("/", dynamic.ThenFunc(app.home))
	mux.Handle("/snippet/view", dynamic.ThenFunc(app.snippetView))
	mux.Handle("/user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("/user/signuping", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("/user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("/user/logining", dynamic.ThenFunc(app.userLoginPost))

	portected := dynamic.Append(app.requireAuthentication)

	mux.Handle("/snippet/create", portected.ThenFunc(app.snippetCreate))
	mux.Handle("/snippet/creating", portected.ThenFunc(app.snippetCreatePost))
	mux.Handle("/user/logout", portected.ThenFunc(app.userLogoutPost))

	standard := aliceNew(app.recoverPanic, app.logRequest, secureHeaders)

	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standard.Then(mux)
}
