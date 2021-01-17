package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
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

	nodeEnv    = "NODE_ENV"
	production = "production"
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
	if !showLogFor(debugLabel) {
		return
	}
	log.Printf(debugLabel+" "+caller()+tmpl, args...)
}

func Infof(tmpl string, args ...interface{}) {
	if !showLogFor(infoLabel) {
		return
	}
	log.Printf(infoLabel+" "+caller()+tmpl, args...)
}

func Warningf(tmpl string, args ...interface{}) {
	if !showLogFor(warningLabel) {
		return
	}
	log.Printf(warningLabel+" "+caller()+tmpl, args...)
}

func Errorf(tmpl string, args ...interface{}) {
	if !showLogFor(errorLabel) {
		return
	}
	log.Printf(errorLabel+" "+caller()+tmpl, args...)
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

func Panicf(fmt string, args ...interface{}) {
	log.Panicf(fmt, args...)
}

func NewClient(parent string, opts ...option.ClientOption) (*Client, error) {
	if !isProduction() {
		return new(Client), nil
	}
	cl, err := logging.NewClient(context.Background(), parent, opts...)
	return &Client{cl}, err
}

type Client struct {
	*logging.Client
}

type Logger struct {
	logID string
	*logging.Logger
}

func (cl *Client) Logger(logID string, opts ...logging.LoggerOption) *Logger {
	if !isProduction() {
		return new(Logger)
	}
	return &Logger{logID: logID, Logger: cl.Client.Logger(logID, opts...)}
}

func (l *Logger) Debugf(tmpl string, args ...interface{}) {
	if !isProduction() {
		Debugf(tmpl, args...)
		return
	}

	if !showLogFor(debugLabel) {
		return
	}

	if l.Logger == nil {
		Warningf("missing logger")
	}

	l.StandardLogger(logging.Debug).Printf(debugLabel+" "+caller()+tmpl, args...)
}

func (l *Logger) StandardLogger(s logging.Severity) *log.Logger {
	return l.Logger.StandardLogger(s)
}

func (l *Logger) Infof(tmpl string, args ...interface{}) {
	if !isProduction() {
		Infof(tmpl, args...)
		return
	}

	if !showLogFor(infoLabel) {
		return
	}

	if l.Logger == nil {
		Warningf("missing logger")
	}

	l.StandardLogger(logging.Info).Printf(debugLabel+" "+caller()+tmpl, args...)
}

func (l *Logger) Warningf(tmpl string, args ...interface{}) {
	if !isProduction() {
		Warningf(tmpl, args...)
		return
	}

	if !showLogFor(infoLabel) {
		return
	}

	if l.Logger == nil {
		Warningf("missing logger")
	}

	l.StandardLogger(logging.Warning).Printf(debugLabel+" "+caller()+tmpl, args...)
}

func (l *Logger) Errorf(tmpl string, args ...interface{}) {
	if !isProduction() {
		Errorf(tmpl, args...)
		return
	}

	if !showLogFor(infoLabel) {
		return
	}

	if l.Logger == nil {
		Warningf("missing logger")
	}

	l.StandardLogger(logging.Error).Printf(debugLabel+" "+caller()+tmpl, args...)
}

// IsProduction returns true if NODE_ENV environment variable is equal to "production".
// GAE sets NODE_ENV environement to "production" on deployment.
// NODE_ENV can be overridden in app.yaml configuration.
func isProduction() bool {
	return os.Getenv(nodeEnv) == production
}
