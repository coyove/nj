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

func CatchError(err *error) {
	if r := recover(); r != nil {
		if os.Getenv("crab_stack") != "" {
			log.Println(string(debug.Stack()))
		}

		x, ok := r.(interface{ IsValue(Node) })
		if ok {
			*err = CatchedError{x}
			return
		}

		*err, _ = r.(error)
		if *err == nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
