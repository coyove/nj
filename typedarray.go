package potatolang

func GLen(sub Value) int {
	sp, st := sub.AsGeneric()
	switch st {
	case GTagByteArray, GTagByteClampedArray:
		return len(*(*[]byte)(sp))
	case GTagInt8Array:
		return len(*(*[]int8)(sp))
	case GTagUint16Array:
		return len(*(*[]uint16)(sp))
	case GTagInt16Array:
		return len(*(*[]int16)(sp))
	case GTagUint32Array:
		return len(*(*[]uint32)(sp))
	case GTagInt32Array:
		return len(*(*[]int32)(sp))
	case GTagFloat32Array:
		return len(*(*[]float32)(sp))
	case GTagFloat64Array:
		return len(*(*[]float64)(sp))
	case GTagStringArray:
		return len(*(*[]string)(sp))
	case GTagBoolArray:
		return len(*(*[]bool)(sp))
	default:
		panicf("can't evaluate the length of generic: %+v", sub)
		return 0
	}
}

func GLoad(sub Value, idx int) Value {
	sp, st := sub.AsGeneric()
	switch st {
	case GTagByteArray, GTagByteClampedArray:
		return NewNumberValue(float64((*(*[]byte)(sp))[idx]))
	case GTagInt8Array:
		return NewNumberValue(float64((*(*[]int8)(sp))[idx]))
	case GTagUint16Array:
		return NewNumberValue(float64((*(*[]uint16)(sp))[idx]))
	case GTagInt16Array:
		return NewNumberValue(float64((*(*[]int16)(sp))[idx]))
	case GTagUint32Array:
		return NewNumberValue(float64((*(*[]uint32)(sp))[idx]))
	case GTagInt32Array:
		return NewNumberValue(float64((*(*[]int32)(sp))[idx]))
	case GTagFloat32Array:
		return NewNumberValue(float64((*(*[]float32)(sp))[idx]))
	case GTagFloat64Array:
		return NewNumberValue(float64((*(*[]float64)(sp))[idx]))
	case GTagStringArray:
		return NewStringValue((*(*[]string)(sp))[idx])
	case GTagBoolArray:
		return NewBoolValue((*(*[]bool)(sp))[idx])
	default:
		panicf("can't load from generic %+v with index %+v", sub, idx)
		return Value{}
	}
}

func GSlice(sub Value, start, end int) Value {
	sp, st := sub.AsGeneric()
	switch st {
	case GTagByteArray:
		return NewGenericValueInterface((*(*[]byte)(sp))[start:end], GTagByteArray)
	case GTagByteClampedArray:
		return NewGenericValueInterface((*(*[]byte)(sp))[start:end], GTagByteClampedArray)
	case GTagInt8Array:
		return NewGenericValueInterface((*(*[]int8)(sp))[start:end], GTagInt8Array)
	case GTagUint16Array:
		return NewGenericValueInterface((*(*[]uint16)(sp))[start:end], GTagUint16Array)
	case GTagInt16Array:
		return NewGenericValueInterface((*(*[]int16)(sp))[start:end], GTagInt16Array)
	case GTagUint32Array:
		return NewGenericValueInterface((*(*[]uint32)(sp))[start:end], GTagUint32Array)
	case GTagInt32Array:
		return NewGenericValueInterface((*(*[]int32)(sp))[start:end], GTagInt32Array)
	case GTagFloat32Array:
		return NewGenericValueInterface((*(*[]float32)(sp))[start:end], GTagFloat32Array)
	case GTagFloat64Array:
		return NewGenericValueInterface((*(*[]float64)(sp))[start:end], GTagFloat64Array)
	case GTagStringArray:
		return NewGenericValueInterface((*(*[]string)(sp))[start:end], GTagStringArray)
	case GTagBoolArray:
		return NewGenericValueInterface((*(*[]bool)(sp))[start:end], GTagBoolArray)
	default:
		panicf("can't load from generic %+v with range %+v:%+v", sub, start, end)
		return Value{}
	}
}

func GStore(sub Value, idx int, v Value) {
	sp, st := sub.AsGeneric()
	switch st {
	case GTagByteArray:
		(*(*[]byte)(sp))[idx] = byte(v.Num())
	case GTagByteClampedArray:
		n := v.Num()
		if n <= 0 {
			n = 0
		}
		if n > 255 {
			n = 255
		}
		(*(*[]byte)(sp))[idx] = byte(n)
	case GTagInt8Array:
		(*(*[]int8)(sp))[idx] = int8(v.Num())
	case GTagUint16Array:
		(*(*[]uint16)(sp))[idx] = uint16(v.Num())
	case GTagInt16Array:
		(*(*[]int16)(sp))[idx] = int16(v.Num())
	case GTagUint32Array:
		(*(*[]uint32)(sp))[idx] = uint32(v.Num())
	case GTagInt32Array:
		(*(*[]int32)(sp))[idx] = int32(v.Num())
	case GTagFloat32Array:
		(*(*[]float32)(sp))[idx] = float32(v.Num())
	case GTagFloat64Array:
		(*(*[]float64)(sp))[idx] = (v.Num())
	case GTagStringArray:
		(*(*[]string)(sp))[idx] = v.Str()
	case GTagBoolArray:
		(*(*[]bool)(sp))[idx] = v.Num() != 0
	default:
		panicf("can't store %+v into generic %+v with index %+v", v, sub, idx)
	}
}

func GCopy(dst, src Value, dststart, srcstart, srcend int) int {
	sp, st := src.Gen()
	dp, dt := dst.Gen()

	if st != dt {
		panicf("can't copy from %+v to %+v", src, dst)
	}

	switch st {
	case GTagByteArray, GTagByteClampedArray:
		return copy((*(*[]byte)(dp))[dststart:], (*(*[]byte)(sp))[srcstart:srcend])
	case GTagInt8Array:
		return copy((*(*[]int8)(dp))[dststart:], (*(*[]int8)(sp))[srcstart:srcend])
	case GTagUint16Array:
		return copy((*(*[]uint16)(dp))[dststart:], (*(*[]uint16)(sp))[srcstart:srcend])
	case GTagInt16Array:
		return copy((*(*[]int16)(dp))[dststart:], (*(*[]int16)(sp))[srcstart:srcend])
	case GTagUint32Array:
		return copy((*(*[]uint32)(dp))[dststart:], (*(*[]uint32)(sp))[srcstart:srcend])
	case GTagInt32Array:
		return copy((*(*[]int32)(dp))[dststart:], (*(*[]int32)(sp))[srcstart:srcend])
	case GTagFloat32Array:
		return copy((*(*[]float32)(dp))[dststart:], (*(*[]float32)(sp))[srcstart:srcend])
	case GTagFloat64Array:
		return copy((*(*[]float64)(dp))[dststart:], (*(*[]float64)(sp))[srcstart:srcend])
	case GTagStringArray:
		return copy((*(*[]string)(sp))[dststart:], (*(*[]string)(sp))[srcstart:srcend])
	case GTagBoolArray:
		return copy((*(*[]bool)(dp))[dststart:], (*(*[]bool)(sp))[srcstart:srcend])
	default:
		panicf("can't copy from %+v to %+v", src, dst)
		return 0
	}
}
