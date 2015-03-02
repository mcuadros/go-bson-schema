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

type ExampleBasic struct {
	Id              bson.ObjectId
	Nil             *ExampleBasic
	Map             map[string]string
	Slice           []string
	BoolYes         bool    `default:"true"`
	BoolNon         bool    `default:"false"`
	Float64         float64 `default:"64.21"`
	String          string  `default:"33"`
	Int64           int64   `default:"64"`
	Int8            int8    `default:"8"`
	Byte            []byte  `default:"bytes"`
	Timestamp       bson.MongoTimestamp
	JavaScript      bson.JavaScript
	JavaScriptScope bson.JavaScript
	RegEx           bson.RegEx
	Symbol          bson.Symbol
	Time            time.Time
}

func (s *BSSuite) TestValidator_ValidateStruct(c *C) {
	e := &ExampleBasic{}
	e.Time = time.Now()
	e.Map = map[string]string{"mapfoo": "qux"}
	e.Slice = []string{"qux", "foo"}
	e.Id = bson.NewObjectId()
	e.Timestamp = bson.MongoTimestamp(258)
	e.JavaScript = bson.JavaScript{Code: "alert('foo')"}
	e.JavaScriptScope = bson.JavaScript{"code", bson.M{"scope": nil}}
	e.RegEx = bson.RegEx{Pattern: "qux", Options: "baz"}
	e.Symbol = bson.Symbol("qux")
	defaults.SetDefaults(e)

	b, err := bson.Marshal(e)
	c.Assert(err, IsNil)

	v := &Validator{}
	v.Validate(bytes.NewReader(b))
}
