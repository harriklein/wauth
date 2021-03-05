package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/harriklein/wauth/config"
	"github.com/patrickmn/go-cache"
)

// THomePageData struct
type THomePageData struct {
	AppName        string
	ErrorMsg       string
	Email          string
	LaunchURL      string
	Token          string
	RefreshTimeout int
}

var (
	homePage *template.Template
)

//   init() initializes the auth elements
func init() {
	homePage = template.Must(template.ParseFiles(config.AppWWWTemplates + "index.html"))
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}
	_session.Save(r, w)

	// Check if user is authenticated
	if _auth, _ok := _session.Values[config.AppSessionAuthenticated].(bool); !_ok || !_auth {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	_tokenStr := ""
	if _session.Values["AuthToken"] != nil {
		_tokenStr = _session.Values["AuthToken"].(string)
	}

	// Get error message to display
	_errorMessage := ""
	if _session.Values[config.AppSessionErrorMessage] != nil {
		_errorMessage = _session.Values[config.AppSessionErrorMessage].(string)
	}

	_val, _found := cacheUsers.Get(_session.Values[config.AppSessionAuthGoogleEmail].(string))
	if !_found {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	_cacheUser := _val.(*TCacheUser)
	cacheUsers.Set(_session.Values[config.AppSessionAuthGoogleEmail].(string), _cacheUser, cache.DefaultExpiration)

	// Set information and load page
	_homePageData := &THomePageData{
		AppName:        config.AppName,
		ErrorMsg:       _errorMessage,
		Email:          _session.Values[config.AppSessionAuthGoogleEmail].(string),
		LaunchURL:      fmt.Sprintf("sftp://%s:%s@%s", _cacheUser.Email, _cacheUser.Password, config.SFTPServer),
		Token:          _tokenStr,
		RefreshTimeout: (config.AppWWWTimeout),
	}
	homePage.Execute(w, _homePageData)
}
