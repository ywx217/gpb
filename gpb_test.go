package gpb

import (
	"testing"

	"github.com/ywx217/gpb/internal/testprotos"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func initGoTestField() *testprotos.GoTestField {
	f := new(testprotos.GoTestField)
	f.Label = proto.String("label")
	f.Type = proto.String("type")
	return f
}

// These are all structurally equivalent but the tag numbers differ.
// (It's remarkable that required, optional, and repeated all have
// 8 letters.)
func initGoTestRequiredGroup() *testprotos.GoTest_RequiredGroup {
	return &testprotos.GoTest_RequiredGroup{
		RequiredField: proto.String("required"),
	}
}

func initGoTestOptionalGroup() *testprotos.GoTest_OptionalGroup {
	return &testprotos.GoTest_OptionalGroup{
		RequiredField: proto.String("optional"),
	}
}

func initGoTestRepeatedGroup() *testprotos.GoTest_RepeatedGroup {
	return &testprotos.GoTest_RepeatedGroup{
		RequiredField: proto.String("repeated"),
	}
}

func initGoTest(setDefaults bool) *testprotos.GoTest {
	msg := &testprotos.GoTest{
		Kind:               testprotos.GoTest_TIME.Enum(),
		RequiredField:      initGoTestField(),
		F_BoolRequired:     proto.Bool(true),
		F_Int32Required:    proto.Int32(3),
		F_Int64Required:    proto.Int64(6),
		F_Fixed32Required:  proto.Uint32(32),
		F_Fixed64Required:  proto.Uint64(64),
		F_Uint32Required:   proto.Uint32(3232),
		F_Uint64Required:   proto.Uint64(6464),
		F_FloatRequired:    proto.Float32(3232),
		F_DoubleRequired:   proto.Float64(6464),
		F_StringRequired:   proto.String("string"),
		F_BytesRequired:    []byte("bytes"),
		F_Sint32Required:   proto.Int32(-32),
		F_Sint64Required:   proto.Int64(-64),
		F_Sfixed32Required: proto.Int32(-32),
		F_Sfixed64Required: proto.Int64(-64),
		Requiredgroup:      initGoTestRequiredGroup(),
	}
	if setDefaults {
		msg.F_BoolDefaulted = proto.Bool(testprotos.Default_GoTest_F_BoolDefaulted)
		msg.F_Int32Defaulted = proto.Int32(testprotos.Default_GoTest_F_Int32Defaulted)
		msg.F_Int64Defaulted = proto.Int64(testprotos.Default_GoTest_F_Int64Defaulted)
		msg.F_Fixed32Defaulted = proto.Uint32(testprotos.Default_GoTest_F_Fixed32Defaulted)
		msg.F_Fixed64Defaulted = proto.Uint64(testprotos.Default_GoTest_F_Fixed64Defaulted)
		msg.F_Uint32Defaulted = proto.Uint32(testprotos.Default_GoTest_F_Uint32Defaulted)
		msg.F_Uint64Defaulted = proto.Uint64(testprotos.Default_GoTest_F_Uint64Defaulted)
		msg.F_FloatDefaulted = proto.Float32(testprotos.Default_GoTest_F_FloatDefaulted)
		msg.F_DoubleDefaulted = proto.Float64(testprotos.Default_GoTest_F_DoubleDefaulted)
		msg.F_StringDefaulted = proto.String(testprotos.Default_GoTest_F_StringDefaulted)
		msg.F_BytesDefaulted = testprotos.Default_GoTest_F_BytesDefaulted
		msg.F_Sint32Defaulted = proto.Int32(testprotos.Default_GoTest_F_Sint32Defaulted)
		msg.F_Sint64Defaulted = proto.Int64(testprotos.Default_GoTest_F_Sint64Defaulted)
		msg.F_Sfixed32Defaulted = proto.Int32(testprotos.Default_GoTest_F_Sfixed32Defaulted)
		msg.F_Sfixed64Defaulted = proto.Int64(testprotos.Default_GoTest_F_Sfixed64Defaulted)
	}
	return msg
}

func verify(t *testing.T, expect, actual proto.Message) {
	require.True(t, proto.Equal(expect, actual))
}

func fverify(t *testing.T, expect proto.Message, unmarshaler func(bs []byte) proto.Message) {
	bs, err := proto.Marshal(expect)
	require.NoError(t, err)
	actual := unmarshaler(bs)
	verify(t, expect, actual)
}

// TestEncodeDecode1 all required fields set, no defaults provided.
func TestEncodeDecode1(t *testing.T) {
	fverify(t, initGoTest(false), func(bs []byte) proto.Message {
		return &testprotos.GoTest{
			Kind: testprotos.GoTest_KIND(GetOne(bs, 1).Int32()).Enum(),
			RequiredField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 4, 1).String()),
				Type:  proto.String(GetOne(bs, 4, 2).String()),
			},
			F_BoolRequired:     proto.Bool(GetOne(bs, 10).Bool()),
			F_Int32Required:    proto.Int32(GetOne(bs, 11).Int32()),
			F_Int64Required:    proto.Int64(GetOne(bs, 12).Int64()),
			F_Fixed32Required:  proto.Uint32(GetOne(bs, 13).Fixed32()),
			F_Fixed64Required:  proto.Uint64(GetOne(bs, 14).Fixed64()),
			F_Uint32Required:   proto.Uint32(GetOne(bs, 15).Uint32()),
			F_Uint64Required:   proto.Uint64(GetOne(bs, 16).Uint64()),
			F_FloatRequired:    proto.Float32(GetOne(bs, 17).Float32()),
			F_DoubleRequired:   proto.Float64(GetOne(bs, 18).Float64()),
			F_StringRequired:   proto.String(GetOne(bs, 19).String()),
			F_BytesRequired:    GetOne(bs, 101).Bytes(),
			F_Sint32Required:   proto.Int32(GetOne(bs, 102).Sint32()),
			F_Sint64Required:   proto.Int64(GetOne(bs, 103).Sint64()),
			F_Sfixed32Required: proto.Int32(GetOne(bs, 104).SFixed32()),
			F_Sfixed64Required: proto.Int64(GetOne(bs, 105).SFixed64()),
			Requiredgroup: &testprotos.GoTest_RequiredGroup{
				RequiredField: proto.String(GetOne(bs, 70, 71).String()),
			},
		}
	})
}

