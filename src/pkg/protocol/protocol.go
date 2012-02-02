// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "log"
    "net"
    "runtime"
)

type Request struct {
    Opcode           uint8
    Table, Key, Body []byte
    Reply            chan Response
}

type Response struct {
    Status           uint16
    Table, Key, Body []byte
    Fatal            bool
}

func (req Request) String() string {
    return fmt.Sprintf("{Request op=%x table=%s key=%s bodylen=%d}",
        req.Opcode, string(req.Table), string(req.Key), len(req.Body))
}

func (resp Response) String() string {
    return fmt.Sprintf("{Response status=%x table=%s key=%s bodylen=%d}",
        resp.Status, string(resp.Table), string(resp.Key), len(resp.Body))
}

func HandleConnection(sock net.Conn, backend chan Request) {
    log.Printf("connection established from %s", sock.RemoteAddr())
    defer hangup(sock)
    for handleRequest(sock, backend) {
    }
}

func hangup(sock net.Conn) {
    if sock != nil {
        log.Printf("closing connection from %s", sock.RemoteAddr())
        sock.Close()
    }
}

func handleRequest(sock net.Conn, backend chan Request) (rv bool) {
    // initialize and read request header
    hdrBytes := make([]byte, HDR_LEN)
    bytesRead, err := io.ReadFull(sock, hdrBytes)
    if err != nil || bytesRead != HDR_LEN {
        log.Printf("error reading message from %s: %s (%d bytes read)",
            sock.RemoteAddr(), err, bytesRead)
        return
    }

    // read and parse request message
    req := parseHeader(hdrBytes)
    readBytes(sock, req.Body)
    log.Printf("processing request %s", req)

    // send newly-formed request message to server for processing
    req.Reply = make(chan Response, 1)
    backend <- req

    // get response from server and dispatch it or die if fatal error 
    // occurred
    resp := <-req.Reply
    rv = !resp.Fatal
    if rv {
        log.Printf("got response %s", resp)
        sendResponse(sock, req, resp)
    } else {
        log.Printf("error during processing on %s; hanging up", sock)
    }
    return
}

func parseHeader(hdrBytes []byte) (rv Request) {
    if hdrBytes[0] != REQ_MAGIC {
        log.Printf("bad magic: 0x%x", hdrBytes[0])
        runtime.Goexit()
    }
    rv.Opcode = hdrBytes[1]
    rv.Body = make([]byte, binary.BigEndian.Uint32(hdrBytes[2:]))
    return
}

func sendResponse(sock net.Conn, req Request, resp Response) {
    out := bufio.NewWriter(sock)
    writeByte(out, RES_MAGIC)
    writeByte(out, req.Opcode)
    writeUint16(out, resp.Status)
    writeUint32(out, uint32(len(resp.Body)))
    writeBytes(out, resp.Body)
    out.Flush()
}
