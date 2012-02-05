// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import (
    "bufio"
    "errors"
    "fmt"
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
    Opcode           uint8
    Status           uint16
    Table, Key, Body []byte
    Fatal            bool
}

func (r Request) String() string {
    return fmt.Sprintf("{Request op=%x table=%s key=%s bodylen=%d}",
        r.Opcode, string(r.Table), string(r.Key), len(r.Body))
}

func (r Request) ReadFrom(rd *bufio.Reader) (err error) {
    var magic uint8

    if magic, err = rd.ReadByte(); err != nil {
        return err
    }
    if magic != REQ_MAGIC {
        return errors.New(fmt.Sprintf("bad magic: 0x%x", magic))
    }
    if r.Opcode, err = rd.ReadByte(); err != nil {
        return err
    }
    if r.Table, err = ReadEntry(rd); err != nil {
        return err
    }
    if r.Key, err = ReadEntry(rd); err != nil {
        return err
    }
    if r.Body, err = ReadEntry(rd); err != nil {
        return err
    }
    return nil
}

func (r Response) String() string {
    return fmt.Sprintf("{Response status=%x table=%s key=%s bodylen=%d}",
        r.Status, string(r.Table), string(r.Key), len(r.Body))
}

func (r Response) WriteTo(w *bufio.Writer) (err error) {
    if err = w.WriteByte(RES_MAGIC); err != nil {
        return err
    }
    if err = w.WriteByte(r.Opcode); err != nil {
        return err
    }
    if err = writeUint16(*w, r.Status); err != nil {
        return err
    }
    if err = WriteEntry(w, r.Table); err != nil {
        return err
    }
    if err = WriteEntry(w, r.Key); err != nil {
        return err
    }
    if err = WriteEntry(w, r.Body); err != nil {
        return err
    }
    if err = w.Flush(); err != nil {
        return err
    }
    return nil
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
    var req Request

    // parse request message
    if err := req.ReadFrom(bufio.NewReader(sock)); err != nil {
        log.Printf("failed to read request: %s", err)
        runtime.Goexit()
    }
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
        if err := resp.WriteTo(bufio.NewWriter(sock)); err != nil {
            log.Printf("error writing response: %s", err)
            return true
        }
    } else {
        log.Printf("error during processing on %s; hanging up", sock)
    }
    return
}
