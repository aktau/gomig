package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Welcome to gomig v.%v.%v\n", GOMIG_MAJ_VERSION, GOMIG_MIN_VERSION)

	conf, err := LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("sucesfully loaded config: %v\n", conf)
}
