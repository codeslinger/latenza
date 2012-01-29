// vim:set ts=4 sw=4 et ai ft=go:
package storage

import . "protocol"

import (
    "log"
)

type Item struct {
    Version uint64
    Data    []byte
}

type backing struct {
    data map[string]Item
}

type handler func(req LtzRequest, b *backing) LtzResponse

var handlers = map[uint8]handler{
    GET:    handleGET,
    PUT:    handlePUT,
    DELETE: handleDELETE,
    STATS:  handleSTATS,
}

func Service(ingress chan LtzRequest) {
    var b backing
    b.data = make(map[string]Item)
    for {
        req := <-ingress
        log.Printf("serving request %s", req)
        req.Reply <- dispatch(req, &b)
    }
}

func dispatch(req LtzRequest, b *backing) (rv LtzResponse) {
    if f, ok := handlers[req.Opcode]; ok {
        rv = f(req, b)
    } else {
        rv.Status = EBADOP
    }
    return
}

func handleGET(req LtzRequest, b *backing) (rv LtzResponse) {
    if item, ok := b.data[string(req.Key)]; ok {
        rv.Status = OK
        rv.Body = item.Data
    } else {
        rv.Status = ENOENT
    }
    return
}

func handlePUT(req LtzRequest, b *backing) (rv LtzResponse) {
    var item Item

    item.Data = req.Body
    item.Version = 0
    rv.Status = OK
    rv.Version = item.Version
    b.data[string(req.Key)] = item
    return
}

func handleDELETE(req LtzRequest, b *backing) (rv LtzResponse) {
    delete(b.data, string(req.Key))
    rv.Status = OK
    return
}

func handleSTATS(req LtzRequest, b *backing) (rv LtzResponse) {
    rv.Status = EINVAL
    return
}
