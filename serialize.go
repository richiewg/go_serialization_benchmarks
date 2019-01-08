package goserbench

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/big"
)

var (
	ErrorInvalidValue       = fmt.Errorf("invalid value")
	ErrorInvalidUint        = fmt.Errorf("invalid uint")
	ErrorInvalidUintRange   = fmt.Errorf("invalid uint range")
	ErrorDataSizeOutOfRange = fmt.Errorf("data size out of range")
)

//WriteUint64 serialize a uint64 value to writer.
//The serialized binary data with prefix 0, and in big-end byte order.
//If the value of uint64 <= 0x80, the serialized result is the value self and the size is 1byte;
//If the value > 0x80    and < 1 << 8,  the serialized result is 0x80+1 + value(without prefix 0 in big-end size) and the size is 2bytes;
//If the value > 1 << 8  and < 1 << 16, the serialized result is 0x80+2 + value(without prefix 0 in big-end size) and the size is 3bytes;
//If the value > 1 << 16 and < 1 << 24, the serialized result is 0x80+3 + value(without prefix 0 in big-end size) and the size is 4bytes;
//If the value > 1 << 24 and < 1 << 32, the serialized result is 0x80+4 + value(without prefix 0 in big-end size) and the size is 5bytes;
//If the value > 1 << 32 and < 1 << 40, the serialized result is 0x80+5 + value(without prefix 0 in big-end size) and the size is 6bytes;
//If the value > 1 << 40 and < 1 << 48, the serialized result is 0x80+6 + value(without prefix 0 in big-end size) and the size is 7bytes;
//If the value > 1 << 48 and < 1 << 56, the serialized result is 0x80+7 + value(without prefix 0 in big-end size) and the size is 8bytes;
//If the value > 1 << 56 and < 1 << 64, the serialized result is 0x80+8 + value(without prefix 0 in big-end size) and the size is 9bytes;
//Example:33        => [00100001]
//Example:33313     => [10000010 10000010 00100001]
//Example:1355849   => [10000011 00010100 10110000 1001001]
//Example:607826416 => [10000100 00100100 00111010 10110001 11110000]
func WriteUint64(w io.Writer, v uint64) error {
	var data [9]byte
	var size int
	switch {
	case v <= 0x80:
		data[0] = byte(v)
		size = 1
	case v < (1 << 8):
		data[0] = 0x80 + 1
		data[1] = byte(v)
		size = 2
	case v < (1 << 16):
		data[0] = 0x80 + 2
		data[1] = byte(v >> 8)
		data[2] = byte(v)
		size = 3
	case v < (1 << 24):
		data[0] = 0x80 + 3
		data[1] = byte(v >> 16)
		data[2] = byte(v >> 8)
		data[3] = byte(v)
		size = 4
	case v < (1 << 32):
		data[0] = 0x80 + 4
		data[1] = byte(v >> 24)
		data[2] = byte(v >> 16)
		data[3] = byte(v >> 8)
		data[4] = byte(v)
		size = 5
	case v < (1 << 40):
		data[0] = 0x80 + 5
		data[1] = byte(v >> 32)
		data[2] = byte(v >> 24)
		data[3] = byte(v >> 16)
		data[4] = byte(v >> 8)
		data[5] = byte(v)
		size = 6
	case v < (1 << 48):
		data[0] = 0x80 + 6
		data[1] = byte(v >> 40)
		data[2] = byte(v >> 32)
		data[3] = byte(v >> 24)
		data[4] = byte(v >> 16)
		data[5] = byte(v >> 8)
		data[6] = byte(v)
		size = 7
	case v < (1 << 56):
		data[0] = 0x80 + 7
		data[1] = byte(v >> 48)
		data[2] = byte(v >> 40)
		data[3] = byte(v >> 32)
		data[4] = byte(v >> 24)
		data[5] = byte(v >> 16)
		data[6] = byte(v >> 8)
		data[7] = byte(v)
		size = 8
	default:
		data[0] = 0x80 + 8
		data[1] = byte(v >> 56)
		data[2] = byte(v >> 48)
		data[3] = byte(v >> 40)
		data[4] = byte(v >> 32)
		data[5] = byte(v >> 24)
		data[6] = byte(v >> 16)
		data[7] = byte(v >> 8)
		data[8] = byte(v)
		size = 9
	}
	_, err := w.Write(data[0:size])
	return err
}

