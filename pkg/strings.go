package pkg

import (
	"errors"
	"unsafe"
)

type MDString struct {
	addr   *MDString
	str    []rune
	Length int
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (m *MDString) copyCheck() {
	if m.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		m.addr = (*MDString)(noescape(unsafe.Pointer(m)))
	} else if m.addr != m {
		panic("strings: illegal use of non-zero MDString copied by value")
	}
}

// Init initializes the instance with the string s.
func (m *MDString) Init(s string) error {
	m.Reset()

	runes := []rune(s)
	return m.WriteRunes(runes)
}

// String returns the accumulated string.
func (m *MDString) String() string {
	return *(*string)(unsafe.Pointer(&m.str))
}

//Len returns the number of accumulated runes (characters)
// O(1)
func (m *MDString) Len() int { return len(m.str) }

//Len2 returns the number of accumulated runes (characters)
// O(1)
func (m *MDString) Len2() int { return m.Length }

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (m *MDString) Cap() int { return cap(m.str) }

// Reset resets the Builder to be empty.
func (m *MDString) Reset() {
	m.addr = nil
	m.str = nil
	m.Length = 0
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// runes of capacity beyond len(m.str).
func (m *MDString) grow(n int) {
	str := make([]rune, len(m.str), 2*cap(m.str)+n)
	copy(str, m.str)
	m.str = str
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n runes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (m *MDString) Grow(n int) {
	m.copyCheck()
	if n < 0 {
		panic("MDString.Grow: negative count")
	}
	if cap(m.str)-len(m.str) < n {
		m.grow(n)
	}
}

func (m *MDString) WriteString(s string) error {
	runes := []rune(s)
	return m.WriteRunes(runes)
}

// WriteRunes appends the UTF-8 encoding of Unicode code point r to str.
func (m *MDString) WriteRunes(r []rune) error {
	m.copyCheck()
	n := len(r)

	m.str = append(m.str, r...)
	m.Length += n
	return nil
}

func (m *MDString) Slice(start int, end int) string {
	if end <= start {
		return ""
	}

	n := m.Length
	if start >= n || end > n {
		return ""
	}
	return string(m.str[start:end])
}

func (m *MDString) CharCodeAt(idx int) (rune, error) {
	if idx < 0 || idx >= m.Length {
		return 0x00, errors.New("idx is out of bounds")
	}
	return m.str[idx], nil
}

func (m *MDString) CharAt(idx int) (string, error) {
	ch, err := m.CharCodeAt(idx)

	if err != nil {
		return "", err
	}

	return string(ch), nil
}
