// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import "fmt"

const (
    HDR_LEN   = 16
    REQ_MAGIC = 0xA5
    RES_MAGIC = 0x5A
)

// request opcodes
const (
    GET    = 0x00
    PUT    = 0x01
    DELETE = 0x02
    STATS  = 0x03
)

// response statuses
const (
    OK          = 0x00
    ENOENT      = 0x01
    E2BIG       = 0x02
    EINVAL      = 0x03
    EBADVERSION = 0x04
    EBADOP      = 0x80
    ENOMEM      = 0x81
)

type LtzRequest struct {
    Opcode    uint8
    Version   uint64
    Key, Body []byte
    Reply     chan LtzResponse
}

type LtzResponse struct {
    Status    uint16
    Version   uint64
    Key, Body []byte
    Fatal     bool
}

func (req LtzRequest) String() string {
    return fmt.Sprintf("{LtzRequest opcode=%x key='%s'}",
        req.Opcode, len(req.Key))
}

func (resp LtzResponse) String() string {
    return fmt.Sprintf("{LtzResponse status=%x version=%x bodylen=%d}",
        resp.Status, resp.Version, len(resp.Body))
}
