package packet

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

// Serializer handler to be able to parse packets into structs
type Serializer interface {
	// Deserialize a provided packet setting the result to the provided ptr to a struct
	Deserialize(data Packet, dest interface{}) error
}

type serializer struct {
	reflections map[string][]packetField
	lock *sync.Mutex
}

type packetField struct {
	packetIdx int
	name      string
}

// NewSerializer creates a new packet parser
func NewSerializer() Serializer {
	return &serializer{
		reflections: make(map[string][]packetField),
		lock: &sync.Mutex{},
	}
}

// Deserialize a provided packet setting the result to the provided ptr to a struct
func (p *serializer) Deserialize(data Packet, dest interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(dest))
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("must pass a struct type into the reflection parser")
	}

	fields, err := p.inspect(val)
	if err != nil {
		return fmt.Errorf("failed to inspect receiver: %v", err)
	}

	for _, field := range fields {
		f := val.FieldByName(field.name)
		err = p.setVal(f, data)
		if err != nil {
			return fmt.Errorf("failed to set %v value: %v", field.name, err)
		}
	}

	return nil
}

// inspect parse out the fields of the packet to read into the struct
func (p *serializer) inspect(val reflect.Value) ([]packetField, error) {
	var fields []packetField

	p.lock.Lock()
	defer p.lock.Unlock()

	t := val.Type()
	structName := t.String()
	cached, hasCached := p.reflections[structName]
	if hasCached {
		return cached, nil
	}

	fields, err := generateSortedFields(val)
	if err != nil {
		return fields, fmt.Errorf("failed to generate fields: %v", err)
	}

	p.reflections[structName] = fields
	return fields, nil
}

// setVal read the value from the packet for a given field
func (p *serializer) setVal(f reflect.Value, data Packet) error {
	if !f.CanSet() {
		return fmt.Errorf("field is not settable")
	}

	kind := f.Kind()
	switch kind {
	case reflect.Struct:
		return p.parseStruct(f, data)
	case reflect.Array:
		return p.parseArray(f, data)
	case reflect.Bool:
		return p.parseBool(f, data)
	case reflect.Float32:
		return p.parseFloat32(f, data)
	case reflect.Uint64:
		return p.parseUint64(f, data)
	case reflect.Uint32:
		return p.parseUint32(f, data)
	case reflect.Uint16:
		return p.parseUint16(f, data)
	case reflect.Uint8:
		return p.parseUint8(f, data)
	case reflect.Int64:
		return p.parseInt64(f, data)
	case reflect.Int32:
		return p.parseInt32(f, data)
	case reflect.Int16:
		return p.parseInt16(f, data)
	case reflect.Int8:
		return p.parseInt8(f, data)
	}
	return fmt.Errorf("unsupported type %v", kind.String())
}

// parseStruct parse a struct out of the packet stream
func (p *serializer) parseStruct(f reflect.Value, d Packet) error {
	structVal := reflect.Indirect(reflect.New(f.Type()))
	fields, err := p.inspect(structVal)
	if err != nil {
		return fmt.Errorf("unable to parse struct: %v", err)
	}
	for _, field := range fields {
		structField := structVal.FieldByName(field.name)
		err := p.setVal(structField, d)
		if err != nil {
			return fmt.Errorf("unable to set struct field %v: %v", field.name, err)
		}
	}
	f.Set(structVal)

	return nil
}

// parseArray parse an array out of the packet stream
func (p *serializer) parseArray(f reflect.Value, d Packet) error {
	size := f.Type().Len()
	for i := 0; i < size; i += 1 {
		idxVal := f.Index(i)

		err := p.parseStruct(idxVal, d)
		if err != nil {
			return fmt.Errorf("unable to parse array item %v: %v", i, err)
		}

		f.Index(i).Set(idxVal)
	}
	return nil
}

// parseBool parse an bool out of the packet stream
func (p *serializer) parseBool(f reflect.Value, d Packet) error {
	val, err := d.Bool()
	if err != nil {
		return fmt.Errorf("unable to parse bool: %v", err)
	}
	f.SetBool(val)
	return nil
}

// parseFloat32 parse an float32 out of the packet stream
func (p *serializer) parseFloat32(f reflect.Value, d Packet) error {
	val, err := d.Float32()
	if err != nil {
		return fmt.Errorf("unable to parse float: %v", err)
	}
	f.SetFloat(float64(val))
	return nil
}

// parseUint64 parse an uint64 out of the packet stream
func (p *serializer) parseUint64(f reflect.Value, d Packet) error {
	val, err := d.Uint64()
	if err != nil {
		return fmt.Errorf("unable to parse int8: %v", err)
	}
	f.SetUint(val)
	return nil
}

// parseUint32 parse an uint32 out of the packet stream
func (p *serializer) parseUint32(f reflect.Value, d Packet) error {
	val, err := d.Uint32()
	if err != nil {
		return fmt.Errorf("unable to parse int8: %v", err)
	}
	f.SetUint(uint64(val))
	return nil
}

// parseUint16 parse an uint16 out of the packet stream
func (p *serializer) parseUint16(f reflect.Value, d Packet) error {
	val, err := d.Uint16()
	if err != nil {
		return fmt.Errorf("unable to parse int8: %v", err)
	}
	f.SetUint(uint64(val))
	return nil
}

// parseUint8 parse an uint8 out of the packet stream
func (p *serializer) parseUint8(f reflect.Value, d Packet) error {
	val, err := d.Uint8()
	if err != nil {
		return fmt.Errorf("unable to parse int8: %v", err)
	}
	f.SetUint(uint64(val))
	return nil
}

// parseInt64 parse an int64 out of the packet stream
func (p *serializer) parseInt64(f reflect.Value, d Packet) error {
	val, err := d.Int64()
	if err != nil {
		return fmt.Errorf("unable to parse int16: %v", err)
	}
	f.SetInt(val)
	return nil
}

// parseInt32 parse an int32 out of the packet stream
func (p *serializer) parseInt32(f reflect.Value, d Packet) error {
	val, err := d.Int32()
	if err != nil {
		return fmt.Errorf("unable to parse int16: %v", err)
	}
	f.SetInt(int64(val))
	return nil
}

// parseInt16 parse an int16 out of the packet stream
func (p *serializer) parseInt16(f reflect.Value, d Packet) error {
	val, err := d.Int16()
	if err != nil {
		return fmt.Errorf("unable to parse int16: %v", err)
	}
	f.SetInt(int64(val))
	return nil
}

// parseInt8 parse an int8 out of the packet stream
func (p *serializer) parseInt8(f reflect.Value, d Packet) error {
	val, err := d.Int8()
	if err != nil {
		return fmt.Errorf("unable to parse int8: %v", err)
	}
	f.SetInt(int64(val))
	return nil
}

func generateSortedFields(v reflect.Value) ([]packetField, error) {
	var fields []packetField

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		packetIdx := f.Tag.Get("packet")
		if packetIdx != "" {
			pInt, err := strconv.ParseInt(packetIdx, 10, 32)
			if err != nil {
				return fields, fmt.Errorf("invalid packet index received: %v", err)
			}
			fields = append(fields, packetField{
				packetIdx: int(pInt),
				name:      f.Name,
			})
		}
	}

	fields = sortFields(fields)
	return fields, nil
}

// sortFields sort the fields of the packet by the packet index
func sortFields(fields []packetField) []packetField {
	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].packetIdx < fields[j].packetIdx
	})

	return fields
}

