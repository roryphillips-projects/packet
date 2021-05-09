package packet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// Packet is allows the reading of properties out of raw binary formats
// Each method on this interface when called will scan forwards through the bytes specified.
type Packet interface {
	// Bool read a bool value from the packet
	// Advances a byte with the packet
	// Errors if the value cannot be converted to boolean
	// Errors if there are no bytes remaining to be read
	Bool() (bool, error)
	// Float32 read a float32 value from the packet
	// Advances 4 bytes with the packet
	// Errors if the value cannot be converted to a float32
	// Errors if there are less then 4 bytes remaining to be read
	Float32() (float32, error)
	// Uint64 read an uint64 value from the packet
	// Advances 8 bytes with the packet
	// Errors if the value cannot be converted to an uint64
	// Errors if there are less than 8 bytes remaining to be read
	Uint64() (uint64, error)
	// Uint32 read an uint32 value from the packet
	// Advances 4 bytes with the packet
	// Errors if the value cannot be converted to an uint32
	// Errors if there are less than 4 bytes remaining to be read
	Uint32() (uint32, error)
	// Uint16 read an uint16 value from the packet
	// Advances 2 bytes with the packet
	// Errors if the value cannot be converted to an uint16
	// Errors if there are less than 2 bytes remaining to be read
	Uint16() (uint16, error)
	// Uint8 read an uint8 value from the packet
	// Advances a byte with the packet
	// Errors if the value cannot be converted to an uint8
	// Errors if there are no bytes remaining to be read
	Uint8() (uint8, error)
	// Int64 read an int64 value from the packet
	// Advances 8 bytes with the packet
	// Errors if the value cannot be converted to an int64
	// Errors if there are less than 8 bytes remaining to be read
	Int64() (int64, error)
	// Int32 read an int32 value from the packet
	// Advances 4 bytes with the packet
	// Errors if the value cannot be converted to an int32
	// Errors if there are less then 4 bytes remaining to be read
	Int32() (int32, error)
	// Int16 read an int16 value from the packet
	// Advances 2 bytes with the packet
	// Errors if the value cannot be converted to an int16
	// Errors if there are less then 2 bytes remaining to be read
	Int16() (int16, error)
	// Int8 read an int8 value from the packet
	// Advances a byte with the packet
	// Errors if the value cannot be converted to an int8
	// Errors if there are no bytes remaining to be read
	Int8() (int8, error)
}

// NewPacket returns a Packet using the byte slice input
func NewPacket(data []byte) Packet {
	return &packet{
		reader: bufio.NewReader(bytes.NewReader(data)),
	}
}

type packet struct {
	reader *bufio.Reader
}

func (p *packet) read() (byte, error) {
	return p.reader.ReadByte()
}

func (p *packet) readN(n int) ([]byte, error) {
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		b, err :=  p.read()
		if err != nil {
			return nil, fmt.Errorf("unable to read %v bytes: %v", n, err)
		}
		out[i] = b
	}

	return out, nil
}

// Bool read a bool value from the packet
// Advances a byte with the packet
// Errors if the value cannot be converted to boolean
// Errors if there are no bytes remaining to be read
func (p *packet) Bool() (bool, error) {
	b, err := p.read()
	if err != nil {
		return false, fmt.Errorf("failed to read byte: %v", err)
	}
	if b < 0 || b > 1 {
		return false, fmt.Errorf("unexpected byte value %v is <0 or >1", b)
	}
	return b == 1, nil
}

// Float32 read a float32 value from the packet
// Advances 4 bytes with the packet
// Errors if the value cannot be converted to a float32
// Errors if there are less then 4 bytes remaining to be read
func (p *packet) Float32() (float32, error) {
	b, err := p.readN(4)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return math.Float32frombits(binary.LittleEndian.Uint32(b)), nil
}

// Uint64 read an uint64 value from the packet
// Advances 8 bytes with the packet
// Errors if the value cannot be converted to an uint64
// Errors if there are less than 8 bytes remaining to be read
func (p *packet) Uint64() (uint64, error) {
	b, err := p.readN(8)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return binary.LittleEndian.Uint64(b), nil
}

// Uint32 read an uint32 value from the packet
// Advances 4 bytes with the packet
// Errors if the value cannot be converted to an uint32
// Errors if there are less than 4 bytes remaining to be read
func (p *packet) Uint32() (uint32, error) {
	b, err := p.readN(4)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return binary.LittleEndian.Uint32(b), nil
}

// Uint16 read an uint16 value from the packet
// Advances 2 bytes with the packet
// Errors if the value cannot be converted to an uint16
// Errors if there are less than 2 bytes remaining to be read
func (p *packet) Uint16() (uint16, error) {
	b, err := p.readN(2)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return binary.LittleEndian.Uint16(b), nil
}

// Uint8 read an uint8 value from the packet
// Advances a byte with the packet
// Errors if the value cannot be converted to an uint8
// Errors if there are no bytes remaining to be read
func (p *packet) Uint8() (uint8, error) {
	b, err := p.read()
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return b, nil
}

// Int64 read an int64 value from the packet
// Advances 8 bytes with the packet
// Errors if the value cannot be converted to an int64
// Errors if there are less than 8 bytes remaining to be read
func (p *packet) Int64() (int64, error) {
	b, err := p.readN(8)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return int64(binary.LittleEndian.Uint64(b)), nil
}

// Int32 read an int32 value from the packet
// Advances 4 bytes with the packet
// Errors if the value cannot be converted to an int32
// Errors if there are less then 4 bytes remaining to be read
func (p *packet) Int32() (int32, error) {
	b, err := p.readN(4)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return int32(binary.LittleEndian.Uint32(b)), nil
}

// Int16 read an int16 value from the packet
// Advances 2 bytes with the packet
// Errors if the value cannot be converted to an int16
// Errors if there are less then 2 bytes remaining to be read
func (p *packet) Int16() (int16, error) {
	b, err := p.readN(2)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return int16(binary.LittleEndian.Uint16(b)), nil
}

// Int8 read an int8 value from the packet
// Advances a byte with the packet
// Errors if the value cannot be converted to an int8
// Errors if there are no bytes remaining to be read
func (p *packet) Int8() (int8, error) {
	b, err := p.read()
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes: %v", err)
	}

	return int8(b), nil
}
