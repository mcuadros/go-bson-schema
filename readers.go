package bsonschema

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const x00 = byte(0)

func readInt32(r io.Reader, i *int32) error {
	return binary.Read(r, binary.LittleEndian, i)
}

func readInt64(r io.Reader, i *int64) error {
	return binary.Read(r, binary.LittleEndian, i)
}

func readFloat64(r io.Reader, f *float64) error {
	return binary.Read(r, binary.LittleEndian, f)
}

func readString(r io.Reader, s *string) error {
	var b []byte
	if err := readBytes(r, &b); err != nil {
		return err
	}

	l := len(b)
	if b[l-1] != x00 {
		return errors.New("non-null terminated")
	}

	*s = string(b[:l-1])
	return nil
}

func readBinary(r io.Reader, b *[]byte) error {
	var l int32
	fmt.Println("len", l)
	if err := readInt32(r, &l); err != nil {
		return err
	}

	var k byte
	if err := readByte(r, &k); err != nil {
		return err
	}

	d := make([]byte, l)
	if _, err := r.Read(d); err != nil {
		return err
	}

	*b = d

	return nil
}

func readBytes(r io.Reader, b *[]byte) error {
	var l int32
	if err := readInt32(r, &l); err != nil {
		return err
	}

	d := make([]byte, l)
	if _, err := r.Read(d); err != nil {
		return err
	}

	*b = d

	return nil
}

func readByte(r io.Reader, b *byte) error {
	var n [1]byte
	if _, err := io.ReadFull(r, n[:]); err != nil {
		return err
	}

	*b = n[0]
	return nil
}

func readCString(r io.Reader, s *[]byte) error {
	var b []byte
	var n [1]byte
	for {
		if _, err := io.ReadFull(r, n[:]); err != nil {
			return err
		}

		if n[0] == x00 {
			*s = b
			return nil
		}

		b = append(b, n[0])
	}
}
