// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import (
    "log"
    "io"
    "net"
    "bufio"
    "runtime"
    "encoding/binary"
)

func HandleConnection(s net.Conn, server chan LtzRequest) {
    log.Print("connection established from %s", s)
    defer hangup(s)
    for handleRequest(s, server) {
    }
}

func hangup(s net.Conn) {
    if s != nil {
        log.Print("closing connection from %s", s)
        s.Close()
    }
}

func handleRequest(sock net.Conn, server chan LtzRequest) (rv bool) {
    // initialize and read request header
    hdrBytes := make([]byte, HDR_LEN)
    bytesRead, err := io.ReadFull(sock, hdrBytes)
    if err != nil || bytesRead != HDR_LEN {
        log.Print("error reading message from %s: %s (%d bytes read)",
            sock, err, bytesRead)
        return
    }

    // read and parse request message
    req := parseHeader(hdrBytes)
    readBytes(sock, req.Key)
    readBytes(sock, req.Body)
    log.Printf("processing request %s", req)

    // send newly-formed request message to server for processing
    req.Reply = make(chan LtzResponse, 1)
    server <- req

    // get response from server and dispatch it or die if fatal error 
    // occurred
    resp := <-req.Reply
    if rv = !resp.Fatal; rv {
        log.Printf("got response %s", resp)
        sendResponse(sock, req, resp)
    } else {
        log.Printf("error during processing on %s; hanging up", sock)
    }
    return
}

func parseHeader(hdrBytes []byte) (rv LtzRequest) {
    if hdrBytes[0] != REQ_MAGIC {
        log.Printf("bad magic: %x", hdrBytes[0])
        runtime.Goexit()
    }
    rv.Opcode = hdrBytes[1]
    rv.Key = make([]byte, binary.BigEndian.Uint16(hdrBytes[2:]))
    bodyLen := binary.BigEndian.Uint32(hdrBytes[4:]) - uint32(len(rv.Key))
    rv.Version = binary.BigEndian.Uint64(hdrBytes[8:])
    rv.Body = make([]byte, bodyLen)
    return
}

func sendResponse(sock net.Conn, req LtzRequest, resp LtzResponse) {
    out := bufio.NewWriter(sock)
    writeByte(out, RES_MAGIC)
    writeByte(out, req.Opcode)
    writeUint16(out, resp.Status)
    writeUint64(out, resp.Version)
    writeUint32(out, uint32(len(resp.Key)))
    writeUint32(out, uint32(len(resp.Body)))
    writeBytes(out, resp.Key)
    writeBytes(out, resp.Body)
    out.Flush()
}
