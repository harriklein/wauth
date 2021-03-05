package log

import (
	"os"
	"strings"

	"github.com/harriklein/wauth/config"
	"github.com/sirupsen/logrus"
)

// Log is the global variable for the customized log
var Log *logrus.Logger

// Init initializes log
func Init() {
	_level, _error := logrus.ParseLevel(config.LogLevel)
	if _error != nil {
		_level = logrus.DebugLevel
	}

	Log = &logrus.Logger{
		Level: _level,
		Out:   os.Stdout,
	}

	if config.IsProduction() {
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Formatter = &logrus.TextFormatter{}
	}
}

// ParseFields attempts to convert a string to a field. eg.: "f:v" to {"f","v"}
func ParseFields(pTags ...string) logrus.Fields {
	_result := make(logrus.Fields, len(pTags))
	for _, _tag := range pTags {
		_els := strings.Split(_tag, ":")
		_result[strings.TrimSpace(_els[0])] = strings.TrimSpace(_els[1])
	}
	return _result
}
