package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

// Handle OAuth routes
func handleOAuth(r *mux.Router) {
	manager := manage.NewDefaultManager()

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	r.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

// Called when the user is redirected to /authorize
func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		return "", err
	}

	// Get the user from the session.
	uid, ok := store.Get("AuthorizeUserID")
	if !ok {
		redirectToLogin(w, r, store)

		// The calling function will just return when userID is empty.
		return "", nil
	}

	// The user gave authorization, make sure it can't be reused.
	userID = uid.(string)
	store.Delete("AuthorizeUserID")
	store.Save()

	return userID, nil
}

// Store the redirect data in the session and send the client to the login page
func redirectToLogin(w http.ResponseWriter, r *http.Request, store session.Store) {
	// Make sure the form data is populated.
	if r.Form == nil {
		r.ParseForm()
	}

	// Store the redirect data in the session.
	store.Set("RedirectData", r.Form)
	store.Save()

	// Redirect to the login page.
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}