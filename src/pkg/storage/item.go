// vim:set ts=4 sw=4 et ai ft=go:
package storage

import (
    "bytes"
    "protocol"
)

type Item struct {
    fields map[string][]byte
}

type Table struct {
    items map[string]Item
}

func (this *Table) Get(key string) (Item, bool) {
    item, ok := this.items[key]
    return item, ok
}

func (this *Table) Set(key string, item Item) {
    this.items[key] = item
}

func (this *Table) Delete(key string) {
    delete(this.items, key)
}

func NewItem() *Item {
    item := new(Item)
    item.fields = make(map[string][]byte)
    return item
}

func (item *Item) Get(field string) []byte {
    if val, ok := item.fields[field]; ok {
        return val
    }
    return nil
}

func (item *Item) Put(field string, value []byte) {
    item.fields[field] = value
}

func (item *Item) Marshal() ([]byte, error) {
    buf := new(bytes.Buffer)
    for k, v := range item.fields {
        if err := protocol.WriteEntry(buf, []byte(k)); err != nil {
            return nil, err
        }
        if err := protocol.WriteEntry(buf, v); err != nil {
            return nil, err
        }
    }
    return buf.Bytes(), nil
}

func (item *Item) Unmarshal(body []byte) (err error) {
    b := bytes.NewBuffer(body[:])
    i := 0
    for i < len(body) {
        var k, v []byte

        if k, err = protocol.ReadEntry(b); err != nil {
            return
        }
        i += len(k)
        if v, err = protocol.ReadEntry(b); err != nil {
            return
        }
        i += len(v)
        item.fields[string(k)] = v
    }
    return
}
