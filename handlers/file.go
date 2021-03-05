package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/harriklein/wauth/config"
	"github.com/harriklein/wauth/log"
	"github.com/patrickmn/go-cache"
)

// FileHandler handles the home page
func FileHandler(w http.ResponseWriter, r *http.Request) {
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
	/*
		_tokenStr := ""
		if _session.Values["AuthToken"] != nil {
			_tokenStr = _session.Values["AuthToken"].(string)
		}

		// Get error message to display
		_errorMessage := ""
		if _session.Values[config.AppSessionErrorMessage] != nil {
			_errorMessage = _session.Values[config.AppSessionErrorMessage].(string)
		}
	*/
	_val, _found := cacheUsers.Get(_session.Values[config.AppSessionAuthGoogleEmail].(string))
	if !_found {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	_cacheUser := _val.(*TCacheUser)
	cacheUsers.Set(_session.Values[config.AppSessionAuthGoogleEmail].(string), _cacheUser, cache.DefaultExpiration)

	log.Log.Debugf("Check URL in : \n %v\n", _cacheUser.URLAllowed)

	// ---------------------------------

	_url := config.AppURL + r.URL.Path

	_allowed := false
	for i := range _cacheUser.URLAllowed {
		if _cacheUser.URLAllowed[i] == _url {
			_allowed = true
			break
		}
	}

	if !_allowed {
		log.Log.Errorf("File not allowed: %s %s\n", _cacheUser.Email, _url)
		http.Error(w, "File not allowed", http.StatusForbidden)
		return
	}

	_path := filepath.Join(config.AppFilesPath, strings.TrimPrefix(_url, config.AppURL))
	log.Log.Printf("File View: %s %s\n", _cacheUser.Email, _path)

	// It is necessary because by default it redirects to "/docs/" when request "/docs" or "/docs/index.html"
	http.ServeFile(w, r, _path)
}
