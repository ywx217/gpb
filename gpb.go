package gpb

import (
	"unsafe"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protowire"
)

// Parse protobuf wire format without idl definition.
// Inspired by tidwall/gjson
//
// Protobuf wire format specification: https://developers.google.com/protocol-buffers/docs/encoding

var (
	ErrUnknownWireType  = errors.New("unknown wire type")
	ErrInvalidLength    = errors.New("invalid length")
	ErrEndGroupNotFound = errors.New("end group not found")
)

const InvalidWireType protowire.Type = -1

type Result struct {
	// WireType the fields wire type enumerate
	//   type==0, used for: int32, int64, uint32, uint64, sint32, sint64, bool, enum
	//   type==1, used for: fixed64, sfixed64, double
	//   type==5, used for: fixed32, sfixed32, float
	//   type==2, used for: string, bytes, embedded messages, packed repeated fields
	WireType protowire.Type

	Varint uint64
	Raw    []byte // message without length header
}

// GetOne gets the first value by the given field numbers. `pbNumbers` indicates
// the path to retrieve the desired field.
// **Attention**: when the field is repeated in proto2, the first item is returned.
//   In proto3 or proto2 packed mode, the first packed group is returned, and
//   the UnpackVarint / UnpackFixed32 / UnpackFixed64 should be called to break
//   a single length-delimited frame into multiple Results.
func GetOne(pb []byte, pbNumbers ...protowire.Number) Result {
	state := Result{Raw: pb}
	return state.GetOne(pbNumbers...)
}

// GetAll gets all the values by the given field numbers. `pbNumbers` indicates
// the path to retrieve the desired field.
func GetAll(pb []byte, pbNumbers ...protowire.Number) []Result {
	state := Result{Raw: pb}
	return state.GetAll(pbNumbers...)
}

// GetOne gets the first field by the given field numbers. `pbNumbers` indicates
// the path to retrieve the desired field. There is no heap-memory allocation in this
// function.
//
// **Attention**: when the field is repeated in proto2, the first item is returned.
//   In proto3 or proto2 packed mode, the first packed group is returned, and
//   the UnpackVarint / UnpackFixed32 / UnpackFixed64 should be called to break
//   a single length-delimited frame into multiple Results.
func (r Result) GetOne(pbNumbers ...protowire.Number) (result Result) {
	result.WireType = InvalidWireType
	// use callback to avoid heap memory allocation
	_ = r.GetIter(func(r Result) bool {
		result = r
		return false
	}, pbNumbers...)
	return
}

// GetAll unlike GetOne, GetAll returns all the values by the given field numbers.
func (r Result) GetAll(pbNumbers ...protowire.Number) []Result {
	results := make([]Result, 0)
	_ = r.GetIter(func(r Result) bool {
		results = append(results, r)
		return true
	}, pbNumbers...)
	return results
}

// GetIter like GetAll, gets all the values until the resultSink returns false.
// Both GetOne and GetAll are implemented by this function.
func (r Result) GetIter(resultSink func(Result) bool, pbNumbers ...protowire.Number) (err error) {
	// use recursion calls to iterate through the data in depth first order.
	// according to the BenchmarkOptimized, dfs has a better performance than bfs.
	var skip bool
	var depth int
	var walkFunc func(Result) bool
	walkFunc = func(it Result) bool {
		if skip {
			return false
		}
		if depth == len(pbNumbers)-1 {
			if !resultSink(it) {
				skip = true
				return false
			}
			return true
		}
		// recursion call
		depth++
		if _, innerErr := it.IterFields(pbNumbers[depth], walkFunc); innerErr != nil {
			err = innerErr
			skip = true
			return false
		}
		depth--
		return true
	}
	if _, innerErr := r.IterFields(pbNumbers[depth], walkFunc); innerErr != nil {
		err = innerErr
	}
	return
}

