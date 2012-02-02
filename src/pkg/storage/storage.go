// vim:set ts=4 sw=4 et ai ft=go:
package storage

import . "protocol"

import (
    "log"
)

type backing struct {
    data map[string]Table
}

type handler func(req Request, table *Table) Response

var handlers = map[uint8]handler{
    GET: handleGET,
    SET: handleSET,
    DEL: handleDEL,
}

func Run(ingress chan Request) {
    var b backing
    b.data = make(map[string]Table)
    for {
        req := <-ingress
        log.Printf("serving request %s", req)
        req.Reply <- dispatch(req, &b)
    }
}

func dispatch(req Request, b *backing) (rv Response) {
    f, ok := handlers[req.Opcode]
    if !ok {
        rv.Status = EBADOP
        return
    }
    if table, ok2 := b.data[string(req.Table)]; ok2 {
        return f(req, &table)
    }
    rv.Status = ENOTABLE
    return
}

func handleGET(req Request, table *Table) (rv Response) {
    return
}

func handleSET(req Request, table *Table) (rv Response) {
    return
}

func handleDEL(req Request, table *Table) (rv Response) {
    return
}
