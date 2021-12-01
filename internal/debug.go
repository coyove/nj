package internal

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

const UnnamedFunc = "<native>"

type TransparentError struct{}

func Panic(msg string, args ...interface{}) {
	panic(fmt.Errorf(msg, args...))
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsDebug() bool {
	return os.Getenv("njd") != ""
}

func processSpecialError(err *error, r interface{}) bool {
	if x, ok := r.(interface{ TransparentError() TransparentError }); ok {
		*err = x.(error)
		return true
	}
	return false
}

func processPanic(err *error, r interface{}) {
	if IsDebug() {
		log.Println(string(debug.Stack()))
	}

	*err, _ = r.(error)
	if *err == nil {
		*err = fmt.Errorf("%v", r)
	}
}

func CatchError(err *error) {
	if r := recover(); r != nil {
		if processSpecialError(err, r) {
			return
		}
		processPanic(err, r)
	}
}

func CatchErrorFuncCall(err *error, f string) {
	if r := recover(); r != nil {
		if processSpecialError(err, r) {
			return
		}
		processPanic(err, fmt.Errorf("%s() %v", f, r))
	}
}