// TestEncodeDecode2 all required fields set, defaults provided.
func TestEncodeDecode2(t *testing.T) {
	fverify(t, initGoTest(true), func(bs []byte) proto.Message {
		return &testprotos.GoTest{
			Kind: testprotos.GoTest_KIND(GetOne(bs, 1).Int32()).Enum(),
			RequiredField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 4, 1).String()),
				Type:  proto.String(GetOne(bs, 4, 2).String()),
			},
			F_BoolRequired:    proto.Bool(GetOne(bs, 10).Bool()),
			F_Int32Required:   proto.Int32(GetOne(bs, 11).Int32()),
			F_Int64Required:   proto.Int64(GetOne(bs, 12).Int64()),
			F_Fixed32Required: proto.Uint32(GetOne(bs, 13).Fixed32()),
			F_Fixed64Required: proto.Uint64(GetOne(bs, 14).Fixed64()),
			F_Uint32Required:  proto.Uint32(GetOne(bs, 15).Uint32()),
			F_Uint64Required:  proto.Uint64(GetOne(bs, 16).Uint64()),
			F_FloatRequired:   proto.Float32(GetOne(bs, 17).Float32()),
			F_DoubleRequired:  proto.Float64(GetOne(bs, 18).Float64()),
			F_StringRequired:  proto.String(GetOne(bs, 19).String()),

			F_BoolDefaulted:    proto.Bool(GetOne(bs, 40).Bool()),
			F_Int32Defaulted:   proto.Int32(GetOne(bs, 41).Int32()),
			F_Int64Defaulted:   proto.Int64(GetOne(bs, 42).Int64()),
			F_Fixed32Defaulted: proto.Uint32(GetOne(bs, 43).Fixed32()),
			F_Fixed64Defaulted: proto.Uint64(GetOne(bs, 44).Fixed64()),
			F_Uint32Defaulted:  proto.Uint32(GetOne(bs, 45).Uint32()),
			F_Uint64Defaulted:  proto.Uint64(GetOne(bs, 46).Uint64()),
			F_FloatDefaulted:   proto.Float32(GetOne(bs, 47).Float32()),
			F_DoubleDefaulted:  proto.Float64(GetOne(bs, 48).Float64()),
			F_StringDefaulted:  proto.String(GetOne(bs, 49).String()),
			Requiredgroup: &testprotos.GoTest_RequiredGroup{
				RequiredField: proto.String(GetOne(bs, 70, 71).String()),
			},

			F_BytesRequired:     GetOne(bs, 101).Bytes(),
			F_Sint32Required:    proto.Int32(GetOne(bs, 102).Sint32()),
			F_Sint64Required:    proto.Int64(GetOne(bs, 103).Sint64()),
			F_Sfixed32Required:  proto.Int32(GetOne(bs, 104).SFixed32()),
			F_Sfixed64Required:  proto.Int64(GetOne(bs, 105).SFixed64()),
			F_BytesDefaulted:    GetOne(bs, 401).Bytes(),
			F_Sint32Defaulted:   proto.Int32(GetOne(bs, 402).Sint32()),
			F_Sint64Defaulted:   proto.Int64(GetOne(bs, 403).Sint64()),
			F_Sfixed32Defaulted: proto.Int32(GetOne(bs, 404).SFixed32()),
			F_Sfixed64Defaulted: proto.Int64(GetOne(bs, 405).SFixed64()),
		}
	})
}

