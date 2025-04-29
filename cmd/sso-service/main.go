package main

import (
	"fmt"
	"github.com/LavaJover/shvark-sso-service/internal/config"
)

func main(){
	cfg := config.MustLoad()
	fmt.Println(cfg)
}