package main

import (
	"fmt"

	"github.com/WaronLimsakul/gator/internal/config"
)

func main() {
	currentConfig, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	currentConfig.SetUser("Ron")
	currentConfig, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", currentConfig)
	return
}
