// vim:set ts=4 sw=4 et ai ft=go:
package protocol

import (
    "encoding/binary"
    "errors"
    "io"
)

func ReadEntry(r io.Reader) ([]byte, error) {
    slice := []byte{0, 0, 0, 0}
    n, err := r.Read(slice[0:4])
    if err != nil {
        return nil, err
    }
    if n != 4 {
        return nil, errors.New("could not read length from buffer")
    }
    size := binary.BigEndian.Uint32(slice[0:4])
    if size == 0 {
        return nil, nil
    }
    val := make([]byte, size)
    if _, err := io.ReadFull(r, val); err != nil {
        return nil, err
    }
    return val, nil
}

func WriteEntry(out io.Writer, data []byte) error {
    size := []byte{0, 0, 0, 0}
    if data == nil {
        if _, err := out.Write(size); err != nil {
            return err
        }
    } else {
        binary.BigEndian.PutUint32(size, uint32(len(data)))
        if _, err := out.Write(size); err != nil {
            return err
        }
        if _, err := out.Write(data); err != nil {
            return err
        }
    }
    return nil
}

func writeUint16(w io.Writer, x uint16) error {
    buf := []byte{0, 0}
    binary.BigEndian.PutUint16(buf, x)
    n, err := w.Write(buf)
    if err != nil || n != 2 {
        return err
    }
    return nil
}
