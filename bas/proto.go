package bas

var ObjectProto Object

var Proto = struct {
	StaticObject,
	Bool,
	Str,
	Bytes,
	Int,
	Float,
	Func,
	Array,
	Error,
	Native,
	NativeMap,
	NativePtr,
	NativeIntf,
	Channel,
	Reader,
	Writer,
	Seeker,
	Closer,
	ReadWriter,
	ReadCloser,
	WriteCloser,
	ReadWriteCloser,
	ReadWriteSeekCloser *Object
}{
	StaticObject:        NewNamedObject("staticobject", 0),        // empty
	Bool:                NewObject(0),                             // filled in lib_init.go
	Str:                 NewObject(0),                             // filled in lib_init.go
	Bytes:               NewObject(0),                             // filled in lib_init.go
	Int:                 NewObject(0),                             // filled in lib_init.go
	Float:               NewObject(0),                             // filled in lib_init.go
	Func:                NewObject(0),                             // filled in lib_init.go
	Array:               NewObject(0),                             // filled in lib_init.go
	Error:               NewObject(0),                             // filled in lib_init.go
	Channel:             NewObject(0),                             // filled in lib_init.go
	Native:              NewObject(0),                             // filled in lib_init.go
	NativeMap:           NewObject(0),                             // filled in lib_init.go
	NativePtr:           NewObject(0),                             // filled in lib_init.go
	NativeIntf:          NewObject(0),                             // filled in lib_init.go
	Reader:              NewNamedObject("Reader", 0),              // filled in io.go
	Writer:              NewNamedObject("Writer", 0),              // filled in io.go
	Seeker:              NewNamedObject("Seeker", 0),              // filled in io.go
	Closer:              NewNamedObject("Closer", 0),              // filled in io.go
	ReadWriter:          NewNamedObject("ReadWriter", 0),          // filled in io.go
	ReadCloser:          NewNamedObject("ReadCloser", 0),          // filled in io.go
	WriteCloser:         NewNamedObject("WriteCloser", 0),         // filled in io.go
	ReadWriteCloser:     NewNamedObject("ReadWriteCloser", 0),     // filled in io.go
	ReadWriteSeekCloser: NewNamedObject("ReadWriteSeekCloser", 0), // filled in io.go
}
