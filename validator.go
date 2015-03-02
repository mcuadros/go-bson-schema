package bsonschema

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const TypeString = 0x02
const TypeObjectId = 0x07
const TypeBool = 0x08
const TypeDate = 0x09
const TypeInt32 = 0x10
const TypeInt64 = 0x12
const TypeDocument = 0x03
const TypeArray = 0x04
const TypeBinary = 0x05
const TypeDouble = 0x01
const TypeCode = 0x0d
const TypeCodeWithScope = 0x0f
const TypeNull = 0x0a
const TypeRegexp = 0x0b
const TypeSymbol = 0x0e
const TypeTimestamp = 0x11
const TypeMaxKey = 0x7F
const TypeMinKey = 0xFF

type Validator struct{}

func (v *Validator) Validate(bson io.Reader) error {
	var len int32
	if err := readInt32(bson, &len); err != nil {
		return nil
	}

	var err error
	for err == nil {
		err = v.decElement(bson)
	}

	return err
}

func (v *Validator) decElement(bson io.Reader) error {
	var typ [1]byte
	if _, err := bson.Read(typ[:]); err != nil {
		return err
	}

	var name []byte
	if err := readCString(bson, &name); err != nil {
		return err
	}

	fmt.Println("--", typ[0] == TypeInt32, string(name))

	switch typ[0] {
	case TypeBool:
		v.validateBool(bson)
	case TypeInt32:
		v.validateInt32(bson)
	case TypeInt64:
		v.validateInt64(bson)
	case TypeString:
		v.validateString(bson)
	case TypeDate:
		v.validateDate(bson)
	case TypeBinary:
		v.validateBinary(bson)
	default:
		fmt.Println("non-supported type", typ)
	}

	return nil
}

func (v *Validator) validateBool(bson io.Reader) error {
	var value byte
	err := readByte(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateInt32(bson io.Reader) error {
	var value int32
	err := readInt32(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateInt64(bson io.Reader) error {
	var value int64
	err := readInt64(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateString(bson io.Reader) error {
	var value string
	err := readString(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateDate(bson io.Reader) error {
	var value int64
	err := readInt64(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateBinary(bson io.Reader) error {
	var value []byte
	err := readBytes(bson, &value)
	fmt.Println("----", value)

	return err
}

const x00 = byte(0)

func readInt32(r io.Reader, i *int32) error {
	return binary.Read(r, binary.LittleEndian, i)
}

func readInt64(r io.Reader, i *int64) error {
	return binary.Read(r, binary.LittleEndian, i)
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

		b = append(b, n[0])
		if n[0] == x00 {
			*s = b
			return nil
		}
	}
}