// IterFields read through the binary data stored in r.Raw field-by-field, skipping all the fields
// not interested in.
func (r Result) IterFields(pbNumber protowire.Number, resultSink func(r Result) bool) (int, error) {
	var field Result
	var consumedLength int
	pb := r.Raw
	// fields are not organized in order, so we need to iterate through all fields
	for len(pb) > 0 {
		fieldNumber, wireType, totalLen := protowire.ConsumeTag(pb)
		if totalLen < 0 {
			// error occurred when totalLen is negative
			return consumedLength, ErrInvalidLength
		}
		pb = pb[totalLen:]
		consumedLength += totalLen

		field.WireType = wireType
		field.Varint = 0
		switch wireType {
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(pb)
			if n < 0 {
				return consumedLength, ErrInvalidLength
			}
			field.Varint = v
			field.Raw = pb[:n]
			pb = pb[n:]
			consumedLength += n
		case protowire.Fixed32Type:
			if len(pb) < 4 {
				return consumedLength, ErrInvalidLength
			}
			field.Raw = pb[:4]
			pb = pb[4:]
			consumedLength += 4
		case protowire.Fixed64Type:
			if len(pb) < 8 {
				return consumedLength, ErrInvalidLength
			}
			field.Raw = pb[:8]
			pb = pb[8:]
			consumedLength += 8
		case protowire.BytesType:
			// consume length varint
			v, n := protowire.ConsumeVarint(pb)
			if n < 0 {
				return consumedLength, ErrInvalidLength
			}
			pb = pb[n:]
			field.Raw = pb[:v]
			pb = pb[v:]
			consumedLength += n + int(v)
		case protowire.StartGroupType:
			// deprecated start group type, we need to consume all values to the end of the group
			var endGroupOccurred bool
			subGroup := Result{
				Raw: pb,
			}
			groupLength, err := subGroup.IterFields(fieldNumber, func(r Result) bool {
				// consume all fields inside the group until end group tag occurs
				if r.WireType == protowire.EndGroupType {
					endGroupOccurred = true
					return false // stop iteration
				}
				return true
			})
			if err != nil {
				return consumedLength, err
			}
			if !endGroupOccurred {
				return consumedLength, ErrEndGroupNotFound
			}
			field.Raw = pb[:groupLength]
			pb = pb[groupLength:]
			consumedLength += groupLength
			// consume end group again
			_, _, endGroupLen := protowire.ConsumeTag(pb)
			if endGroupLen < 0 {
				return consumedLength, ErrInvalidLength
			}
			// no need to check the tag wire type and field number,
			// because we already checked it in subGroup.IterFields
			pb = pb[endGroupLen:]
			consumedLength += endGroupLen
		case protowire.EndGroupType:
			// end group type, we need to consume the end group tag and give it to result consumer
			consumedLength -= totalLen
		default:
			return consumedLength, errors.WithMessagef(ErrUnknownWireType, "wire_type=%d", wireType)
		}

		if fieldNumber != pbNumber {
			// field number not match, read for the following fields
			continue
		}
		if !resultSink(field) {
			return consumedLength, nil
		}
	}
	return consumedLength, nil
}

// Varints - normal

func (r Result) Int32() int32 {
	if r.WireType != protowire.VarintType {
		return 0
	}
	return *(*int32)(unsafe.Pointer(&r.Varint))
}

func (r Result) Int64() int64 {
	if r.WireType != protowire.VarintType {
		return 0
	}
	return *(*int64)(unsafe.Pointer(&r.Varint))
}

func (r Result) Uint32() uint32 {
	if r.WireType != protowire.VarintType {
		return 0
	}
	return *(*uint32)(unsafe.Pointer(&r.Varint))
}

func (r Result) Uint64() uint64 {
	if r.WireType != protowire.VarintType {
		return 0
	}
	return r.Varint
}

func (r Result) Bool() bool {
	if r.WireType != protowire.VarintType {
		return false
	}
	return protowire.DecodeBool(r.Varint)
}

// Varints - zigzag (sint32 | sint64)

func (r Result) Sint32() int32 {
	if r.WireType != protowire.VarintType {
		return 0
	}
	v := protowire.DecodeZigZag(r.Varint)
	return *(*int32)(unsafe.Pointer(&v))
}

func (r Result) Sint64() int64 {
	if r.WireType != protowire.VarintType {
		return 0
	}
	return protowire.DecodeZigZag(r.Varint)
}

// Fixed32

func (r Result) Float32() float32 {
	if r.WireType != protowire.Fixed32Type {
		return 0
	}
	v, n := protowire.ConsumeFixed32(r.Raw)
	if n < 0 {
		return 0
	}
	return *(*float32)(unsafe.Pointer(&v))
}

