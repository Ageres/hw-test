package main

import (
	"fmt"

	us "github.com/Ageres/hw-test/hw02_unpack_string"
)

func main() {
	//in := "a4bc2d5e"
	//in := "a33a4bc2d5e"
	//in := "abccd"
	//in := "i0"
	//in := "a🙃0"
	//in := "aa🙃1"
	in := `qwe\4\5`

	out, err := us.Unpack(in)
	fmt.Println("out:", out)
	fmt.Println("err", err)

	/*
		var in2 = `qwe\45`
		out, err = us.Unpack(in2)
		fmt.Println("out:", out)
		fmt.Println("err", err)
	*/
}
