// vim:set ts=4 sw=4 et ai ft=go:
package storage

type Table struct {
    Items map[string]Item
}

type Item struct {
    Fields map[string][]byte
}
