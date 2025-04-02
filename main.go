package main

import (
	"fmt"
	"go_project/config"
)

func main() {
	config := config.NewConfig()
	config.ReadConfig()
	fmt.Println(config.Server.Part)
}
