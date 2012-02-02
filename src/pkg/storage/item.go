// vim:set ts=4 sw=4 et ai ft=go:
package storage

type Table struct {
    items map[string]Item
}

type Item struct {
    fields map[string][]byte
}

func (this *Table) GetItem(key string) (Item, bool) {
    item, ok := this.items[key]
    return item, ok
}
