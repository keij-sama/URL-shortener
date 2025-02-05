package main

import (
	"UrlShort/internal/config"
	"fmt"
)

func main() {

	//init config
	cfg := config.MustLoad()

	fmt.Println(cfg)
	//init logger

	//init storage

	//init router

	//run server
}
