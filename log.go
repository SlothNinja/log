package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	logLevel = "LOGLEVEL"

	LvlNone    = "NONE"
	LvlDebug   = "DEBUG"
	LvlInfo    = "INFO"
	LvlWarning = "WARNING"
	LvlError   = "ERROR"

	debugLabel   = "[DEBUG]"
	infoLabel    = "[INFO]"
	warningLabel = "[WARNING]"
	errorLabel   = "[ERROR]"
)

var DefaultLevel = LvlDebug

func showLogFor(label string) bool {
	switch getLevel() {
	case LvlNone:
		return false
	case LvlDebug:
		return (label == debugLabel) || (label == infoLabel) || (label == warningLabel) || (label == errorLabel)
	case LvlInfo:
		return (label == infoLabel) || (label == warningLabel) || (label == errorLabel)
	case LvlWarning:
		return (label == warningLabel) || (label == errorLabel)
	case LvlError:
		return (label == errorLabel)
	default:
		return true
	}
}

func getLevel() string {
	v, found := os.LookupEnv(logLevel)
	if !found {
		return DefaultLevel
	}
	switch v {
	case LvlNone:
		return LvlNone
	case LvlDebug:
		return LvlDebug
	case LvlInfo:
		return LvlInfo
	case LvlWarning:
		return LvlWarning
	case LvlError:
		return LvlError
	default:
		return DefaultLevel
	}
}

func Debugf(tmpl string, args ...interface{}) {
	if showLogFor(debugLabel) {
		log.Printf(debugLabel+" "+caller()+tmpl, args...)
	}
}

func Infof(tmpl string, args ...interface{}) {
	if showLogFor(infoLabel) {
		log.Printf(infoLabel+" "+caller()+tmpl, args...)
	}
}

func Warningf(tmpl string, args ...interface{}) {
	if showLogFor(warningLabel) {
		log.Printf(warningLabel+" "+caller()+tmpl, args...)
	}
}

func Errorf(tmpl string, args ...interface{}) {
	if showLogFor(errorLabel) {
		log.Printf(errorLabel+" "+caller()+tmpl, args...)
	}
}

func caller() string {
	pc, file, line, _ := runtime.Caller(2)
	files := strings.Split(file, "/")
	if lenFiles := len(files); lenFiles > 1 {
		file = files[lenFiles-1]
	}
	fun := runtime.FuncForPC(pc).Name()
	funs := strings.Split(fun, "/")
	if lenFuns := len(funs); lenFuns > 2 {
		fun = strings.Join(funs[len(funs)-2:], "/")
	}
	return fmt.Sprintf("%v#%v(L: %v)\n\t => ", file, fun, line)
}

func Printf(fmt string, args ...interface{}) {
	log.Printf(caller()+fmt, args...)
}
