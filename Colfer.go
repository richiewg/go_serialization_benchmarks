package goserbench

// This file was generated by colf(1); DO NOT EDIT

import (
	"fmt"
	"io"
	"math"
	"time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = math.E
var _ = time.RFC3339

// ColferContinue signals a data continuation as a byte index.
type ColferContinue int

func (i ColferContinue) Error() string {
	return fmt.Sprintf("colfer: data continuation at byte %d", i)
}

// ColferError signals a data mismatch as as a byte index.
type ColferError int

func (i ColferError) Error() string {
	return fmt.Sprintf("colfer: unknown header at byte %d", i)
}

type ColferA struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int32
	Spouse   bool
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
func (o *ColferA) MarshalTo(buf []byte) int {
	if o == nil {
		return 0
	}

	buf[0] = 0x80
	i := 1

	if v := o.Name; len(v) != 0 {
		buf[i] = 0
		i++
		x := uint(len(v))
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		to := i + len(v)
		copy(buf[i:], v)
		i = to
	}

	if v := o.BirthDay; !v.IsZero() {
		buf[i] = 1
		s, ns := v.Unix(), v.Nanosecond()
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(s>>56), byte(s>>48), byte(s>>40), byte(s>>32)
		buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
		if ns == 0 {
			i += 9
		} else {
			buf[i] |= 0x80
			buf[i+9], buf[i+10], buf[i+11], buf[i+12] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 13
		}
	}

	if v := o.Phone; len(v) != 0 {
		buf[i] = 2
		i++
		x := uint(len(v))
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		to := i + len(v)
		copy(buf[i:], v)
		i = to
	}

	if v := o.Siblings; v != 0 {
		x := uint32(v)
		if v >= 0 {
			buf[i] = 3
		} else {
			x = ^x + 1
			buf[i] = 3 | 0x80
		}
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if o.Spouse {
		buf[i] = 4
		i++
	}


	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
func (o *ColferA) MarshalLen() int {
	if o == nil {
		return 0
	}

	l := 2

	if x := len(o.Name); x != 0 {
		l += x
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
	}

	if v := o.BirthDay; !v.IsZero() {
		if v.Nanosecond() == 0 {
			l += 9
		} else {
			l += 13
		}
	}

	if x := len(o.Phone); x != 0 {
		l += x
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
	}

	if v := o.Siblings; v != 0 {
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
	}

	if o.Spouse {
		l++
	}


	return l
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return is always nil.
func (o *ColferA) MarshalBinary() (data []byte, err error) {
	data = make([]byte, o.MarshalLen())
	o.MarshalTo(data)
	return data, nil
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, goserbench.ColferError, and goserbench.ColferContinue.
func (o *ColferA) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	if data[0] != 0x80 {
		return ColferError(0)
	}

	if len(data) == 1 {
		return io.EOF
	}
	header := data[1]
	i := 2

	if header == 0 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if i == len(data) {
				return io.EOF
			}
			b := data[i]
			i++
			if shift == 28 {
				x |= uint32(b) << 28
				break
			}
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
		}

		to := i + int(x)
		if to >= len(data) {
			return io.EOF
		}
		o.Name = string(data[i:to])

		header = data[to]
		i = to + 1
	}

	if header == 1 {
		if i+8 >= len(data) {
			return io.EOF
		}
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		o.BirthDay = time.Unix(int64(sec), 0)

		header = data[i+8]
		i += 9
	} else if header == 1|0x80 {
		if i+12 >= len(data) {
			return io.EOF
		}
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		nsec := int64(uint(data[i+8])<<24 | uint(data[i+9])<<16 | uint(data[i+10])<<8 | uint(data[i+11]))
		o.BirthDay = time.Unix(int64(sec), nsec)

		header = data[i+12]
		i += 13
	}

	if header == 2 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if i == len(data) {
				return io.EOF
			}
			b := data[i]
			i++
			if shift == 28 {
				x |= uint32(b) << 28
				break
			}
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
		}

		to := i + int(x)
		if to >= len(data) {
			return io.EOF
		}
		o.Phone = string(data[i:to])

		header = data[to]
		i = to + 1
	}

	if header == 3 || header == 3|0x80 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if i == len(data) {
				return io.EOF
			}
			b := data[i]
			i++
			if shift == 28 {
				x |= uint32(b) << 28
				break
			}
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
		}
		if header&0x80 != 0 {
			x = ^x + 1
		}
		o.Siblings = int32(x)

		if i == len(data) {
			return io.EOF
		}
		header = data[i]
		i++
	}

	if header == 4 {
		o.Spouse = true

		if i == len(data) {
			return io.EOF
		}
		header = data[i]
		i++
	}

	if header == 5 {
		if i+8 >= len(data) {
			return io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])

		header = data[i+8]
		i += 9
	}

	if header != 0x7f {
		return ColferError(i - 1)
	}
	if i != len(data) {
		return ColferContinue(i)
	}
	return nil
}
