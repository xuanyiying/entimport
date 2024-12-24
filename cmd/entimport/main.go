package main

import (
	"flag"
)

var (
	addr = flag.String("addr", ":8080", "Web server address")
)

/*
func main() {
	flag.Parse()

	server, _ := web.NewServer(&web.Config{})
	log.Printf("Starting web server at %s", *addr)
	if err := server.Start(*addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}*/
