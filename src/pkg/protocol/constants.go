// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import "fmt"

const (
    HDR_LEN   = 6
    REQ_MAGIC = 0xA5
    RES_MAGIC = 0x5A
)

// request opcodes
const (
    GET    = 0x00
    PUT    = 0x01
    DELETE = 0x02
)

// response statuses
const (
    OK      = 0x0000
    ENOENT  = 0x0001
    E2BIG   = 0x0002
    EINVAL  = 0x0003
    EBADOP  = 0x0080
)

type LtzRequest struct {
    Opcode  uint8
    Body    []byte
    Reply   chan LtzResponse
}

type LtzResponse struct {
    Status  uint16
    Body    []byte
    Fatal   bool
}

func (req LtzRequest) String() string {
    return fmt.Sprintf("{LtzRequest opcode=%x bodylen=%d}",
        req.Opcode, len(req.Body))
}

func (resp LtzResponse) String() string {
    return fmt.Sprintf("{LtzResponse status=%x version=%x bodylen=%d}",
        resp.Status, len(resp.Body))
}
