// vim:set ts=4 sw=4 et ai ft=go:
package protocol

const (
    HDR_LEN         = 14
    REQ_MAGIC       = 0xA5
    RES_MAGIC       = 0x5A
    MAX_KEY_LENGTH  = 1024
    MAX_ITEM_LENGTH = 1048576
    MAX_BODY_LENGTH = 1048576
)

// request opcodes
const (
    GET    = 0x01
    SET    = 0x02
    DEL    = 0x03
    CREATE = 0x04
    DROP   = 0x05
    MGET   = 0x06
    INCR   = 0x07
    ADD    = 0x08
    SUB    = 0x09
    CARD   = 0x10
    MEMBER = 0x11
    STATS  = 0x12
)

// response statuses
const (
    OK       = 0x0000
    ENOENT   = 0x0001
    E2BIG    = 0x0002
    EINVAL   = 0x0003
    EBADSIZE = 0x0004
    ENOTABLE = 0x0005
    EBADOP   = 0x0006
)
