package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	rnd "math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/harriklein/wauth/config"
	"github.com/harriklein/wauth/log"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// TLoginPageData struct
type TLoginPageData struct {
	AppName  string
	ErrorMsg string
}

// TVirtualFolder struct
type TVirtualFolder struct {
	ID          int    `json:"id"`
	MappedPath  string `json:"mapped_path"`
	VirtualPath string `json:"virtual_path"`
}

// TUser struct
type TUser struct {
	ID             int              `json:"id"`
	Username       string           `json:"username"`
	Status         int              `json:"status"`
	HomeDir        string           `json:"home_dir"`
	VirtualFolders []TVirtualFolder `json:"virtual_folders"`
}

// TCacheUser struct
type TCacheUser struct {
	Email      string
	IP         string
	Password   string
	User       TUser
	URLAllowed []string
}

var (
	loginPage        *template.Template
	store            *sessions.CookieStore
	configAuthGoogle *oauth2.Config
	cacheUsers       *cache.Cache
)

//   init() initializes the auth elements
func init() {
	loginPage = template.Must(template.ParseFiles(config.AppWWWTemplates + "login.html"))

	//key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	//key = []byte("super-secret-key")

	_key := []byte(config.AppSessionKey)
	store = sessions.NewCookieStore(_key)
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: config.AppSessionTimeout,
	}
	store.MaxAge(store.Options.MaxAge)
	cacheUsers = cache.New(
		(time.Duration)(config.AppSessionTimeout)*time.Second,
		(time.Duration)(config.AppSessionTimeout)*time.Second)

	configAuthGoogle = &oauth2.Config{
		ClientID:     config.AuthGoogleClientID,
		ClientSecret: config.AuthGoogleClientSecret,
		RedirectURL:  config.AuthGoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

// randToken() generates random 32 bytes token
func getRandomToken32() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// LoginGetHandler gets the login page
func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	// Get session instance
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	// Get error message to display
	_errorMessage := ""
	if _session.Values[config.AppSessionErrorMessage] != nil {
		_errorMessage = _session.Values[config.AppSessionErrorMessage].(string)
	}

	// Revoke users authentication
	_session.Values[config.AppSessionErrorMessage] = ""
	_session.Values[config.AppSessionAuthenticated] = false
	_session.Save(r, w)

	// Set information and load page
	_loginPageData := &TLoginPageData{
		AppName:  config.AppName,
		ErrorMsg: _errorMessage,
	}
	loginPage.Execute(w, _loginPageData)
}

// LoginHandler post the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	var op oauth2.AuthCodeOption = oauth2.SetAuthURLParam("prompt", "select_account")

	_state := getRandomToken32()
	_urlAuth := configAuthGoogle.AuthCodeURL(_state, op)
	// 	log.Log.Debugln(_urlAuth)

	// Update session info
	_session.Values[config.AppSessionAuthGoogleState] = _state
	_session.Save(r, w)
	// Redirect page to google login page
	http.Redirect(w, r, _urlAuth, http.StatusSeeOther)
}

// LogoutHandler handles the logout page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	_tokenStr := r.URL.Query().Get("s")

	// Revoke users authentication
	_session.Options.MaxAge = -1
	_session.Values[config.AppSessionAuthenticated] = false
	_session.Save(r, w)

	_resp, _error := http.PostForm("https://accounts.google.com/o/oauth2/revoke?token="+_tokenStr, url.Values{"key": {"token"}, "id": {_tokenStr}})
	if _error != nil {
		log.Log.Errorf("Could not revoke token: %s\n", _error.Error())
	}
	defer _resp.Body.Close()

	http.Redirect(w, r, "https://www.google.com/accounts/Logout?continue=https://appengine.google.com/_ah/logout?continue="+url.PathEscape(config.AppURL), http.StatusFound)
}