//ReadUint64 deserialize uint64 from a reader.
//The data from reader should in big-end byte order
func ReadUint64(r io.Reader) (uint64, error) {
	d := make([]byte, 1, 1)
	_, err := io.ReadFull(r, d)
	if err != nil {
		return 0, err
	}
	var p = uint64(d[0])
	if p <= 0x80 {
		return p, nil
	}
	var v uint64
	switch p - 0x80 {
	case 1:
		d = make([]byte, 1, 1)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(d[0])
	case 2:
		d = make([]byte, 2, 2)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(uint16(d[1]) | uint16(d[0])<<8)
	case 3:
		d = make([]byte, 3, 3)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(uint32(d[2]) | uint32(d[1])<<8 | uint32(d[0])<<16)
	case 4:
		d = make([]byte, 4, 4)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(uint32(d[3]) | uint32(d[2])<<8 | uint32(d[1])<<16 | uint32(d[0])<<24)
	case 5:
		d = make([]byte, 5, 5)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(d[4]) | uint64(d[3]<<8) | uint64(d[2])<<16 | uint64(d[1])<<24 |
			uint64(d[0])<<32
	case 6:
		d = make([]byte, 6, 6)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(d[5]) | uint64(d[4])<<8 | uint64(d[3])<<16 | uint64(d[2])<<24 |
			uint64(d[1])<<32 | uint64(d[0])<<40
	case 7:
		d = make([]byte, 7, 7)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(d[6]) | uint64(d[5])<<8 | uint64(d[4])<<16 | uint64(d[3])<<24 |
			uint64(d[2])<<32 | uint64(d[1])<<40 | uint64(d[0])<<48
	case 8:
		d = make([]byte, 8, 8)
		_, err = io.ReadFull(r, d)
		if err != nil {
			return 0, err
		}
		v = uint64(d[7]) | uint64(d[6])<<8 | uint64(d[5])<<16 | uint64(d[4])<<24 |
			uint64(d[3])<<32 | uint64(d[2])<<40 | uint64(d[1])<<48 | uint64(d[0])<<56
	default:
		return 0, ErrorInvalidUint
	}
	return v, nil
}

//WriteUint8 serialize uint8 to a writer.
//WriteUint8 is a convenient wrap function of WriteUint64
func WriteUint8(w io.Writer, v uint8) error {
	return WriteUint64(w, uint64(v))
}

//WriteUint16 serialize uint16 to a writer.
//WriteUint16 is a convenient wrap function of WriteUint64
func WriteUint16(w io.Writer, v uint16) error {
	return WriteUint64(w, uint64(v))
}

//WriteUint32 serialize uint32 to a writer.
//WriteUint32 is a convenient wrap function of WriteUint64
func WriteUint32(w io.Writer, v uint32) error {
	return WriteUint64(w, uint64(v))
}

//ReadUint8 deserialize uint8 from a reader.
//ReadUint8 is a convenient wrap function of ReadUint64
func ReadUint8(r io.Reader) (uint8, error) {
	v, err := ReadUint64(r)
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint8 {
		return 0, ErrorInvalidUintRange
	}
	return uint8(v), nil
}

//ReadUint16 deserialize uint16 from a reader.
//ReadUint16 is a convenient wrap function of ReadUint64
func ReadUint16(r io.Reader) (uint16, error) {
	v, err := ReadUint64(r)
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint16 {
		return 0, ErrorInvalidUintRange
	}
	return uint16(v), nil
}

//ReadUint32 deserialize uint32 from a reader.
//ReadUint32 is a convenient wrap function of ReadUint64
func ReadUint32(r io.Reader) (uint32, error) {
	v, err := ReadUint64(r)
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint32 {
		return 0, ErrorInvalidUintRange
	}
	return uint32(v), nil
}

