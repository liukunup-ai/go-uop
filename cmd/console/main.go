package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/liukunup/go-uop/internal/console"
)

var (
	addr        string
	openBrowser bool
	devMode     bool
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "HTTP server address")
	flag.BoolVar(&openBrowser, "open", true, "Open browser on start")
	flag.BoolVar(&devMode, "dev", false, "Development mode")
	flag.Parse()

	server, err := console.NewServer(addr, devMode)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🚀 go-uop Console starting on %s\n", addr)
	if openBrowser {
		go func() {
			// 延迟打开浏览器
		}()
	}

	log.Fatal(server.Start())
}