// TestEncodeDecode3 all default fields set to their default value by hand.
func TestEncodeDecode3(t *testing.T) {
	msg := initGoTest(false)
	msg.F_BoolDefaulted = proto.Bool(true)
	msg.F_Int32Defaulted = proto.Int32(32)
	msg.F_Int64Defaulted = proto.Int64(64)
	msg.F_Fixed32Defaulted = proto.Uint32(320)
	msg.F_Fixed64Defaulted = proto.Uint64(640)
	msg.F_Uint32Defaulted = proto.Uint32(3200)
	msg.F_Uint64Defaulted = proto.Uint64(6400)
	msg.F_FloatDefaulted = proto.Float32(314159)
	msg.F_DoubleDefaulted = proto.Float64(271828)
	msg.F_StringDefaulted = proto.String("hello, \"world!\"\n")
	msg.F_BytesDefaulted = []byte("Bignose")
	msg.F_Sint32Defaulted = proto.Int32(-32)
	msg.F_Sint64Defaulted = proto.Int64(-64)
	msg.F_Sfixed32Defaulted = proto.Int32(-32)
	msg.F_Sfixed64Defaulted = proto.Int64(-64)

	fverify(t, msg, func(bs []byte) proto.Message {
		return &testprotos.GoTest{
			Kind: testprotos.GoTest_KIND(GetOne(bs, 1).Int32()).Enum(),
			RequiredField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 4, 1).String()),
				Type:  proto.String(GetOne(bs, 4, 2).String()),
			},
			F_BoolRequired:    proto.Bool(GetOne(bs, 10).Bool()),
			F_Int32Required:   proto.Int32(GetOne(bs, 11).Int32()),
			F_Int64Required:   proto.Int64(GetOne(bs, 12).Int64()),
			F_Fixed32Required: proto.Uint32(GetOne(bs, 13).Fixed32()),
			F_Fixed64Required: proto.Uint64(GetOne(bs, 14).Fixed64()),
			F_Uint32Required:  proto.Uint32(GetOne(bs, 15).Uint32()),
			F_Uint64Required:  proto.Uint64(GetOne(bs, 16).Uint64()),
			F_FloatRequired:   proto.Float32(GetOne(bs, 17).Float32()),
			F_DoubleRequired:  proto.Float64(GetOne(bs, 18).Float64()),
			F_StringRequired:  proto.String(GetOne(bs, 19).String()),

			F_BoolDefaulted:    proto.Bool(GetOne(bs, 40).Bool()),
			F_Int32Defaulted:   proto.Int32(GetOne(bs, 41).Int32()),
			F_Int64Defaulted:   proto.Int64(GetOne(bs, 42).Int64()),
			F_Fixed32Defaulted: proto.Uint32(GetOne(bs, 43).Fixed32()),
			F_Fixed64Defaulted: proto.Uint64(GetOne(bs, 44).Fixed64()),
			F_Uint32Defaulted:  proto.Uint32(GetOne(bs, 45).Uint32()),
			F_Uint64Defaulted:  proto.Uint64(GetOne(bs, 46).Uint64()),
			F_FloatDefaulted:   proto.Float32(GetOne(bs, 47).Float32()),
			F_DoubleDefaulted:  proto.Float64(GetOne(bs, 48).Float64()),
			F_StringDefaulted:  proto.String(GetOne(bs, 49).String()),
			Requiredgroup: &testprotos.GoTest_RequiredGroup{
				RequiredField: proto.String(GetOne(bs, 70, 71).String()),
			},

			F_BytesRequired:     GetOne(bs, 101).Bytes(),
			F_Sint32Required:    proto.Int32(GetOne(bs, 102).Sint32()),
			F_Sint64Required:    proto.Int64(GetOne(bs, 103).Sint64()),
			F_Sfixed32Required:  proto.Int32(GetOne(bs, 104).SFixed32()),
			F_Sfixed64Required:  proto.Int64(GetOne(bs, 105).SFixed64()),
			F_BytesDefaulted:    GetOne(bs, 401).Bytes(),
			F_Sint32Defaulted:   proto.Int32(GetOne(bs, 402).Sint32()),
			F_Sint64Defaulted:   proto.Int64(GetOne(bs, 403).Sint64()),
			F_Sfixed32Defaulted: proto.Int32(GetOne(bs, 404).SFixed32()),
			F_Sfixed64Defaulted: proto.Int64(GetOne(bs, 405).SFixed64()),
		}
	})
}

