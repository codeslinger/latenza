// vim:set ts=4 sw=4 et ai ft=go:
package storage

import . "protocol"

import (
    "log"
)

type store struct {
    data map[string]Table
}

type handler func(req Request, table *Table) Response

var handlers = map[uint8]handler{
    GET: handleGET,
    SET: handleSET,
    DEL: handleDEL,
}

func Run(ingress chan Request) {
    var s store
    s.data = make(map[string]Table)
    for {
        req := <-ingress
        log.Printf("serving request %s", req)
        req.Reply <- dispatch(req, &s)
    }
}

func dispatch(req Request, s *store) (rv Response) {
    f, ok := handlers[req.Opcode]
    if !ok {
        rv.Status = EBADOP
        return
    }
    if table, ok2 := s.data[string(req.Table)]; ok2 {
        return f(req, &table)
    }
    rv.Status = ENOTABLE
    return
}

func handleGET(req Request, table *Table) (rv Response) {
    _, ok := table.GetItem(string(req.Key))
    if !ok {
        rv.Status = ENOENT
        return
    }
    rv.Status = OK
    return
}

func handleSET(req Request, table *Table) (rv Response) {
    item, err := UnmarshalItem(req.Body)
    if err != nil {
        rv.Status = EINVAL
        return
    }
    table.SetItem(string(req.Key), *item)
    rv.Status = OK
    return
}

func handleDEL(req Request, table *Table) (rv Response) {
    table.DeleteItem(string(req.Key))
    rv.Status = OK
    return
}
