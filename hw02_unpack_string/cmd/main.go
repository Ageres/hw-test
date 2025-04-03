package main

import (
	"fmt"

	us "github.com/Ageres/hw-test/hw02_unpack_string"
)

func main() {

	//in := ""
	//in := "3abc"
	//in := `qwe\\\55` // "qwe\55555"
	//in := `qw\e2\\\55` // "qwe\55555"
	in := `qw\e2\\\55a` // "qwe\55555a"

	//---------------------------

	//in := `qwe\\55` // ошибка
	//in := `qwe\\5` // "qwe\\\\\"
	// // "qwe\55555"
	//in := "a4bc2d5e"
	//in := "aaa10b"
	//in :=  `qwe\\5`
	//in := "a33a4bc2d5e"
	//in := "abccd"
	//in := "i0"
	//in := "a🙃0"
	//in := "aa🙃1"
	//in := `qwe\4\5`
	//in :=  `qwe\\5`
	//in2 = `qwe\45`

	out, err := us.Unpack(in)
	fmt.Println("out:", out)
	fmt.Println("err", err)

	/*
		var
		out, err = us.Unpack(in2)
		fmt.Println("out:", out)
		fmt.Println("err", err)
	*/
}
