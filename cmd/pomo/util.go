package main

import (
	"fmt"
	"os"
)

func maybe(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
