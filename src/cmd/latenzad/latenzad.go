// vim:set ts=4 sw=4 ai ft=go:
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"protocol"
	"storage"
)

var port *int = flag.Int("p", 7007, "Port on which to listen")

func startServer() chan protocol.Request {
	backend := make(chan protocol.Request)
	go storage.Run(backend)
	return backend
}

func listen(port int) net.Listener {
	s, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to bind to port %d: %s", port, err)
		return nil
	}
	log.Printf("listening on port %d", port)
	return s
}

func accept(lsock net.Listener, backend chan protocol.Request) {
	for {
		if sock, err := lsock.Accept(); err == nil {
			go protocol.HandleConnection(sock, backend)
		} else {
			log.Printf("error accepting connection: %s", err)
		}
	}
}

func main() {
	flag.Parse()
	accept(listen(*port), startServer())
}