// TestEncodeDecode4 all required fields set, defaults provided, all non-defaulted optional fields have values.
func TestEncodeDecode4(t *testing.T) {
	msg := initGoTest(true)
	msg.Table = proto.String("hello")
	msg.Param = proto.Int32(7)
	msg.OptionalField = initGoTestField()
	msg.F_BoolOptional = proto.Bool(true)
	msg.F_Int32Optional = proto.Int32(32)
	msg.F_Int64Optional = proto.Int64(64)
	msg.F_Fixed32Optional = proto.Uint32(3232)
	msg.F_Fixed64Optional = proto.Uint64(6464)
	msg.F_Uint32Optional = proto.Uint32(323232)
	msg.F_Uint64Optional = proto.Uint64(646464)
	msg.F_FloatOptional = proto.Float32(32.)
	msg.F_DoubleOptional = proto.Float64(64.)
	msg.F_StringOptional = proto.String("hello")
	msg.F_BytesOptional = []byte("Bignose")
	msg.F_Sint32Optional = proto.Int32(-32)
	msg.F_Sint64Optional = proto.Int64(-64)
	msg.F_Sfixed32Optional = proto.Int32(-32)
	msg.F_Sfixed64Optional = proto.Int64(-64)
	msg.Optionalgroup = initGoTestOptionalGroup()

	fverify(t, msg, func(bs []byte) proto.Message {
		return &testprotos.GoTest{
			Kind:  testprotos.GoTest_KIND(GetOne(bs, 1).Int32()).Enum(),
			Table: proto.String(GetOne(bs, 2).String()),
			Param: proto.Int32(GetOne(bs, 3).Int32()),
			RequiredField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 4, 1).String()),
				Type:  proto.String(GetOne(bs, 4, 2).String()),
			},
			OptionalField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 6, 1).String()),
				Type:  proto.String(GetOne(bs, 6, 2).String()),
			},

			F_BoolRequired:    proto.Bool(GetOne(bs, 10).Bool()),
			F_Int32Required:   proto.Int32(GetOne(bs, 11).Int32()),
			F_Int64Required:   proto.Int64(GetOne(bs, 12).Int64()),
			F_Fixed32Required: proto.Uint32(GetOne(bs, 13).Fixed32()),
			F_Fixed64Required: proto.Uint64(GetOne(bs, 14).Fixed64()),
			F_Uint32Required:  proto.Uint32(GetOne(bs, 15).Uint32()),
			F_Uint64Required:  proto.Uint64(GetOne(bs, 16).Uint64()),
			F_FloatRequired:   proto.Float32(GetOne(bs, 17).Float32()),
			F_DoubleRequired:  proto.Float64(GetOne(bs, 18).Float64()),
			F_StringRequired:  proto.String(GetOne(bs, 19).String()),

			F_BoolOptional:    proto.Bool(GetOne(bs, 30).Bool()),
			F_Int32Optional:   proto.Int32(GetOne(bs, 31).Int32()),
			F_Int64Optional:   proto.Int64(GetOne(bs, 32).Int64()),
			F_Fixed32Optional: proto.Uint32(GetOne(bs, 33).Fixed32()),
			F_Fixed64Optional: proto.Uint64(GetOne(bs, 34).Fixed64()),
			F_Uint32Optional:  proto.Uint32(GetOne(bs, 35).Uint32()),
			F_Uint64Optional:  proto.Uint64(GetOne(bs, 36).Uint64()),
			F_FloatOptional:   proto.Float32(GetOne(bs, 37).Float32()),
			F_DoubleOptional:  proto.Float64(GetOne(bs, 38).Float64()),
			F_StringOptional:  proto.String(GetOne(bs, 39).String()),

			F_BoolDefaulted:    proto.Bool(GetOne(bs, 40).Bool()),
			F_Int32Defaulted:   proto.Int32(GetOne(bs, 41).Int32()),
			F_Int64Defaulted:   proto.Int64(GetOne(bs, 42).Int64()),
			F_Fixed32Defaulted: proto.Uint32(GetOne(bs, 43).Fixed32()),
			F_Fixed64Defaulted: proto.Uint64(GetOne(bs, 44).Fixed64()),
			F_Uint32Defaulted:  proto.Uint32(GetOne(bs, 45).Uint32()),
			F_Uint64Defaulted:  proto.Uint64(GetOne(bs, 46).Uint64()),
			F_FloatDefaulted:   proto.Float32(GetOne(bs, 47).Float32()),
			F_DoubleDefaulted:  proto.Float64(GetOne(bs, 48).Float64()),
			F_StringDefaulted:  proto.String(GetOne(bs, 49).String()),
			Requiredgroup: &testprotos.GoTest_RequiredGroup{
				RequiredField: proto.String(GetOne(bs, 70, 71).String()),
			},

			Optionalgroup: &testprotos.GoTest_OptionalGroup{
				RequiredField: proto.String(GetOne(bs, 90, 91).String()),
			},

			F_BytesRequired:     GetOne(bs, 101).Bytes(),
			F_Sint32Required:    proto.Int32(GetOne(bs, 102).Sint32()),
			F_Sint64Required:    proto.Int64(GetOne(bs, 103).Sint64()),
			F_Sfixed32Required:  proto.Int32(GetOne(bs, 104).SFixed32()),
			F_Sfixed64Required:  proto.Int64(GetOne(bs, 105).SFixed64()),
			F_BytesOptional:     GetOne(bs, 301).Bytes(),
			F_Sint32Optional:    proto.Int32(GetOne(bs, 302).Sint32()),
			F_Sint64Optional:    proto.Int64(GetOne(bs, 303).Sint64()),
			F_Sfixed32Optional:  proto.Int32(GetOne(bs, 304).SFixed32()),
			F_Sfixed64Optional:  proto.Int64(GetOne(bs, 305).SFixed64()),
			F_BytesDefaulted:    GetOne(bs, 401).Bytes(),
			F_Sint32Defaulted:   proto.Int32(GetOne(bs, 402).Sint32()),
			F_Sint64Defaulted:   proto.Int64(GetOne(bs, 403).Sint64()),
			F_Sfixed32Defaulted: proto.Int32(GetOne(bs, 404).SFixed32()),
			F_Sfixed64Defaulted: proto.Int64(GetOne(bs, 405).SFixed64()),
		}
	})
}

