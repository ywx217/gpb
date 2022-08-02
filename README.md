<p align="center">
<a href="https://godoc.org/github.com/ywx217/gpb"><img src="https://img.shields.io/badge/api-reference-blue.svg?style=flat-square" alt="GoDoc"></a>
</p>

GPB is a Go package inspired by [tidwall/gjson](https://github.com/tidwall/gjson), that provides a fast and simple way
to get fields from a protobuf binary message.

Getting Started
===============

## Install

Install the package by simply run `go get`:

```shell
$ go get -u github.com/ywx217/gpb
```

## Performance

Benchmarks of GPB alongside [golang/protobuf](https://github.com/golang/protobuf) is in [gpb_test.go](./gpb_test.go),
each operation gets a single int32 value from the encoded protobuf data.

These benchmarks were run on a MacBook Pro 15" Intel Core i7@2.20GHz

```shell
➜  gpb git:(master) ✗ make bench-compare 
name                time/op
GoProtobufTiny-12    140ns ± 2%
GpbTiny-12          42.9ns ± 8%
GoProtobufSmall-12  3.82µs ± 4%
GpbSmall-12         77.8ns ± 4%

name                alloc/op
GoProtobufTiny-12    52.0B ± 0%
GpbTiny-12           0.00B     
GoProtobufSmall-12  2.50kB ± 0%
GpbSmall-12          0.00B     

name                allocs/op
GoProtobufTiny-12     2.00 ± 0%
GpbTiny-12            0.00     
GoProtobufSmall-12    89.0 ± 0%
GpbSmall-12           0.00     
```

### Test data

#### Tidy payload, size=2B

```text
GoEnum {
    foo:FOO1
}
```

#### Small payload, size=451B

```text
GoTest {
    Kind: TIME
    RequiredField: {
      Label: "label"
      Type: "type"
    }
    RepeatedField: {
      Label: "label"
      Type: "type"
    }
    RepeatedField: {
      Label: "label"
      Type: "type"
    }
    F_Bool_required: true
    F_Int32_required: 3
    F_Int64_required: 6
    F_Fixed32_required: 32
    F_Fixed64_required: 64
    F_Uint32_required: 3232
    F_Uint64_required: 6464
    F_Float_required: 3232
    F_Double_required: 6464
    F_String_required: "string"
    F_Bytes_required: "bytes"
    F_Sint32_required: -32
    F_Sint64_required: -64
    F_Sfixed32_required: -32
    F_Sfixed64_required: -64
    F_Bool_repeated: false
    F_Bool_repeated: true
    F_Int32_repeated: 32
    F_Int32_repeated: 33
    F_Int64_repeated: 64
    F_Int64_repeated: 65
    F_Fixed32_repeated: 3232
    F_Fixed32_repeated: 3333
    F_Fixed64_repeated: 6464
    F_Fixed64_repeated: 6565
    F_Uint32_repeated: 323232
    F_Uint32_repeated: 333333
    F_Uint64_repeated: 646464
    F_Uint64_repeated: 656565
    F_Float_repeated: 32
    F_Float_repeated: 33
    F_Double_repeated: 64
    F_Double_repeated: 65
    F_String_repeated: "hello"
    F_String_repeated: "sailor"
    F_Bytes_repeated: "big"
    F_Bytes_repeated: "nose"
    F_Sint32_repeated: 32
    F_Sint32_repeated: -32
    F_Sint64_repeated: 64
    F_Sint64_repeated: -64
    F_Sfixed32_repeated: 32
    F_Sfixed32_repeated: -32
    F_Sfixed64_repeated: 64
    F_Sfixed64_repeated: -64
    F_Bool_defaulted: true
    F_Int32_defaulted: 32
    F_Int64_defaulted: 64
    F_Fixed32_defaulted: 320
    F_Fixed64_defaulted: 640
    F_Uint32_defaulted: 3200
    F_Uint64_defaulted: 6400
    F_Float_defaulted: 314159
    F_Double_defaulted: 271828
    F_String_defaulted: "hello, \"world!\"\n"
    F_Bytes_defaulted: "Bignose"
    F_Sint32_defaulted: -32
    F_Sint64_defaulted: -64
    F_Sfixed32_defaulted: -32
    F_Sfixed64_defaulted: -64
    RequiredGroup: {
      RequiredField: "required"
    }
    RepeatedGroup: {
      RequiredField: "repeated"
    }
    RepeatedGroup: {
      RequiredField: "repeated"
    }
}
```
