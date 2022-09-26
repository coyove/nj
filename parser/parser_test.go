package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestHashString(t *testing.T) {
	c, _ := Parse(`b(1, 2+2,  a...) + "a"`, "")
	c.Dump(os.Stdout)

	fmt.Println(ParseJSON("{1:true, a=[1,2]}"))
}
