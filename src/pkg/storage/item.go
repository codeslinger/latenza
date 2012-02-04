// vim:set ts=4 sw=4 et ai ft=go:
package storage

import (
    "bytes"
    "encoding/binary"
    "errors"
    "io"
)

type Item struct {
    fields map[string][]byte
    size   uint32
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
    item.size = 0
    return item
}

func (item *Item) Get(field string) []byte {
    if val, ok := item.fields[field]; ok {
        return val
    }
    return nil
}

func (item *Item) Marshal() ([]byte, error) {
    buf := new(bytes.Buffer)
    for k, v := range item.fields {
        err := binary.Write(buf, binary.BigEndian, uint32(len(k)))
        if err != nil {
            return nil, err
        }
        buf.Write([]byte(k))
        err = binary.Write(buf, binary.BigEndian, uint32(len(v)))
        if err != nil {
            return nil, err
        }
        buf.Write(v)
    }
    return buf.Bytes(), nil
}

func (item *Item) Unmarshal(body []byte) error {
    b := bytes.NewBuffer(body[:])
    for item.size < uint32(len(body)) {
        k, err := readEntry(b)
        if err != nil {
            return err
        }
        v, err := readEntry(b)
        if err != nil {
            return err
        }
        item.size += uint32(len(k) + len(v))
        item.fields[string(k)] = v
    }
    return nil
}

func readEntry(b *bytes.Buffer) ([]byte, error) {
    slice := make([]byte, 4)
    n, err := b.Read(slice[0:4])
    if err != nil {
        return nil, err
    }
    if n != 4 {
        return nil, errors.New("could not read length from buffer")
    }
    size := binary.BigEndian.Uint32(slice[0:4])
    val := make([]byte, size)
    if _, err := io.ReadFull(b, val); err != nil {
        return nil, err
    }
    return val, nil
}
