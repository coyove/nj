package typ

import (
	"fmt"
	"testing"
)

func TestShape(t *testing.T) {
	fmt.Println(Shape.Any().Bool().Array(Shape.Array(Shape.Str()).Repeat()))
}
