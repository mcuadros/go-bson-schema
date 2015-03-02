package bsonschema

import (
	"bytes"
	"testing"
	"time"

	"github.com/mcuadros/go-defaults"
	. "gopkg.in/check.v1"
	"labix.org/v2/mgo/bson"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type BSSuite struct{}

var _ = Suite(&BSSuite{})

func (s *BSSuite) TestValidator_ValidateMap(c *C) {
	b, err := bson.Marshal(map[string]int{"foo_qux": 1, "bar": 2})
	c.Assert(err, IsNil)

	v := &Validator{}
	v.Validate(bytes.NewReader(b))
}

type ExampleBasic struct {
	BoolYes bool   `default:"true"`
	BoolNon bool   `default:"false"`
	String  string `default:"33"`
	Int64   int64  `default:"64"`
	Int8    int8   `default:"8"`
	Byte    []byte `default:"bytes"`
	Time    time.Time
}

func (s *BSSuite) TestValidator_ValidateStruct(c *C) {
	e := &ExampleBasic{}
	e.Time = time.Now()
	defaults.SetDefaults(e)

	b, err := bson.Marshal(e)
	c.Assert(err, IsNil)

	v := &Validator{}
	v.Validate(bytes.NewReader(b))
}
