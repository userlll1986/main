package main

import (
	"fmt"
	"userlll1986/config"
)

func main() {
	config := config.NewConfig()
	config.ReadConfig()
	fmt.Println(config.Server.Part)
}
