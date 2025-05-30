package fuzzy

import (
	"encoding/json"
	"io"
	"math"
	"strings"
	"unsafe"

	"github.com/json-iterator/go"
	"github.com/samber/lo"
)

func RegisterPointerFuzzyDecoders() {
	jsoniter.RegisterTypeDecoder("bool", &fuzzyBoolDecoder{})
	jsoniter.RegisterTypeDecoder("*bool", &fuzzyPointerBoolDecoder{})
	jsoniter.RegisterTypeDecoder("*string", &fuzzyPointerStringDecoder{})
	jsoniter.RegisterTypeDecoder("*float32", &fuzzyPointerFloat32Decoder{})
	jsoniter.RegisterTypeDecoder("*float64", &fuzzyPointerFloat64Decoder{})
	jsoniter.RegisterTypeDecoder("*int", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxInt) || val < float64(minInt) {
				iter.ReportError("fuzzy decode *int", "exceed range")
				return
			}
			*((**int)(ptr)) = lo.ToPtr(int(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int)(ptr)) = lo.ToPtr(iter.ReadInt())
			}
		}
	}})

	jsoniter.RegisterTypeDecoder("*uint", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxUint) || val < 0 {
				iter.ReportError("fuzzy decode *uint", "exceed range")
				return
			}
			*((**uint)(ptr)) = lo.ToPtr(uint(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint)(ptr)) = lo.ToPtr(iter.ReadUint())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int8", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt8) || val < float64(math.MinInt8) {
				iter.ReportError("fuzzy decode *int8", "exceed range")
				return
			}
			*((**int8)(ptr)) = lo.ToPtr(int8(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int8)(ptr)) = lo.ToPtr(iter.ReadInt8())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint8", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint8) || val < 0 {
				iter.ReportError("fuzzy decode *uint8", "exceed range")
				return
			}
			*((*uint8)(ptr)) = uint8(val)
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint8)(ptr)) = lo.ToPtr(iter.ReadUint8())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int16", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt16) || val < float64(math.MinInt16) {
				iter.ReportError("fuzzy decode *int16", "exceed range")
				return
			}
			*((**uint16)(ptr)) = lo.ToPtr(uint16(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int16)(ptr)) = lo.ToPtr(iter.ReadInt16())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint16", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint16) || val < 0 {
				iter.ReportError("fuzzy decode *uint16", "exceed range")
				return
			}
			*((**uint16)(ptr)) = lo.ToPtr(uint16(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint16)(ptr)) = lo.ToPtr(iter.ReadUint16())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int32", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt32) || val < float64(math.MinInt32) {
				iter.ReportError("fuzzy decode *int32", "exceed range")
				return
			}
			*((**int32)(ptr)) = lo.ToPtr(int32(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int32)(ptr)) = lo.ToPtr(iter.ReadInt32())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint32", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint32) || val < 0 {
				iter.ReportError("fuzzy decode *uint32", "exceed range")
				return
			}
			*((**uint32)(ptr)) = lo.ToPtr(uint32(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint32)(ptr)) = lo.ToPtr(iter.ReadUint32())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int64", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt64) || val < float64(math.MinInt64) {
				iter.ReportError("fuzzy decode *int64", "exceed range")
				return
			}
			*((**int64)(ptr)) = lo.ToPtr(int64(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int64)(ptr)) = lo.ToPtr(iter.ReadInt64())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint64", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint64) || val < 0 {
				iter.ReportError("fuzzy decode *uint64", "exceed range")
				return
			}
			*((**uint64)(ptr)) = lo.ToPtr(uint64(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint64)(ptr)) = lo.ToPtr(iter.ReadUint64())
			}
		}
	}})
}

type fuzzyBoolDecoder struct{}

type fuzzyPointerBoolDecoder struct{}

func (f fuzzyPointerBoolDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.BoolValue:
		*((**bool)(ptr)) = lo.ToPtr(iter.ReadBool())
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		if number == "1" {
			*((**bool)(ptr)) = lo.ToPtr(true)
		} else {
			*((**bool)(ptr)) = lo.ToPtr(false)
		}
	case jsoniter.StringValue:
		value := iter.ReadString()
		switch value {
		case "1", "true":
			*((**bool)(ptr)) = lo.ToPtr(true)
		case "0", "false":
			*((**bool)(ptr)) = lo.ToPtr(false)
		default:
			iter.ReportError("fuzzyPointerBoolDecoder", "not bool")
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((**bool)(ptr)) = nil
	default:
		iter.ReportError("fuzzyPointerBoolDecoder", "not number or string")
	}
}

