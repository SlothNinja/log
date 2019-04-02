package log

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Debugf(tmpl string, args ...interface{}) {
	log.Printf("[DEBUG] "+caller()+tmpl, args...)
}

// func Debugf(ctx context.Context, tmpl string, args ...interface{}) {
// 	// l(ctx).Debugf(caller()+tmpl, args...)
// 	log.Printf("[DEBUG] "+caller()+tmpl, args...)
// }

func Errorf(tmpl string, args ...interface{}) {
	log.Printf("[ERROR] "+caller()+tmpl, args...)
}

// func Errorf(ctx context.Context, tmpl string, args ...interface{}) {
// 	// l(ctx).Errorf(caller()+tmpl, args...)
// 	log.Printf("[ERROR] "+caller()+tmpl, args...)
// }

func Infof(tmpl string, args ...interface{}) {
	log.Printf("[INFO] "+caller()+tmpl, args...)
}

// func Infof(ctx context.Context, tmpl string, args ...interface{}) {
// 	// l(ctx).Infof(caller()+tmpl, args...)
// 	log.Printf("[INFO] "+caller()+tmpl, args...)
// }

func Warningf(tmpl string, args ...interface{}) {
	log.Printf("[WARNING] "+caller()+tmpl, args...)
}

// func Warningf(ctx context.Context, tmpl string, args ...interface{}) {
// 	// l(ctx).Warningf(caller()+tmpl, args...)
// 	log.Printf("[WARNING] "+caller()+tmpl, args...)
// }

// func l(ctx context.Context) logging.Logger {
// 	return logging.Get(ctx)
// }

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