// TestEncodeDecode5 all required fields set, defaults provided, all repeated fields given two values.
func TestEncodeDecode5(t *testing.T) {
	msg := initGoTest(true)
	msg.RepeatedField = []*testprotos.GoTestField{initGoTestField(), initGoTestField()}
	msg.F_BoolRepeated = []bool{false, true}
	msg.F_Int32Repeated = []int32{32, 33}
	msg.F_Int64Repeated = []int64{64, 65}
	msg.F_Fixed32Repeated = []uint32{3232, 3333}
	msg.F_Fixed64Repeated = []uint64{6464, 6565}
	msg.F_Uint32Repeated = []uint32{323232, 333333}
	msg.F_Uint64Repeated = []uint64{646464, 656565}
	msg.F_FloatRepeated = []float32{32., 33.}
	msg.F_DoubleRepeated = []float64{64., 65.}
	msg.F_StringRepeated = []string{"hello", "sailor"}
	msg.F_BytesRepeated = [][]byte{[]byte("big"), []byte("nose")}
	msg.F_Sint32Repeated = []int32{32, -32}
	msg.F_Sint64Repeated = []int64{64, -64}
	msg.F_Sfixed32Repeated = []int32{32, -32}
	msg.F_Sfixed64Repeated = []int64{64, -64}
	msg.Repeatedgroup = []*testprotos.GoTest_RepeatedGroup{initGoTestRepeatedGroup(), initGoTestRepeatedGroup()}

	fverify(t, msg, func(bs []byte) proto.Message {
		return &testprotos.GoTest{
			Kind: testprotos.GoTest_KIND(GetOne(bs, 1).Int32()).Enum(),
			RequiredField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 4, 1).String()),
				Type:  proto.String(GetOne(bs, 4, 2).String()),
			},
			RepeatedField: lo.Map(GetAll(bs, 5), func(r Result, _ int) *testprotos.GoTestField {
				return &testprotos.GoTestField{
					Label: proto.String(r.GetOne(1).String()),
					Type:  proto.String(r.GetOne(2).String()),
				}
			}),

			F_BoolRequired:    proto.Bool(GetOne(bs, 10).Bool()),
			F_Int32Required:   proto.Int32(GetOne(bs, 11).Int32()),
			F_Int64Required:   proto.Int64(GetOne(bs, 12).Int64()),
			F_Fixed32Required: proto.Uint32(GetOne(bs, 13).Fixed32()),
			F_Fixed64Required: proto.Uint64(GetOne(bs, 14).Fixed64()),
			F_Uint32Required:  proto.Uint32(GetOne(bs, 15).Uint32()),
			F_Uint64Required:  proto.Uint64(GetOne(bs, 16).Uint64()),
			F_FloatRequired:   proto.Float32(GetOne(bs, 17).Float32()),
			F_DoubleRequired:  proto.Float64(GetOne(bs, 18).Float64()),
			F_StringRequired:  proto.String(GetOne(bs, 19).String()),

			// Repeated fields of all basic types
			F_BoolRepeated:    lo.Map(GetAll(bs, 20), func(r Result, _ int) bool { return r.Bool() }),
			F_Int32Repeated:   lo.Map(GetAll(bs, 21), func(r Result, _ int) int32 { return r.Int32() }),
			F_Int64Repeated:   lo.Map(GetAll(bs, 22), func(r Result, _ int) int64 { return r.Int64() }),
			F_Fixed32Repeated: lo.Map(GetAll(bs, 23), func(r Result, _ int) uint32 { return r.Fixed32() }),
			F_Fixed64Repeated: lo.Map(GetAll(bs, 24), func(r Result, _ int) uint64 { return r.Fixed64() }),
			F_Uint32Repeated:  lo.Map(GetAll(bs, 25), func(r Result, _ int) uint32 { return r.Uint32() }),
			F_Uint64Repeated:  lo.Map(GetAll(bs, 26), func(r Result, _ int) uint64 { return r.Uint64() }),
			F_FloatRepeated:   lo.Map(GetAll(bs, 27), func(r Result, _ int) float32 { return r.Float32() }),
			F_DoubleRepeated:  lo.Map(GetAll(bs, 28), func(r Result, _ int) float64 { return r.Float64() }),
			F_StringRepeated:  lo.Map(GetAll(bs, 29), func(r Result, _ int) string { return r.String() }),

			F_BoolDefaulted:    proto.Bool(GetOne(bs, 40).Bool()),
			F_Int32Defaulted:   proto.Int32(GetOne(bs, 41).Int32()),
			F_Int64Defaulted:   proto.Int64(GetOne(bs, 42).Int64()),
			F_Fixed32Defaulted: proto.Uint32(GetOne(bs, 43).Fixed32()),
			F_Fixed64Defaulted: proto.Uint64(GetOne(bs, 44).Fixed64()),
			F_Uint32Defaulted:  proto.Uint32(GetOne(bs, 45).Uint32()),
			F_Uint64Defaulted:  proto.Uint64(GetOne(bs, 46).Uint64()),
			F_FloatDefaulted:   proto.Float32(GetOne(bs, 47).Float32()),
			F_DoubleDefaulted:  proto.Float64(GetOne(bs, 48).Float64()),
			F_StringDefaulted:  proto.String(GetOne(bs, 49).String()),
			Requiredgroup: &testprotos.GoTest_RequiredGroup{
				RequiredField: proto.String(GetOne(bs, 70, 71).String()),
			},
			Repeatedgroup: lo.Map(GetAll(bs, 80), func(r Result, _ int) *testprotos.GoTest_RepeatedGroup {
				return &testprotos.GoTest_RepeatedGroup{
					RequiredField: proto.String(r.GetOne(81).String()),
				}
			}),

			F_BytesRequired:     GetOne(bs, 101).Bytes(),
			F_Sint32Required:    proto.Int32(GetOne(bs, 102).Sint32()),
			F_Sint64Required:    proto.Int64(GetOne(bs, 103).Sint64()),
			F_Sfixed32Required:  proto.Int32(GetOne(bs, 104).SFixed32()),
			F_Sfixed64Required:  proto.Int64(GetOne(bs, 105).SFixed64()),
			F_BytesRepeated:     lo.Map(GetAll(bs, 201), func(r Result, _ int) []byte { return r.Bytes() }),
			F_Sint32Repeated:    lo.Map(GetAll(bs, 202), func(r Result, _ int) int32 { return r.Sint32() }),
			F_Sint64Repeated:    lo.Map(GetAll(bs, 203), func(r Result, _ int) int64 { return r.Sint64() }),
			F_Sfixed32Repeated:  lo.Map(GetAll(bs, 204), func(r Result, _ int) int32 { return r.SFixed32() }),
			F_Sfixed64Repeated:  lo.Map(GetAll(bs, 205), func(r Result, _ int) int64 { return r.SFixed64() }),
			F_BytesDefaulted:    GetOne(bs, 401).Bytes(),
			F_Sint32Defaulted:   proto.Int32(GetOne(bs, 402).Sint32()),
			F_Sint64Defaulted:   proto.Int64(GetOne(bs, 403).Sint64()),
			F_Sfixed32Defaulted: proto.Int32(GetOne(bs, 404).SFixed32()),
			F_Sfixed64Defaulted: proto.Int64(GetOne(bs, 405).SFixed64()),
		}
	})
}

