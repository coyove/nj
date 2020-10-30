package parser

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func CatchError(err *error) {
	if r := recover(); r != nil {
		*err, _ = r.(error)
		if *err == nil {
			*err = fmt.Errorf("%v", r)
		}
		if os.Getenv("crab_stack") != "" {
			log.Println(string(debug.Stack()))
		}
	}
}
