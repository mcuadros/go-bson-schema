package bsonschema

import (
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
	return v.decDocument(bson)
}

func (v *Validator) decDocument(bson io.Reader) error {
	var len int32
	if err := readInt32(bson, &len); err != nil {
		return nil
	}

	var err error
	for err == nil {
		err = v.decElement(bson)
	}

	fmt.Println("error", err)

	return nil
}

func (v *Validator) decElement(bson io.Reader) error {
	var typ [1]byte
	if _, err := bson.Read(typ[:]); err != nil {
		return err
	}

	if typ[0] == x00 {
		return io.EOF
	}

	var name []byte
	if err := readCString(bson, &name); err != nil {
		return err
	}

	fmt.Printf("-->  %q  <--\n", string(name))

	switch typ[0] {
	case TypeDocument, TypeArray:
		v.decDocument(bson)
	case TypeObjectId:
		v.validateObjectId(bson)
	case TypeBool:
		v.validateBool(bson)
	case TypeInt32:
		v.validateInt32(bson)
	case TypeInt64:
		v.validateInt64(bson)
	case TypeTimestamp:
		v.validateInt64(bson)
	case TypeDouble:
		v.validateDouble(bson)
	case TypeDate:
		v.validateDate(bson)
	case TypeBinary:
		v.validateBinary(bson)
	case TypeString, TypeCode, TypeSymbol:
		v.validateString(bson)
	case TypeCodeWithScope:
		v.validateString(bson)
		v.decDocument(bson)
	case TypeRegexp:
		v.validateRegexp(bson)
	case TypeNull:
		v.validateNil()
	case TypeMaxKey, TypeMinKey:
	default:
		fmt.Println("non-supported type", typ)
	}

	return nil
}

func (v *Validator) validateObjectId(bson io.Reader) error {
	var value [12]byte
	_, err := bson.Read(value[:])
	fmt.Println("----", value)

	return err
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

func (v *Validator) validateDouble(bson io.Reader) error {
	var value float64
	err := readFloat64(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateString(bson io.Reader) error {
	var value string
	err := readString(bson, &value)
	fmt.Println("----", value)

	return err
}

func (v *Validator) validateRegexp(bson io.Reader) error {
	var pattern, options []byte
	err := readCString(bson, &pattern)
	if err != nil {
		return err
	}

	err = readCString(bson, &options)
	fmt.Println("----", string(pattern), string(options))

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
	err := readBinary(bson, &value)
	fmt.Println("----", string(value))

	return err
}

func (v *Validator) validateNil() error {
	return nil
}
