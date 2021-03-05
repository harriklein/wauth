package app

import (
	"net/http"

	"github.com/harriklein/wauth/config"
	"github.com/harriklein/wauth/handlers"
)

// MapUrls defines the URLs to handle
func MapUrls() {
	srvMux.HandleFunc("/", handlers.HomeHandler)
	srvMux.HandleFunc("/viewer", handlers.ViewerHandler)
	srvMux.HandleFunc("/login", handlers.LoginGetHandler).Methods(http.MethodGet)
	srvMux.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost)
	srvMux.HandleFunc("/logout", handlers.LogoutHandler)
	srvMux.HandleFunc("/auth/google/callback", handlers.AuthGoogleCallbackHandler)

	srvMux.HandleFunc("/api/v1/validate202009", handlers.ValidateHandler)
	//srvMux.HandleFunc("/keepalive", handlers.KeepHandler)

	srvMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.AppWWWStatic))))

	srvMux.PathPrefix("/files/").HandlerFunc(handlers.FileHandler)

	/*
		// region DUMMY HANDLER ------------------------------------
		// Initialize handler and register routers
		_dummyHandler := handlers.NewDummyHandler()
		_dummyRouter := srvMux.PathPrefix("/api/v1/dummies").Subrouter()
		_dummyRouter.HandleFunc("", _dummyHandler.Read).Methods(http.MethodGet)
		_dummyRouter.HandleFunc("", _dummyHandler.CreateOrApplyUpdates).Methods(http.MethodPost)
		_dummyRouter.HandleFunc("", _dummyHandler.Update).Methods(http.MethodPut)
		_dummyRouter.HandleFunc("", _dummyHandler.Delete).Methods(http.MethodDelete)
		_dummyRouter.HandleFunc("/{id}", _dummyHandler.Read).Methods(http.MethodGet)
		_dummyRouter.HandleFunc("/{id}", _dummyHandler.CreateOrApplyUpdates).Methods(http.MethodPost)
		_dummyRouter.HandleFunc("/{id}", _dummyHandler.Update).Methods(http.MethodPut)
		_dummyRouter.HandleFunc("/{id}", _dummyHandler.Delete).Methods(http.MethodDelete)
		// endregion -----------------------------------------------

		// region DOCUMENTATION HANDLER ----------------------------
		_docHandler := handlers.NewDocHandler()
		srvMux.HandleFunc("/docs/{file}", _docHandler.Get).Methods(http.MethodGet)
		srvMux.HandleFunc("/docs/", _docHandler.Get).Methods(http.MethodGet)
		srvMux.HandleFunc("/docs", _docHandler.Get).Methods(http.MethodGet)
		// endregion -----------------------------------------------
	*/
}
