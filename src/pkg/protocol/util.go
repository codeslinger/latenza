// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import (
    "log"
    "io"
    "bufio"
    "net"
    "runtime"
    "encoding/binary"
)

func readBytes(s net.Conn, buf []byte) {
    if _, err := io.ReadFull(s, buf); err != nil {
        log.Printf("error reading message: %s", err)
        runtime.Goexit()
    }
}

func writeBytes(s *bufio.Writer, x []byte) {
    if len(x) > 0 {
        wrote, err := s.Write(x)
        if err != nil || wrote != len(x) {
            log.Printf("error writing to stream: %s", err)
            runtime.Goexit()
        }
    }
}

func writeByte(s *bufio.Writer, x byte) {
    buf := make([]byte, 1)
    buf[0] = x
    writeBytes(s, buf)
}

func writeUint16(s *bufio.Writer, x uint16) {
    buf := []byte{0,0}
    binary.BigEndian.PutUint16(buf, x)
    writeBytes(s, buf)
}

func writeUint32(s *bufio.Writer, x uint32) {
    buf := []byte{0,0,0,0}
    binary.BigEndian.PutUint32(buf, x)
    writeBytes(s, buf)
}

func writeUint64(s *bufio.Writer, x uint64) {
    buf := []byte{0,0,0,0,0,0,0,0}
    binary.BigEndian.PutUint64(buf, x)
    writeBytes(s, buf)
}
