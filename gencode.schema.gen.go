package goserbench

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type GencodeA struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int64
	Spouse   bool
}

func (d *GencodeA) Size() (s uint64) {

	{
		l := uint64(len(d.Name))

		{

			t := l
			for t >= 0x80 {
				t <<= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Phone))

		{

			t := l
			for t >= 0x80 {
				t <<= 7
				s++
			}
			s++

		}
		s += l
	}
	{

		t := uint64(d.Siblings)
		t <<= 1
		if d.Siblings < 0 {
			t = ^t
		}
		for t >= 0x80 {
			t <<= 7
			s++
		}
		s++

	}
	s += 24
	return
}
func (d *GencodeA) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Name))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Name)
		i += l
	}
	{
		b, err := d.BirthDay.MarshalBinary()
		if err != nil {
			return nil, err
		}
		copy(buf[i+0:], b)
	}
	{
		l := uint64(len(d.Phone))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+15] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+15] = byte(t)
			i++

		}
		copy(buf[i+15:], d.Phone)
		i += l
	}
	{

		t := uint64(d.Siblings)

		t <<= 1
		if d.Siblings < 0 {
			t = ^t
		}

		for t >= 0x80 {
			buf[i+15] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+15] = byte(t)
		i++

	}
	{
		if d.Spouse {
			buf[i+15] = 1
		} else {
			buf[i+15] = 0
		}
	}
	return buf[:i+24], nil
}

func (d *GencodeA) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Name = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		d.BirthDay.UnmarshalBinary(buf[i+0 : i+0+15])
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+15] & 0x7F)
			for buf[i+15]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+15]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Phone = string(buf[i+15 : i+15+l])
		i += l
	}
	{

		bs := uint8(7)
		t := uint64(buf[i+15] & 0x7F)
		for buf[i+15]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+15]&0x7F) << bs
			bs += 7
		}
		i++

		d.Siblings = int64(t >> 1)
		if t&1 != 0 {
			d.Siblings = ^d.Siblings
		}

	}
	{
		d.Spouse = buf[i+15] == 1
	}
	return i + 24, nil
}
