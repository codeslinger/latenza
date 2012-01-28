// vim:set ts=4 sw=4 et ai ft=go:
package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "protocol"
)

var port *int = flag.Int("p", 7777, "Port on which to listen")

func acceptLoop(l net.Listener) {
    reqChannel := make(chan protocol.LtzRequest)
    log.Printf("listening on port %d", *port)
    for {
        if s, e := l.Accept(); e == nil {
            go protocol.HandleConnection(s, reqChannel)
        } else {
            log.Printf("error accepting from %s", l)
        }
    }
}

func main() {
    flag.Parse()
    s, e := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if e != nil {
        log.Fatalf("failed to bind to port %d: %s", *port, e)
    }
    acceptLoop(s)
}
