package config

import (
	"os"
	"strconv"
)

const (
	envLogLevel               = "WAUTH_LOG_LEVEL"
	envAppName                = "WAUTH_APP_NAME"
	envAppWWWTemplates        = "WAUTH_APP_WWW_TEMPLATES"
	envAppWWWStatic           = "WAUTH_APP_WWW_STATIC"
	envAppWWWTimeout          = "WAUTH_APP_WWW_TIMEOUT"
	envAppURL                 = "WAUTH_APP_URL"
	envGoEnvironment          = "WAUTH_GO_ENVIRONMENT"
	envBindAddress            = "WAUTH_BIND_ADDRESS"
	envAuthGoogleClientID     = "WAUTH_AUTH_GOOGLE_CLIENT_ID"
	envAuthGoogleClientSecret = "WAUTH_AUTH_GOOGLE_CLIENT_SECRET"
	envAuthGoogleRedirectURL  = "WAUTH_AUTH_GOOGLE_REDIRECT_URL"
	envAppSessionName         = "WAUTH_APP_SESSION_NAME"
	envAppSessionKey          = "WAUTH_APP_SESSION_KEY"
	envAppSessionTimeout      = "WAUTH_APP_SESSION_TIMEOUT"
	envSFTPServer             = "WAUTH_SFTP_SERVER"
	envSFTPServerURL          = "WAUTH_SFTP_SERVER_URL"
	envSFTPServerURLUser      = "WAUTH_SFTP_SERVER_URL_USER"
	envSFTPServerURLPass      = "WAUTH_SFTP_SERVER_URL_PASS"
	envSFTPAuthCheckIP        = "WAUTH_SFTP_AUTH_CHECK_IP"
	envAppFilesPath           = "WAUTH_APP_FILES_PATH"

	defLogLevel               = "info"
	defAppName                = "AltoFS"
	defAppWWWTemplates        = "./www/templates/"
	defAppWWWStatic           = "./www/static/"
	defAppWWWTimeout          = 30 // seconds
	defAppURL                 = "http://altofs.org"
	defGoEnvironment          = "dev"
	defBindAddress            = ":80"
	defAuthGoogleClientID     = ""
	defAuthGoogleClientSecret = ""
	defAuthGoogleRedirectURL  = defAppURL + "/auth/google/callback"
	defAppSessionName         = "WAUTH-ALTOFS"
	defAppSessionKey          = "wauth-altofs-key"
	defAppSessionTimeout      = 35 // seconds
	defSFTPServer             = "altofs.org:2322"
	defSFTPServerURL          = "http://127.0.0.1:7070"
	defSFTPServerURLUser      = ""
	defSFTPServerURLPass      = ""
	defSFTPAuthCheckIP        = 1
	defAppFilesPath           = "/home/altofs/storage"

	productionMode = "production"

	AppSessionErrorMessage            = "ErrorMessage"
	AppSessionAuthenticated           = "Authenticated"
	AppSessionAuthGoogleState         = "State"
	AppSessionAuthGoogleID            = "AuthID"
	AppSessionAuthGoogleEmail         = "AuthEmail"
	AppSessionAuthGoogleEmailVerified = "AuthEmailVerified"
	AppSessionAuthGooglePicture       = "AuthPicture"
)

var (
	// LogLevel defines de log level
	LogLevel               = envStr(envLogLevel, defLogLevel)
	AppName                = envStr(envAppName, defAppName)
	AppWWWTemplates        = envStr(envAppWWWTemplates, defAppWWWTemplates)
	AppWWWStatic           = envStr(envAppWWWStatic, defAppWWWStatic)
	AppWWWTimeout          = envInt(envAppWWWTimeout, defAppWWWTimeout)
	AppURL                 = envStr(envAppURL, defAppURL)
	GoEnvironment          = envStr(envGoEnvironment, defGoEnvironment)
	BindAddress            = envStr(envBindAddress, defBindAddress)
	AuthGoogleClientID     = envStr(envAuthGoogleClientID, defAuthGoogleClientID)
	AuthGoogleClientSecret = envStr(envAuthGoogleClientSecret, defAuthGoogleClientSecret)
	AuthGoogleRedirectURL  = envStr(envAuthGoogleRedirectURL, defAuthGoogleRedirectURL)
	AppSessionName         = envStr(envAppSessionName, defAppSessionName)
	AppSessionKey          = envStr(envAppSessionKey, defAppSessionKey)
	AppSessionTimeout      = envInt(envAppSessionTimeout, defAppSessionTimeout)
	SFTPServer             = envStr(envSFTPServer, defSFTPServer)
	SFTPServerURL          = envStr(envSFTPServerURL, defSFTPServerURL)
	SFTPServerURLUser      = envStr(envSFTPServerURLUser, defSFTPServerURLUser)
	SFTPServerURLPass      = envStr(envSFTPServerURLPass, defSFTPServerURLPass)
	SFTPAuthCheckIP        = envInt(envSFTPAuthCheckIP, defSFTPAuthCheckIP)
	AppFilesPath           = envStr(envAppFilesPath, defAppFilesPath)
)

// Init()
func init() {

}

// PStr2Str func
func envStr(pKey string, pDef string) string {
	_val, _ok := os.LookupEnv(pKey)
	if _ok {
		return _val
	}
	return pDef
}

// PInt2Int func
func envInt(pKey string, pDef int) int {
	_val, _ok := os.LookupEnv(pKey)
	if _ok {
		_i, err := strconv.Atoi(_val)
		if err == nil {
			return _i
		}
	}
	return pDef
}

// IsProduction return true if it is on production environment
func IsProduction() bool {
	return GoEnvironment == productionMode
}
