package main

import (
	"flag"

	"22dojo-online/pkg/infrastructure/server"
)

var (
	// Listenするアドレス+ポート
	addr string
)

// https://www.spinute.org/go-by-example/command-line-flags.html
func init() {
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
}

func main() {
	server.Serve(addr)
}
