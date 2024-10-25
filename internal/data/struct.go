package data

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type FieldConfig struct {
	Endian Endian `json:"endian"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Len  any    `json:"len"`

	FieldConfig
}

type Struct struct {
	Name       string      `json:"name"`
	Global     FieldConfig `json:"global"`
	Fields     []Field     `json:"fields"`
	Escape     bool        `json:"escape"`
	EscapeChar int         `json:"escape_char"`
}

func (fc *FieldConfig) patch(field *Field) {
	if field.Endian == "" && fc.Endian != "" {
		field.Endian = fc.Endian
	}
}

func (s *Struct) unescape(b []byte) []byte {
	if !s.Escape {
		return b
	}

	var r []byte
	for i := 0; i < len(b); {
		if b[i] == byte(s.EscapeChar) {
			r = append(r, b[i+1])
			i += 2
		} else {
			r = append(r, b[i])
			i++
		}
	}

	return r
}

func (s *Struct) Parse(b []byte, ss []*Struct) (*Data, error) {
	data := NewData()
	b = s.unescape(b)
	len := 0

	for _, f := range s.Fields {
		l := 0
		s.Global.patch(&f)

		switch f.Len.(type) {
		case string:
			n, err := data.GetInt(f.Len.(string))
			if err != nil {
				return nil, err
			}
			l = n
		case json.Number:
			n, err := f.Len.(json.Number).Int64()
			if err != nil {
				return nil, err
			}
			l = int(n)
		default:
			return nil, fmt.Errorf("unknown type: %s", reflect.TypeOf(f.Len))
		}

		bp := b[len : len+l]

		switch f.Type {
		case "int":
			if err := data.SetInt(f.Name, bp); err != nil {
				return nil, err
			}
		case "bytes":
			if err := data.SetBytes(f.Name, bp); err != nil {
				return nil, err
			}
		case "time":
			if err := data.SetBytes(f.Name, bp); err != nil {
				return nil, err
			}
		default:
			findStruct := false

			for _, s := range ss {
				if s.Name == f.Type {
					d, err := s.Parse(bp, ss)
					if err != nil {
						return nil, err
					}
					data.keys = append(data.keys, f.Name)
					data.data[f.Name] = d
					findStruct = true
					break
				}
			}
			if !findStruct {
				return nil, fmt.Errorf("unknown type: %s", f.Type)
			}
		}

		len += l
	}

	return data, nil
}
