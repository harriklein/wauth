package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/harriklein/wauth/config"
	"github.com/harriklein/wauth/log"
	"github.com/patrickmn/go-cache"
)

// TViewerPageData struct
type TViewerPageData struct {
	AppName        string
	ErrorMsg       string
	Email          string
	LaunchURL      string
	Token          string
	RefreshTimeout int
	DataSource     string
}

var (
	viewerPage *template.Template
)

// TFileNode list
type TFileNode struct {
	Name     string      `json:"title"`
	Path     string      `json:"-"`
	URL      string      `json:"link"`
	IsDir    bool        `json:"folder"`
	SubNodes []TFileNode `json:"children"`
}

//   init() initializes the auth elements
func init() {
	viewerPage = template.Must(template.ParseFiles(config.AppWWWTemplates + "viewer.html"))
}

// ViewerHandler handles the home page
func ViewerHandler(w http.ResponseWriter, r *http.Request) {
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

	// ---------------------------------

	var _URLAllowed []string
	var _fileNodes []TFileNode

	log.Log.Debugf("  VF: %v\n", _cacheUser.User.VirtualFolders)

	for _, _virtualFolder := range _cacheUser.User.VirtualFolders {
		_fileNodes = append(_fileNodes, TFileNode{
			Name:  strings.ReplaceAll(_virtualFolder.VirtualPath, "/", ""),
			IsDir: true,
			Path:  _virtualFolder.MappedPath,
			URL:   strings.Replace(_virtualFolder.MappedPath, config.AppFilesPath, config.AppURL, -1)})
	}

	log.Log.Debugf("  FN: %v\n", _fileNodes)

	for _n := range _fileNodes {
		_fileNodes[_n].addSubNodes(&_URLAllowed)
	}

	log.Log.Debugf("  FN2: %v\n", _fileNodes)

	log.Log.Debugf("\n\n  URLAllowed: \n%v\n\n\n", _URLAllowed)

	_cacheUser.URLAllowed = _URLAllowed
	cacheUsers.Set(_session.Values[config.AppSessionAuthGoogleEmail].(string), _cacheUser, cache.DefaultExpiration)

	_jsonFileNodes, _error := json.Marshal(_fileNodes)
	if _error != nil {
		log.Log.Errorf("Folder Structure Error Marshal: %s", _error)
	}

	_srcData := string(_jsonFileNodes)

	log.Log.Debugf("  srcData: %s\n", _srcData)

	// Set information and load page
	_viewerPageData := &TViewerPageData{
		AppName:        config.AppName,
		ErrorMsg:       _errorMessage,
		Email:          _session.Values[config.AppSessionAuthGoogleEmail].(string),
		LaunchURL:      fmt.Sprintf("sftp://%s:%s@%s", _cacheUser.Email, _cacheUser.Password, config.SFTPServer),
		Token:          _tokenStr,
		RefreshTimeout: (config.AppWWWTimeout),
		DataSource:     _srcData,
	}
	viewerPage.Execute(w, _viewerPageData)
}

func (oNode *TFileNode) addSubNodes(pURLAllowed *[]string) error {
	_files, _error := ioutil.ReadDir(oNode.Path)
	if _error != nil {
		return _error
	}
	for _, _file := range _files {
		_subNode := TFileNode{
			Name: _file.Name(),
		}
		_subNode.Path = filepath.Join(oNode.Path, _subNode.Name)
		_subNode.URL = strings.Replace(_subNode.Path, config.AppFilesPath, config.AppURL, -1)

		if !_file.IsDir() {
			_subNode.IsDir = false
			*pURLAllowed = append(*pURLAllowed, _subNode.URL)
		} else {
			_subNode.IsDir = true
			_subNode.addSubNodes(pURLAllowed)
		}
		oNode.SubNodes = append(oNode.SubNodes, _subNode)
	}
	return nil
}
