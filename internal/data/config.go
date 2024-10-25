package data

import (
	"encoding/json"
	"fmt"
	"os"
)

type Endian string

const (
	BigEndian    Endian = "big"
	LittleEndian Endian = "little"
)

type Config struct {
	TopName string    `json:"top_struct"`
	Top     *Struct   `json:"-"`
	Structs []*Struct `json:"-"`
}

func NewConfig(configPath string, structsPath string) (*Config, error) {
	c, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{Structs: make([]*Struct, 0)}
	if err := json.Unmarshal(c, config); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(structsPath)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		s, err := os.Open(structsPath + "/" + e.Name())
		if err != nil {
			return nil, err
		}
		defer s.Close()

		dec := json.NewDecoder(s)
		dec.UseNumber()

		st := &Struct{}
		if err := dec.Decode(st); err != nil {
			return nil, err
		}

		if st.Name == config.TopName {
			config.Top = st
		}

		config.Structs = append(config.Structs, st)
	}

	if config.Top == nil {
		return nil, fmt.Errorf("top struct %s not found", config.TopName)
	}

	return config, nil
}

func (c *Config) Parse(b []byte) (*Data, error) {
	return c.Top.Parse(b, c.Structs)
}
