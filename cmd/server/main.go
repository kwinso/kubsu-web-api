package main

import (
	_ "github.com/joho/godotenv/autoload"
	server "github.com/kwinso/kubsu-web-api"
)

func main() {
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
