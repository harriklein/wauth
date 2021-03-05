package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/harriklein/wauth/config"
	"github.com/harriklein/wauth/log"
)

// ValidateHandler gets the login page
func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	_content, _ := ioutil.ReadAll(r.Body)
	log.Log.Debugf("Validate202009 %s -> %s:%s\n           %s\n", r.RemoteAddr, r.Method, r.URL.String(), _content)

	var _contentResponse struct {
		Username string `json:"username"`
		IP       string `json:"ip"`
		Password string `json:"password"`
		Protocol string `json:"protocol"`
	}

	_status := 0

	if _error := json.Unmarshal(_content, &_contentResponse); _error != nil {
		log.Log.Errorf("Could not decode response: %s\n", _error.Error())
		http.Error(w, _error.Error(), http.StatusInternalServerError)
		return
	}

	log.Log.Debugf("  contentResponse: %v\n", _contentResponse)

	_val, _found := cacheUsers.Get(_contentResponse.Username)
	if _found {
		_cacheUser := _val.(*TCacheUser)
		log.Log.Debugf("  Cache found: %s %v \n", _contentResponse.Username, _cacheUser)
		if _contentResponse.Password == _cacheUser.Password {
			if config.SFTPAuthCheckIP > 0 {
				if _contentResponse.IP == _cacheUser.IP {
					_status = 1
				} else {
					log.Log.Infoln("  Invalid IP: %s <> %s \n", _contentResponse.IP, _cacheUser.IP)
				}
			} else {
				_status = 1
			}
		}
	} else {
		log.Log.Debugf("  Cache not found: %s\n", _contentResponse.Username)
	}

	if _status == 1 {
		log.Log.Infoln("Login OK: " + _contentResponse.Username)
	} else {
		log.Log.Infoln("Login Failed: " + _contentResponse.Username)
	}

	fmt.Fprintf(w, "{\"status\":%d}", _status)
}
