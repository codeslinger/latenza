// vim:set ts=4 sw=4 ai ft=go:
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
	item, ok := table.Get(string(req.Key))
	if !ok {
		rv.Status = ENOENT
		return
	}
	body, err := item.Marshal()
	if err != nil {
		rv.Status = ENOMEM
		return
	}
	rv.Body = body
	rv.Table = req.Table
	rv.Key = req.Key
	rv.Status = OK
	return
}

func handleSET(req Request, table *Table) (rv Response) {
	var item Item

	if err := item.Unmarshal(req.Body); err != nil {
		rv.Status = EINVAL
		return
	}
	table.Set(string(req.Key), item)
	rv.Table = req.Table
	rv.Key = req.Key
	rv.Status = OK
	return
}

func handleDEL(req Request, table *Table) (rv Response) {
	table.Delete(string(req.Key))
	rv.Status = OK
	return
}
