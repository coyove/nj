package potatolang

func LenGeneric(sub Value) int {
	switch sub.a {
	case 1, 11:
		return len(*(*[]byte)(sub.ptr))
	case 2:
		return len(*(*[]int8)(sub.ptr))
	case 3:
		return len(*(*[]uint16)(sub.ptr))
	case 4:
		return len(*(*[]int16)(sub.ptr))
	case 5:
		return len(*(*[]uint32)(sub.ptr))
	case 6:
		return len(*(*[]int32)(sub.ptr))
	case 7:
		return len(*(*[]float32)(sub.ptr))
	case 8:
		return len(*(*[]float64)(sub.ptr))
	case 9:
		return len(*(*[]string)(sub.ptr))
	case 10:
		return len(*(*[]bool)(sub.ptr))
	default:
		panicf("can't evaluate the length of generic: %+v", sub)
		return 0
	}
}

func LoadFromGeneric(sub Value, idx int) Value {
	switch sub.a {
	case 1, 11:
		return NewNumberValue(float64((*(*[]byte)(sub.ptr))[idx]))
	case 2:
		return NewNumberValue(float64((*(*[]int8)(sub.ptr))[idx]))
	case 3:
		return NewNumberValue(float64((*(*[]uint16)(sub.ptr))[idx]))
	case 4:
		return NewNumberValue(float64((*(*[]int16)(sub.ptr))[idx]))
	case 5:
		return NewNumberValue(float64((*(*[]uint32)(sub.ptr))[idx]))
	case 6:
		return NewNumberValue(float64((*(*[]int32)(sub.ptr))[idx]))
	case 7:
		return NewNumberValue(float64((*(*[]float32)(sub.ptr))[idx]))
	case 8:
		return NewNumberValue(float64((*(*[]float64)(sub.ptr))[idx]))
	case 9:
		return NewStringValue((*(*[]string)(sub.ptr))[idx])
	case 10:
		return NewBoolValue((*(*[]bool)(sub.ptr))[idx])
	default:
		panicf("can't load from generic %+v with index %+v", sub, idx)
		return Value{}
	}
}

func StoreToGeneric(sub Value, idx int, v Value) {
	switch sub.a {
	case 1:
		(*(*[]byte)(sub.ptr))[idx] = byte(v.Num())
	case 11:
		n := v.Num()
		if n <= 0 {
			n = 0
		}
		if n > 255 {
			n = 255
		}
		(*(*[]byte)(sub.ptr))[idx] = byte(n)
	case 2:
		(*(*[]int8)(sub.ptr))[idx] = int8(v.Num())
	case 3:
		(*(*[]uint16)(sub.ptr))[idx] = uint16(v.Num())
	case 4:
		(*(*[]int16)(sub.ptr))[idx] = int16(v.Num())
	case 5:
		(*(*[]uint32)(sub.ptr))[idx] = uint32(v.Num())
	case 6:
		(*(*[]int32)(sub.ptr))[idx] = int32(v.Num())
	case 7:
		(*(*[]float32)(sub.ptr))[idx] = float32(v.Num())
	case 8:
		(*(*[]float64)(sub.ptr))[idx] = (v.Num())
	case 9:
		(*(*[]string)(sub.ptr))[idx] = v.Str()
	case 10:
		(*(*[]bool)(sub.ptr))[idx] = v.Num() != 0
	default:
		panicf("can't store %+v into generic %+v with index %+v", v, sub, idx)
	}
}
