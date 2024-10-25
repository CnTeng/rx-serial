package data

import (
	"fmt"
	"math/big"
	"strings"
)

type Data struct {
	keys []string
	data map[string]any
}

func NewData() *Data {
	return &Data{
		keys: make([]string, 0),
		data: make(map[string]any),
	}
}

func (d *Data) GetInt(name string) (int, error) {
	for k, v := range d.data {
		if k != name {
			continue
		}

		v, ok := v.(int)
		if ok {
			return v, nil
		}
	}

	return 0, fmt.Errorf("field %s not found", name)
}

func (d *Data) GetString(name string) (string, error) {
	for k, v := range d.data {
		if k != name {
			continue
		}

		v, ok := v.(string)
		if ok {
			return v, nil
		}
	}

	return "", fmt.Errorf("field %s not found", name)
}

func (d *Data) GetBytes(name string) ([]byte, error) {
	for k, v := range d.data {
		if k != name {
			continue
		}

		v, ok := v.([]byte)
		if ok {
			return v, nil
		}
	}

	return nil, fmt.Errorf("field %s not found", name)
}

func (d *Data) SetInt(name string, b []byte) error {
	v := int(big.NewInt(0).SetBytes(b).Int64())
	d.data[name] = v
	d.keys = append(d.keys, name)
	return nil
}

func (d *Data) SetBytes(name string, b []byte) error {
	d.data[name] = b
	d.keys = append(d.keys, name)
	return nil
}

func (d *Data) String() string {
	builder := strings.Builder{}

	for _, k := range d.keys {
		v := d.data[k]
		switch v := v.(type) {
		case []byte:
			builder.WriteString(fmt.Sprintf("%s: ", k))
			for _, b := range v {
				builder.WriteString(fmt.Sprintf("%02x ", b))
			}
			builder.WriteString("\n")
		case *Data:
			builder.WriteString(fmt.Sprintf("%s: \n%v", k, v))
		default:
			builder.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
	}

	return builder.String()
}
