package main

import (
	"fmt"

	"github.com/coyove/script"
)

type Body struct {
	X, Y, Z    float64
	Vx, Vy, Vz float64
	Mass       float64
}

func main() {
	p, err := script.LoadFile("n-body.lua", "body", func() *Body {
		return &Body{}
	})

	fmt.Println(err)
	fmt.Println(p.Run())
}
