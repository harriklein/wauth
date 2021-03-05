package handlers

import (
	"net/http"
	"strings"

	"github.com/harriklein/wauth/log"
)

// StaticGet gets static files
func StaticGet(pResponse http.ResponseWriter, pRequest *http.Request) {

	const _staticPath = "./static/"
	_url := strings.TrimPrefix(pRequest.URL.Path, "/static")

	log.Log.Debugln(_url)

	// It is necessary because by default it redirects to "/static/" when request "/static" or "/static/index.html"
	if pRequest.URL.Path == "/static/" {
		http.ServeFile(pResponse, pRequest, _staticPath+"index.html")
	} else {
		http.ServeFile(pResponse, pRequest, _staticPath+_url)
	}
}
