package main

import (
	"log"

	"github.com/HAHLIK/AuthService/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log.Println(cfg)
}
