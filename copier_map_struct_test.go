package copier_test

import (
	"github.com/jinzhu/copier"
	"strconv"
	"testing"
)

type Struct struct {
	FieldString    string
	FieldInt       int
	FieldFloat64   float64
	FieldNotExist  int
	SuperRule      string
	FieldTagIgnore string `copier:"-"`
}

type StructTagName struct {
	FieldA string `copier:"FieldB"`
}

func (s *Struct) Role(role string) {
	s.SuperRule = "Super " + role
}

type StructHasSlice struct {
	FieldSliceString []string
	FieldSliceInt    []int
}

type StructNested struct {
	Struct
	FieldA int64
}

func TestCopyMap2Struct(t *testing.T) {
	ts := &Struct{}
	src := map[string]interface{}{
		"FieldString":  "aaa",
		"FieldInt":     "11",
		"FieldFloat64": "1.22",
		"Role":         "role",
	}

	opt := copier.Option{Converters: []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: copier.Int,
			Fn: func(src interface{}) (dst interface{}, err error) {
				return strconv.Atoi(src.(string))
			},
		}, {
			SrcType: copier.String,
			DstType: copier.Float64,
			Fn: func(src interface{}) (dst interface{}, err error) {
				return strconv.ParseFloat(src.(string), 64)
			},
		},
	}}
	if err := copier.CopyWithOption(&ts, src, opt); err != nil {
		t.Errorf("err:%+v", err)
	}

	fieldInt, _ := strconv.Atoi(src["FieldInt"].(string))
	fieldFloat64, _ := strconv.ParseFloat(src["FieldFloat64"].(string), 64)
	if ts.FieldString != src["FieldString"] ||
		ts.FieldInt != fieldInt ||
		ts.FieldFloat64 != fieldFloat64 ||
		ts.SuperRule != "Super "+src["Role"].(string) {
		t.Errorf("Should be able to copy from ts to ts")
	}
}

func TestCopyMap2StructTagIgnore(t *testing.T) {
	ts := &Struct{FieldTagIgnore: "ewqewq"}
	src := map[string]interface{}{
		"FieldString":    "aaa",
		"FieldInt":       "11",
		"FieldFloat64":   "1.22",
		"Role":           "role",
		"FieldTagIgnore": "iii",
	}

	opt := copier.Option{Converters: []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: copier.Int,
			Fn: func(src interface{}) (dst interface{}, err error) {
				return strconv.Atoi(src.(string))
			},
		}, {
			SrcType: copier.String,
			DstType: copier.Float64,
			Fn: func(src interface{}) (dst interface{}, err error) {
				return strconv.ParseFloat(src.(string), 64)
			},
		},
	}}
	if err := copier.CopyWithOption(&ts, src, opt); err != nil {
		t.Errorf("err:%+v", err)
	}

	if ts.FieldTagIgnore == src["FieldTagIgnore"] {
		t.Error("Was not expected to copy FieldTagIgnore")
	}
	if ts.FieldTagIgnore != "ewqewq" {
		t.Error("Original FieldTagIgnore was overwritten")
	}
}

func TestCopyMap2StructTagName(t *testing.T) {
	ts := &StructTagName{}
	src := map[string]interface{}{
		"FieldB": "aaa",
	}

	if err := copier.Copy(&ts, src); err != nil {
		t.Errorf("err:%+v", err)
	}

	if ts.FieldA != src["FieldB"] {
		t.Errorf("Should be able to copy from ts to ts")
	}
}

func TestCopyMap2StructNested(t *testing.T) {
	ts := &StructNested{}
	src := map[string]interface{}{
		"Struct": map[string]interface{}{
			"FieldString":  "aaa",
			"FieldInt":     "11",
			"FieldFloat64": "1.22",
			"Role":         "role",
		},
		"FieldA": int64(11111),
	}

	opt := copier.Option{Converters: []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: copier.Int,
			Fn: func(src interface{}) (dst interface{}, err error) {
				return strconv.Atoi(src.(string))
			},
		}, {
			SrcType: copier.String,
			DstType: copier.Float64,
			Fn: func(src interface{}) (dst interface{}, err error) {
				return strconv.ParseFloat(src.(string), 64)
			},
		},
	}}
	if err := copier.CopyWithOption(&ts, src, opt); err != nil {
		t.Errorf("err:%+v", err)
	}

	nestedMap := src["Struct"].(map[string]interface{})
	fieldInt, _ := strconv.Atoi(nestedMap["FieldInt"].(string))
	fieldFloat64, _ := strconv.ParseFloat(nestedMap["FieldFloat64"].(string), 64)
	if ts.FieldString != nestedMap["FieldString"] ||
		ts.FieldInt != fieldInt ||
		ts.FieldFloat64 != fieldFloat64 ||
		ts.SuperRule != "Super "+nestedMap["Role"].(string) ||
		ts.FieldA != src["FieldA"].(int64) {
		t.Errorf("Should be able to copy from ts to ts nested struct")
	}
	if ts.FieldA != src["FieldA"].(int64) {
		t.Errorf("Should be able to copy from ts to ts")
	}
}
