package main

import "github.com/kwinso/kubsu-web-api/server"

func main() {
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
