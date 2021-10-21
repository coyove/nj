package parser

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type CatchedError struct {
	Original interface{}
}

func (e CatchedError) Error() string {
	return fmt.Sprint(e.Original)
}

func IsDebug() bool {
	return os.Getenv("crab_stack") != ""
}

func processSpecialError(err *error, r interface{}) bool {
	if x, ok := r.(interface{ IsValue(Node) }); ok {
		*err = CatchedError{x}
		return true
	}
	if x, ok := r.(interface{ GetRootPanic() interface{} }); ok {
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
