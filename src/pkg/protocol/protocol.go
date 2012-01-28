// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import (
    "log"
    "net"
)

func HandleConnection(s net.Conn, rc chan LtzRequest) {
    log.Print("connection established from %s", s)
    defer hangup(s)
    for handleRequest(s, rc) {
    }
}

func hangup(s net.Conn) {
    log.Print("closing connection from %s", s)
    s.Close()
}

func handleRequest(s net.Conn, rc chan LtzRequest) (rv bool) {
    rv = false
    return
}

