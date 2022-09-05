package bas

var ObjectProto Object

var Proto = struct {
	Bool            *Object
	Str             *Object
	Bytes           *Object
	Int             *Object
	Float           *Object
	Func            *Object
	Array           *Object
	Error           *Object
	Native          *Object
	NativeMap       *Object
	NativePtr       *Object
	NativeIntf      *Object
	Channel         *Object
	Reader          *NativeMeta
	Writer          *NativeMeta
	Closer          *NativeMeta
	ReadWriter      *NativeMeta
	ReadCloser      *NativeMeta
	WriteCloser     *NativeMeta
	ReadWriteCloser *NativeMeta
}{
	Bool:            NewObject(0),                                                // filled in lib_init.go
	Str:             NewObject(0),                                                // filled in lib_init.go
	Bytes:           NewObject(0),                                                // filled in lib_init.go
	Int:             NewObject(0),                                                // filled in lib_init.go
	Float:           NewObject(0),                                                // filled in lib_init.go
	Func:            NewObject(0),                                                // filled in lib_init.go
	Array:           NewObject(0),                                                // filled in lib_init.go
	Error:           NewObject(0),                                                // filled in lib_init.go
	Channel:         NewObject(0),                                                // filled in lib_init.go
	Native:          NewObject(0),                                                // filled in lib_init.go
	NativeMap:       NewObject(0),                                                // filled in lib_init.go
	NativePtr:       NewObject(0),                                                // filled in lib_init.go
	NativeIntf:      NewObject(0),                                                // filled in lib_init.go
	Reader:          newEmptyNativeMetaInternal("Reader", NewObject(0)),          // filled in io.go
	Writer:          newEmptyNativeMetaInternal("Writer", NewObject(0)),          // filled in io.go
	Closer:          newEmptyNativeMetaInternal("Closer", NewObject(0)),          // filled in io.go
	ReadWriter:      newEmptyNativeMetaInternal("ReadWriter", NewObject(0)),      // filled in io.go
	ReadCloser:      newEmptyNativeMetaInternal("ReadCloser", NewObject(0)),      // filled in io.go
	WriteCloser:     newEmptyNativeMetaInternal("WriteCloser", NewObject(0)),     // filled in io.go
	ReadWriteCloser: newEmptyNativeMetaInternal("ReadWriteCloser", NewObject(0)), // filled in io.go
}