// TestEncodeDecode6 all required fields set, all packed repeated fields given two values.
func TestEncodeDecode6(t *testing.T) {
	msg := initGoTest(false)
	msg.F_BoolRepeatedPacked = []bool{false, true}
	msg.F_Int32RepeatedPacked = []int32{32, 33}
	msg.F_Int64RepeatedPacked = []int64{64, 65}
	msg.F_Fixed32RepeatedPacked = []uint32{3232, 3333}
	msg.F_Fixed64RepeatedPacked = []uint64{6464, 6565}
	msg.F_Uint32RepeatedPacked = []uint32{323232, 333333}
	msg.F_Uint64RepeatedPacked = []uint64{646464, 656565}
	msg.F_FloatRepeatedPacked = []float32{32., 33.}
	msg.F_DoubleRepeatedPacked = []float64{64., 65.}
	msg.F_Sint32RepeatedPacked = []int32{32, -32}
	msg.F_Sint64RepeatedPacked = []int64{64, -64}
	msg.F_Sfixed32RepeatedPacked = []int32{32, -32}
	msg.F_Sfixed64RepeatedPacked = []int64{64, -64}

	fverify(t, msg, func(bs []byte) proto.Message {
		return &testprotos.GoTest{
			Kind: testprotos.GoTest_KIND(GetOne(bs, 1).Int32()).Enum(),
			RequiredField: &testprotos.GoTestField{
				Label: proto.String(GetOne(bs, 4, 1).String()),
				Type:  proto.String(GetOne(bs, 4, 2).String()),
			},

			F_BoolRequired:    proto.Bool(GetOne(bs, 10).Bool()),
			F_Int32Required:   proto.Int32(GetOne(bs, 11).Int32()),
			F_Int64Required:   proto.Int64(GetOne(bs, 12).Int64()),
			F_Fixed32Required: proto.Uint32(GetOne(bs, 13).Fixed32()),
			F_Fixed64Required: proto.Uint64(GetOne(bs, 14).Fixed64()),
			F_Uint32Required:  proto.Uint32(GetOne(bs, 15).Uint32()),
			F_Uint64Required:  proto.Uint64(GetOne(bs, 16).Uint64()),
			F_FloatRequired:   proto.Float32(GetOne(bs, 17).Float32()),
			F_DoubleRequired:  proto.Float64(GetOne(bs, 18).Float64()),
			F_StringRequired:  proto.String(GetOne(bs, 19).String()),

			// Packed repeated fields (no string or bytes).
			F_BoolRepeatedPacked:    lo.Map(GetOne(bs, 50).UnpackVarint(), func(r Result, _ int) bool { return r.Bool() }),
			F_Int32RepeatedPacked:   lo.Map(GetOne(bs, 51).UnpackVarint(), func(r Result, _ int) int32 { return r.Int32() }),
			F_Int64RepeatedPacked:   lo.Map(GetOne(bs, 52).UnpackVarint(), func(r Result, _ int) int64 { return r.Int64() }),
			F_Fixed32RepeatedPacked: lo.Map(GetOne(bs, 53).UnpackFixed32(), func(r Result, _ int) uint32 { return r.Fixed32() }),
			F_Fixed64RepeatedPacked: lo.Map(GetOne(bs, 54).UnpackFixed64(), func(r Result, _ int) uint64 { return r.Fixed64() }),
			F_Uint32RepeatedPacked:  lo.Map(GetOne(bs, 55).UnpackVarint(), func(r Result, _ int) uint32 { return r.Uint32() }),
			F_Uint64RepeatedPacked:  lo.Map(GetOne(bs, 56).UnpackVarint(), func(r Result, _ int) uint64 { return r.Uint64() }),
			F_FloatRepeatedPacked:   lo.Map(GetOne(bs, 57).UnpackFixed32(), func(r Result, _ int) float32 { return r.Float32() }),
			F_DoubleRepeatedPacked:  lo.Map(GetOne(bs, 58).UnpackFixed64(), func(r Result, _ int) float64 { return r.Float64() }),

			Requiredgroup: &testprotos.GoTest_RequiredGroup{
				RequiredField: proto.String(GetOne(bs, 70, 71).String()),
			},

			F_BytesRequired:    GetOne(bs, 101).Bytes(),
			F_Sint32Required:   proto.Int32(GetOne(bs, 102).Sint32()),
			F_Sint64Required:   proto.Int64(GetOne(bs, 103).Sint64()),
			F_Sfixed32Required: proto.Int32(GetOne(bs, 104).SFixed32()),
			F_Sfixed64Required: proto.Int64(GetOne(bs, 105).SFixed64()),

			F_Sint32RepeatedPacked:   lo.Map(GetOne(bs, 502).UnpackVarint(), func(r Result, _ int) int32 { return r.Sint32() }),
			F_Sint64RepeatedPacked:   lo.Map(GetOne(bs, 503).UnpackVarint(), func(r Result, _ int) int64 { return r.Sint64() }),
			F_Sfixed32RepeatedPacked: lo.Map(GetOne(bs, 504).UnpackFixed32(), func(r Result, _ int) int32 { return r.SFixed32() }),
			F_Sfixed64RepeatedPacked: lo.Map(GetOne(bs, 505).UnpackFixed64(), func(r Result, _ int) int64 { return r.SFixed64() }),
		}
	})
}

