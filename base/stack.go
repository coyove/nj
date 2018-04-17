package base

const (
	INIT_CAPACITY = 16
)

type Stack struct {
	data []Value
}

func NewStack() *Stack {
	return &Stack{
		data: make([]Value, 0, INIT_CAPACITY),
	}
}

func (s *Stack) grow(newSize int) {
	if newSize > cap(s.data) {
		old := s.data
		s.data = make([]Value, newSize, newSize*3/2)
		copy(s.data, old)
	}
	s.data = s.data[:newSize]
}

func (s *Stack) Size() int {
	return len(s.data)
}

func (s *Stack) Get(index int) Value {
	// if (shouldLock) rw.lock();
	// if (shouldLock) rw.unlock();
	if index >= len(s.data) {
		return NewValue()
	}

	return s.data[index]
}

func (s *Stack) Set(index int, value Value) {
	// if (shouldLock) rw.lock();

	if index >= len(s.data) {
		s.grow(index + 1)
	}

	s.data[index] = value
	// if (shouldLock) rw.unlock();
}

func (s *Stack) Add(value Value) {
	s.Set(len(s.data), value)
}

func (s *Stack) Clear() {
	s.data = s.data[:0]
}

func (s *Stack) InsertStack(index int, s2 *Stack) {
	s.Insert(index, s2.data)
}

func (s *Stack) Insert(index int, data []Value) {
	// if (shouldLock) rw.lock();
	if index <= len(s.data) {
		ln := len(s.data)
		s.grow(ln + len(data))
		copy(s.data[len(s.data)-(ln-index):], s.data[index:])
	} else {
		s.grow(index + len(data))
	}
	copy(s.data[index:], data)
	// if (shouldLock) rw.unlock();
}

func (s *Stack) Values() []Value {
	return s.data
}
