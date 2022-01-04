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
	StaticObject:        NamedObject("staticobject", 0),        // empty
	Bool:                NewObject(0),                          // filled in lib_init.go
	Str:                 NewObject(0),                          // filled in lib_init.go
	Bytes:               NewObject(0),                          // filled in lib_init.go
	Int:                 NewObject(0),                          // filled in lib_init.go
	Float:               NewObject(0),                          // filled in lib_init.go
	Func:                NewObject(0),                          // filled in lib_init.go
	Array:               NewObject(0),                          // filled in lib_init.go
	Error:               NewObject(0),                          // filled in lib_init.go
	Channel:             NewObject(0),                          // filled in lib_init.go
	Reader:              NamedObject("Reader", 0),              // filled in io.go
	Writer:              NamedObject("Writer", 0),              // filled in io.go
	Seeker:              NamedObject("Seeker", 0),              // filled in io.go
	Closer:              NamedObject("Closer", 0),              // filled in io.go
	ReadWriter:          NamedObject("ReadWriter", 0),          // filled in io.go
	ReadCloser:          NamedObject("ReadCloser", 0),          // filled in io.go
	WriteCloser:         NamedObject("WriteCloser", 0),         // filled in io.go
	ReadWriteCloser:     NamedObject("ReadWriteCloser", 0),     // filled in io.go
	ReadWriteSeekCloser: NamedObject("ReadWriteSeekCloser", 0), // filled in io.go
}