// payload size: 2 bytes
func benchmarkTiny(b *testing.B, run func(*testing.B, []byte)) {
	b.StopTimer()
	msg := &testprotos.GoEnum{
		Foo: testprotos.FOO_FOO1.Enum(),
	}
	raw, err := proto.Marshal(msg)
	require.NoError(b, err)

	b.StartTimer()
	run(b, raw)
}

func BenchmarkGoProtobufTiny(b *testing.B) {
	benchmarkTiny(b, func(b *testing.B, raw []byte) {
		var sum int
		for i := 0; i < b.N; i++ {
			var decoded testprotos.GoEnum
			_ = proto.Unmarshal(raw, &decoded)
			sum += int(decoded.GetFoo())
		}
	})
}

func BenchmarkGpbTiny(b *testing.B) {
	benchmarkTiny(b, func(b *testing.B, raw []byte) {
		var sum int
		for i := 0; i < b.N; i++ {
			sum += int(GetOne(raw, 1).Int32())
		}
	})
}

// payload size: 451 bytes
func benchmarkSmall(b *testing.B, run func(*testing.B, []byte)) {
	b.StopTimer()
	msg := initGoTest(true)
	msg.RepeatedField = []*testprotos.GoTestField{initGoTestField(), initGoTestField()}
	msg.F_BoolRepeated = []bool{false, true}
	msg.F_Int32Repeated = []int32{32, 33}
	msg.F_Int64Repeated = []int64{64, 65}
	msg.F_Fixed32Repeated = []uint32{3232, 3333}
	msg.F_Fixed64Repeated = []uint64{6464, 6565}
	msg.F_Uint32Repeated = []uint32{323232, 333333}
	msg.F_Uint64Repeated = []uint64{646464, 656565}
	msg.F_FloatRepeated = []float32{32., 33.}
	msg.F_DoubleRepeated = []float64{64., 65.}
	msg.F_StringRepeated = []string{"hello", "sailor"}
	msg.F_BytesRepeated = [][]byte{[]byte("big"), []byte("nose")}
	msg.F_Sint32Repeated = []int32{32, -32}
	msg.F_Sint64Repeated = []int64{64, -64}
	msg.F_Sfixed32Repeated = []int32{32, -32}
	msg.F_Sfixed64Repeated = []int64{64, -64}
	msg.Repeatedgroup = []*testprotos.GoTest_RepeatedGroup{initGoTestRepeatedGroup(), initGoTestRepeatedGroup()}
	raw, err := proto.Marshal(msg)
	require.NoError(b, err)

	b.StartTimer()
	run(b, raw)
}

func BenchmarkGoProtobufSmall(b *testing.B) {
	benchmarkSmall(b, func(b *testing.B, raw []byte) {
		var sum int
		for i := 0; i < b.N; i++ {
			var decoded testprotos.GoTest
			_ = proto.Unmarshal(raw, &decoded)
			sum += int(decoded.GetF_Int32Required())
		}
	})
}

func BenchmarkGpbSmall(b *testing.B) {
	benchmarkSmall(b, func(b *testing.B, raw []byte) {
		var sum int
		for i := 0; i < b.N; i++ {
			sum += int(GetOne(raw, 11).Int32())
		}
	})
}
