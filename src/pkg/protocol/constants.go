// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import "fmt"

type OperationCode int

const (
    GET    = 0x00
    PUT    = 0x01
    DELETE = 0x02
    STATS  = 0x03
)

type StatusCode int

const (
    OK     = 0x00
    ENOENT = 0x01
    E2BIG  = 0x02
    EINVAL = 0x03
    EBADOP = 0x80
    ENOMEM = 0x81
)

type LtzRequest struct {
    Opcode          uint8
    Version         uint64
    Key, Body       []byte
    ResponseChannel chan LtzResponse
}

type LtzResponse struct {
    Status    uint16
    Version   uint64
    Key, Body []byte
}

func (req LtzRequest) String() string {
    return fmt.Sprintf("{LtzRequest opcode=%x key='%s'}",
        req.Opcode, len(req.Key))
}

func (resp LtzResponse) String() string {
    return fmt.Sprintf("{LtzResponse status=%x version=%x bodylen=%d}",
        resp.Status, resp.Version, len(resp.Body))
}
