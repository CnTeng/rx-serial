package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/CnTeng/rx-serial/internal/types"
)

type Endian string

const (
	LittleEndian Endian = "little"
	BigEndian    Endian = "big"
)

type FieldConfig struct {
	IsEscape  *bool            `json:"is_escape"`
	ByteOrder *types.ByteOrder `json:"byte_order"`
}

type Field struct {
	Name   string           `json:"name"`
	Length json.RawMessage  `json:"length"`
	Type   string           `json:"type"`
	Value  *json.RawMessage `json:"value"`
	Config *FieldConfig     `json:"config"`
	Extra  *json.RawMessage `json:"extra"`
}

func (fc *FieldConfig) patch(tc *types.TypeConfig) {
	if fc.IsEscape != nil {
		tc.IsEscape = *fc.IsEscape
	}
	if fc.ByteOrder != nil {
		tc.ByteOrder = *fc.ByteOrder
	}
}

func unmarshalLength(raw json.RawMessage) (*types.IntType, error) {
	var len int
	if err := json.Unmarshal(raw, &len); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	}
	return types.NewConstInt(len), nil
}

func unmarshalDynamicLength(raw json.RawMessage, ct *types.ComplexType) (*types.IntType, error) {
	var len int
	if err := json.Unmarshal(raw, &len); err == nil {
		return types.NewConstInt(len), nil
	}

	var name string
	if err := json.Unmarshal(raw, &name); err == nil {
		t, err := ct.GetType(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get type: %w", err)
		}
		i, ok := t.(*types.IntType)
		if !ok {
			return nil, fmt.Errorf("invalid type: %T", t)
		}
		return i, nil
	}

	return nil, fmt.Errorf("failed to unmarshal length")
}

func (f *Field) toType(ct *types.ComplexType) (types.Type, error) {
	switch f.Type {
	case "complex":
		return f.toComplex(ct)
	case "const":
		return f.toConst(ct)
	case "int":
		return f.toInt(ct)
	case "bytes":
		return f.toBytes(ct)
	case "string":
		return f.toString(ct)
	case "timestamp":
		return f.toTimeStamp(ct)
	case "crc":
		return f.toCRC16(ct)
	default:
		return nil, fmt.Errorf("unknown type: %s", f.Type)
	}
}

func (f *Field) toComplex(ct *types.ComplexType) (*types.ComplexType, error) {
	c := &types.ComplexType{
		TypeMeta: types.NewTypeMeta().
			WithName(f.Name).
			WithConfig(ct.Config),
	}

	if f.Config != nil {
		f.Config.patch(&c.Config)
	}

	return c, nil
}

func (f *Field) toConst(ct *types.ComplexType) (*types.ConstType, error) {
	c := &types.ConstType{
		TypeMeta: types.NewTypeMeta().
			WithName(f.Name).
			WithConfig(ct.Config),
	}

	if f.Config != nil {
		f.Config.patch(&c.Config)
	}

	if len, err := unmarshalLength(f.Length); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	} else {
		c.WithLength(len)
	}

	if err := json.Unmarshal(*f.Value, &c.Value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return c, nil
}

func (f *Field) toInt(ct *types.ComplexType) (*types.IntType, error) {
	i := &types.IntType{
		TypeMeta: types.NewTypeMeta().
			WithName(f.Name).
			WithConfig(ct.Config),
	}

	if f.Config != nil {
		f.Config.patch(&i.Config)
	}

	if len, err := unmarshalDynamicLength(f.Length, ct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	} else {
		i.WithLength(len)
	}

	return i, nil
}

func (f *Field) toBytes(ct *types.ComplexType) (*types.BytesType, error) {
	bs := &types.BytesType{
		TypeMeta: types.NewTypeMeta().
			WithName(f.Name).
			WithConfig(ct.Config),
	}

	if f.Config != nil {
		f.Config.patch(&bs.Config)
	}

	if len, err := unmarshalDynamicLength(f.Length, ct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	} else {
		bs.WithLength(len)
	}

	return bs, nil
}

func (f *Field) toString(ct *types.ComplexType) (*types.StringType, error) {
	s := &types.StringType{
		TypeMeta: types.NewTypeMeta().
			WithName(f.Name).
			WithConfig(ct.Config),
	}

	if f.Config != nil {
		f.Config.patch(&s.Config)
	}

	if len, err := unmarshalDynamicLength(f.Length, ct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	} else {
		s.WithLength(len)
	}

	return s, nil
}

func (f *Field) toTimeStamp(ct *types.ComplexType) (*types.TimeStamp, error) {
	ts := &types.TimeStamp{
		TypeMeta: types.NewTypeMeta().
			WithName(f.Name).
			WithConfig(ct.Config),
	}

	if f.Config != nil {
		f.Config.patch(&ts.Config)
	}

	if len, err := unmarshalDynamicLength(f.Length, ct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	} else {
		ts.WithLength(len)
	}

	return ts, nil
}

func (f *Field) toCRC16(ct *types.ComplexType) (*types.CRC16Type, error) {
	crc := &types.CRC16Type{
		IntType: types.IntType{
			TypeMeta: types.NewTypeMeta().
				WithName(f.Name).
				WithConfig(ct.Config),
		},
	}

	if f.Config != nil {
		f.Config.patch(&crc.Config)
	}

	if len, err := unmarshalDynamicLength(f.Length, ct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal length: %w", err)
	} else {
		crc.WithLength(len)
	}

	scoop := []string{}
	if err := json.Unmarshal(*f.Extra, &scoop); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scoop: %w", err)
	}

	for _, s := range scoop {
		t, err := ct.GetType(s)
		if err != nil {
			return nil, fmt.Errorf("failed to get type: %w", err)
		}
		crc.Scoop = append(crc.Scoop, t)
	}

	return crc, nil
}
