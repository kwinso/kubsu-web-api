package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kwinso/kubsu-web-api/internal/server"
)

func main() {
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