func (f fuzzyBoolDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.BoolValue:
		*((*bool)(ptr)) = iter.ReadBool()
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		if number == "1" {
			*((*bool)(ptr)) = true
		} else {
			*((*bool)(ptr)) = false
		}
	case jsoniter.StringValue:
		value := iter.ReadString()
		switch value {
		case "1", "true":
			*((*bool)(ptr)) = true
		case "0", "false":
			*((*bool)(ptr)) = false
		default:
			iter.ReportError("fuzzyBoolDecoder", "not bool")
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((*bool)(ptr)) = false
	default:
		iter.ReportError("fuzzyBoolDecoder", "not number or string")
	}
}

type fuzzyPointerStringDecoder struct{}

func (decoder *fuzzyPointerStringDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		*((**string)(ptr)) = lo.ToPtr(string(number))
	case jsoniter.StringValue:
		if EnableXssProtection {
			*((**string)(ptr)) = lo.ToPtr(BlueMondayPolicy.Sanitize(iter.ReadString()))
		} else {
			*((**string)(ptr)) = lo.ToPtr(iter.ReadString())
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((**string)(ptr)) = nil
	default:
		iter.ReportError("fuzzyStringDecoder", "not number or string")
	}
}

type fuzzyPointerIntegerDecoder struct {
	fun func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator)
}

func (decoder *fuzzyPointerIntegerDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		str = string(number)
	case jsoniter.StringValue:
		str = iter.ReadString()
		if str == "" {
			str = "null"
		}
	case jsoniter.BoolValue:
		if iter.ReadBool() {
			str = "1"
		} else {
			str = "0"
		}
	case jsoniter.NilValue:
		iter.Skip()
		str = "null"
	default:
		iter.ReportError("fuzzyPointerIntegerDecoder", "not number or string")
	}

	newIter := iter.Pool().BorrowIterator([]byte(str))
	defer iter.Pool().ReturnIterator(newIter)
	isFloat := strings.IndexByte(str, '.') != -1
	decoder.fun(isFloat, ptr, newIter)
	if newIter.Error != nil && newIter.Error != io.EOF {
		iter.Error = newIter.Error
	}
}

type fuzzyPointerFloat32Decoder struct{}

func (decoder *fuzzyPointerFloat32Decoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		*((**float32)(ptr)) = lo.ToPtr(iter.ReadFloat32())
	case jsoniter.StringValue:
		str = iter.ReadString()
		newIter := iter.Pool().BorrowIterator([]byte(str))
		defer iter.Pool().ReturnIterator(newIter)
		*((**float32)(ptr)) = lo.ToPtr(newIter.ReadFloat32())
		if newIter.Error != nil && newIter.Error != io.EOF {
			iter.Error = newIter.Error
		}
	case jsoniter.BoolValue:
		// support bool to float32
		if iter.ReadBool() {
			*((**float32)(ptr)) = lo.ToPtr(float32(1))
		} else {
			*((**float32)(ptr)) = lo.ToPtr(float32(0))
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((**float32)(ptr)) = nil
	default:
		iter.ReportError("fuzzyPointerFloat32Decoder", "not number or string")
	}
}

type fuzzyPointerFloat64Decoder struct{}

func (decoder *fuzzyPointerFloat64Decoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		*((**float64)(ptr)) = lo.ToPtr(iter.ReadFloat64())
	case jsoniter.StringValue:
		str = iter.ReadString()
		newIter := iter.Pool().BorrowIterator([]byte(str))
		defer iter.Pool().ReturnIterator(newIter)
		*((**float64)(ptr)) = lo.ToPtr(newIter.ReadFloat64())
		if newIter.Error != nil && newIter.Error != io.EOF {
			iter.Error = newIter.Error
		}
	case jsoniter.BoolValue:
		// support bool to float64
		if iter.ReadBool() {
			*((**float64)(ptr)) = lo.ToPtr(float64(1))
		} else {
			*((**float64)(ptr)) = lo.ToPtr(float64(0))
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((**float64)(ptr)) = nil
	default:
		iter.ReportError("fuzzyFloat64Decoder", "not number or string")
	}
}
