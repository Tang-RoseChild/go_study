package main

import (
	"fmt"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	fmt.Println(" pwd : ", wd)

	if err := os.Chdir(" .."); err != nil {
		fmt.Println("err when chdir : ", err.Error())
		return
	}
	wd, _ = os.Getwd()
	fmt.Println("after change dir, wd is : ", wd)

}