func (r Result) Fixed32() uint32 {
	if r.WireType != protowire.Fixed32Type {
		return 0
	}
	v, n := protowire.ConsumeFixed32(r.Raw)
	if n < 0 {
		return 0
	}
	return v
}

func (r Result) SFixed32() int32 {
	if r.WireType != protowire.Fixed32Type {
		return 0
	}
	v, n := protowire.ConsumeFixed32(r.Raw)
	if n < 0 {
		return 0
	}
	return *(*int32)(unsafe.Pointer(&v))
}

// Fixed64

// Float64 parses the 64-bit raw bytes into float64. When the wire type is not Fixed64Type, it returns 0.
func (r Result) Float64() float64 {
	if r.WireType != protowire.Fixed64Type {
		return 0
	}
	v, n := protowire.ConsumeFixed64(r.Raw)
	if n < 0 {
		return 0
	}
	return *(*float64)(unsafe.Pointer(&v))
}

// Fixed64 parses the 64-bit raw bytes into uint64. When the wire type is not Fixed64Type, it returns 0.
func (r Result) Fixed64() uint64 {
	if r.WireType != protowire.Fixed64Type {
		return 0
	}
	v, n := protowire.ConsumeFixed64(r.Raw)
	if n < 0 {
		return 0
	}
	return v
}

// SFixed64 parses the 64-bit raw bytes into int64. When the wire type is not Fixed64Type, it returns 0.
func (r Result) SFixed64() int64 {
	if r.WireType != protowire.Fixed64Type {
		return 0
	}
	v, n := protowire.ConsumeFixed64(r.Raw)
	if n < 0 {
		return 0
	}
	return *(*int64)(unsafe.Pointer(&v))
}

// Length-delimited

// String parses the result as a string. When the wire type is not BytesType, empty string will be returned.
func (r Result) String() string {
	if r.WireType != protowire.BytesType {
		return ""
	}
	return string(r.Raw)
}

// Bytes parses the result as a byte slice. When the wire type is not BytesType, nil will be returned.
func (r Result) Bytes() []byte {
	if r.WireType != protowire.BytesType {
		return nil
	}
	return r.Raw
}

// Exist checks if it is a valid result.
func (r Result) Exist() bool {
	return r.WireType != InvalidWireType
}

// Unpack make packed scalars into []Result.
// In proto2, repeated fields of primitive numeric types can be declared as packed.
// In proto3, the fields are packed by default.
func (r Result) Unpack(itemType protowire.Type) []Result {
	switch itemType {
	case protowire.VarintType:
		return r.UnpackVarint()
	case protowire.Fixed32Type:
		return r.UnpackFixed32()
	case protowire.Fixed64Type:
		return r.UnpackFixed64()
	default:
		return nil
	}
}

// UnpackVarint unpacks the length-delimited varint-encoded binary data into separate results.
func (r Result) UnpackVarint() []Result {
	if r.WireType != protowire.BytesType {
		return nil
	}
	pb := r.Raw
	results := make([]Result, 0, len(pb))
	for len(pb) > 0 {
		v, n := protowire.ConsumeVarint(pb)
		if n < 0 {
			break
		}
		results = append(results, Result{
			WireType: protowire.VarintType,
			Varint:   v,
			Raw:      pb[:n],
		})
		pb = pb[n:]
	}
	return results
}

// UnpackFixed32 unpacks the length-delimited fixed32-encoded binary data into separate results.
func (r Result) UnpackFixed32() []Result {
	if r.WireType != protowire.BytesType {
		return nil
	}
	pb := r.Raw
	results := make([]Result, 0, len(pb)/4)
	for len(pb) >= 4 {
		results = append(results, Result{
			WireType: protowire.Fixed32Type,
			Raw:      pb[:4],
		})
		pb = pb[4:]
	}
	return results
}

// UnpackFixed64 unpacks the length-delimited fixed64-encoded binary data into separate results.
func (r Result) UnpackFixed64() []Result {
	if r.WireType != protowire.BytesType {
		return nil
	}
	pb := r.Raw
	results := make([]Result, 0, len(pb)/8)
	for len(pb) >= 8 {
		results = append(results, Result{
			WireType: protowire.Fixed64Type,
			Raw:      pb[:8],
		})
		pb = pb[8:]
	}
	return results
}