// AuthGoogleCallbackHandler handles the home page
func AuthGoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	// log.Log.Debugln("URL: " + r.Method + " " + r.URL.String())

	if r.FormValue("state") != _session.Values[config.AppSessionAuthGoogleState] {
		if _session.Values[config.AppSessionAuthGoogleState] == nil {
			log.Log.Errorf("State: form(%s) <> session(%s)\n", r.FormValue("state"), "nil")
		} else {
			log.Log.Errorf("State: form(%s) <> session(%s)\n", r.FormValue("state"), _session.Values[config.AppSessionAuthGoogleState].(string))
		}
		_session.Values[config.AppSessionErrorMessage] = "Invalid Session"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	_token, _error := configAuthGoogle.Exchange(oauth2.NoContext, r.FormValue("code"))
	if _error != nil {
		log.Log.Errorf("Could not get token: %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Invalid Token"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	_resp, _error := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + _token.AccessToken)
	if _error != nil {
		log.Log.Errorf("Could not create get request: %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Invalid Token Access"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer _resp.Body.Close()
	_content, _error := ioutil.ReadAll(_resp.Body)
	if _error != nil {
		log.Log.Errorf("Could not parse response: %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Invalid Token Response"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//log.Log.Debugf("Auth Google Response : %s\n", _content)

	var _contentResponse struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"verified_email"`
		Picture       string `json:"picture"`
	}

	if _error := json.Unmarshal(_content, &_contentResponse); _error != nil {
		log.Log.Errorf("Could not decode response: %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Invalid Content Response"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Log.Debugf("LOGIN AutoFS\n  ID   : %s\n  Email: %s\n  Verif: %t\n  Pic  : %s\n",
		_contentResponse.ID, _contentResponse.Email, _contentResponse.EmailVerified, _contentResponse.Picture)

	_session.Values[config.AppSessionAuthGoogleID] = _contentResponse.ID
	_session.Values[config.AppSessionAuthGoogleEmail] = _contentResponse.Email
	_session.Values[config.AppSessionAuthGoogleEmailVerified] = _contentResponse.EmailVerified
	_session.Values[config.AppSessionAuthGooglePicture] = _contentResponse.Picture
	_accountExists := false

	var _respServ *http.Response

	// TODO: VERIFY FTP ACCOUNT
	if config.SFTPServerURLUser == "" {
		_respServ, _error = http.Get(config.SFTPServerURL + "/api/v1/user?username=" + _contentResponse.Email)
	} else {
		_req, _ := http.NewRequest("GET", config.SFTPServerURL+"/api/v1/user?username="+_contentResponse.Email, nil)
		_req.SetBasicAuth(config.SFTPServerURLUser, config.SFTPServerURLPass)
		_client := &http.Client{}
		_respServ, _error = _client.Do(_req)
	}

	if _error != nil {
		log.Log.Errorf("Could not create get request: %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Internal Error: Server"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer _respServ.Body.Close()
	_contentServ, _error := ioutil.ReadAll(_respServ.Body)
	if _error != nil {
		log.Log.Errorf("Could not parse response: %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Invalid Token Response"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var _contentUsers []TUser

	if _error := json.Unmarshal(_contentServ, &_contentUsers); _error != nil {
		log.Log.Errorf("Could not decode response (user): %s\n", _error.Error())
		_session.Values[config.AppSessionErrorMessage] = "Invalid Content Response (user)"
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(_contentUsers) > 0 {
		log.Log.Debugf("User Info: %v\n", _contentUsers)

		if _contentUsers[0].Status == 0 {
			log.Log.Errorf("Account is inactive %s\n", _contentResponse.Email)
			_session.Values[config.AppSessionErrorMessage] = "The account '" + _contentResponse.Email + "' has been blocked. Please contact the administrator."
			_session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		_accountExists = true
	}

	if _accountExists == false {
		log.Log.Errorf("Account does not exist: %s\n", _contentResponse.Email)
		_session.Values[config.AppSessionErrorMessage] = "The account '" + _contentResponse.Email + "' does not exist in " + config.AppName + ". Please contact the administrator."
		_session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	_session.Values["AuthToken"] = _token.AccessToken
	_session.Values[config.AppSessionAuthenticated] = true
	_session.Values[config.AppSessionErrorMessage] = ""
	_session.Save(r, w)

	log.Log.Infof("Account logged: %s\n", _contentResponse.Email)

	_ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	_pwd := rnd.Intn(999999-100000) + 100000
	_cacheUser := &TCacheUser{
		Email:    _contentResponse.Email,
		IP:       _ip,
		Password: strconv.Itoa(_pwd),
		User:     _contentUsers[0],
	}

	log.Log.Debugf("Set Cache: %v\n", _cacheUser)
	cacheUsers.Set(_contentResponse.Email, _cacheUser, cache.DefaultExpiration)

	// Print secret message
	http.Redirect(w, r, "/", http.StatusFound)
}

// secretHandler handles simple demo
func secretHandler(w http.ResponseWriter, r *http.Request) {
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is authenticated
	if _auth, _ok := _session.Values[config.AppSessionAuthenticated].(bool); !_ok || !_auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

// KeepHandler handles simple demo
func KeepHandler(w http.ResponseWriter, r *http.Request) {
	_session, _error := store.Get(r, config.AppSessionName)
	if _error != nil {
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	log.Log.Debugln(_session.Values[config.AppSessionAuthGoogleEmail])

	// Print secret message
	fmt.Fprintln(w, "ok")
}