//WriteBool serialize bool to a writer.
//True => 1;
//False => 0;
func WriteBool(w io.Writer, v bool) error {
	d := make([]byte, 1, 1)
	switch v {
	case true:
		d[0] = 1
	case false:
		d[0] = 0
	}
	_, err := w.Write(d)
	return err
}

//ReadBool deserialize bool from a reader.
//1 => True;
//0 => False;
func ReadBool(r io.Reader) (bool, error) {
	d := make([]byte, 1, 1)
	_, err := io.ReadFull(r, d)
	if err != nil {
		return false, err
	}
	switch uint8(d[0]) {
	case 1:
		return true, nil
	case 0:
		return false, nil
	default:
		return false, ErrorInvalidValue
	}
}

//WriteBytes serialize []byte to a writer.
//The serialize result is length([]byte) + []byte.
func WriteBytes(w io.Writer, data []byte) error {
	size := uint64(len(data))
	err := WriteUint64(w, size)
	if err != nil {
		return err
	}
	return WriteFixedBytes(w, data)
}

//WriteFixedBytes serialize []byte to a writer.
//The serialize result is []byte self without len([]byte)
//The max size of []byte is math.MaxUint32
func WriteFixedBytes(w io.Writer, data []byte) error {
	if len(data) >= math.MaxUint32 {
		return ErrorDataSizeOutOfRange
	}
	_, err := w.Write(data)
	return err
}

//ReadBytes deserialize []byte from a reader.
//The serialized data is len([]byte) + []byte
func ReadBytes(r io.Reader) ([]byte, error) {
	n, err := ReadUint32(r)
	if err != nil {
		return nil, err
	}
	return ReadFixedBytes(r, n)
}

//ReadFixedBytes deserialize []byte from a reader.
//The serialized data is []byte without len(data)
//The max size of []byte is math.MaxUint32
func ReadFixedBytes(r io.Reader, size uint32) ([]byte, error) {
	if size == 0 {
		return nil, nil
	}
	if size <= 1<<20 {
		d := make([]byte, size, size)
		_, err := io.ReadFull(r, d)
		if err != nil {
			return nil, err
		}
		return d, nil
	} else if size >= math.MaxUint32 {
		return nil, ErrorDataSizeOutOfRange
	}

	lr := io.LimitReader(r, int64(size))
	buf := &bytes.Buffer{}
	s, err := buf.ReadFrom(lr)
	if err != nil {
		return nil, err
	}
	if s != int64(size) {
		return nil, ErrorInvalidValue
	}
	return buf.Bytes(), nil
}

//WriteString serialize string to a write.
//WriteString is a convenient wrap function of WriteBytes
func WriteString(w io.Writer, data string) error {
	return WriteBytes(w, []byte(data))
}

//ReadString deserialize string from a reader.
//ReadString is a convenient wrap function of ReadBytes
func ReadString(r io.Reader) (string, error) {
	b, err := ReadBytes(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//WriteByte serialize byte to a writer.
func WriteByte(w io.Writer, v byte) error {
	_, err := w.Write([]byte{v})
	return err
}

//ReadByte deserialize byte from a reader.
func ReadByte(r io.Reader) (byte, error) {
	d := make([]byte, 1, 1)
	_, err := io.ReadFull(r, d)
	if err != nil {
		return 0, err
	}
	return d[0], nil
}

//WriteBigInt serialize big int to a writer.
//If the value of big int < 0, will return error.
func WriteBigInt(w io.Writer, v *big.Int) error {
	if v.Cmp(new(big.Int).SetUint64(0)) < 0 {
		return fmt.Errorf("cannot write negative big.Int")
	}
	return WriteBytes(w, v.Bytes())
}

//ReadBigInt deserialize big int from a reader.
func ReadBigInt(r io.Reader) (*big.Int, error) {
	d, err := ReadBytes(r)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(d), nil
}
