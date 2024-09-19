package huacache

type ByteView struct {
	B []byte
}

func (v ByteView) Len() int {
	return len(v.B)
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.B)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.B)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
