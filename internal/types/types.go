package types

import "github.com/google/go-cmp/cmp"

type Type interface {
	Name() string
	Len() int
	IsEscape() bool
	Parse(data []byte) error
	Serialize() []byte
	Equal(other Type) bool
}

type ByteOrder string

const (
	LittleEndian ByteOrder = "little"
	BigEndian    ByteOrder = "big"
)

var defaultConfig = TypeConfig{
	ByteOrder: LittleEndian,
}

type TypeConfig struct {
	IsEscape  bool      `json:"is_escape"`
	ByteOrder ByteOrder `json:"byte_order"`
}

type TypeMeta struct {
	name   string
	Length *IntType
	Config TypeConfig
}

func NewTypeMeta() *TypeMeta {
	return &TypeMeta{Config: defaultConfig}
}

func (tm *TypeMeta) Name() string {
	return tm.name
}

func (tm *TypeMeta) Len() int {
	return tm.Length.Value
}

func (tm *TypeMeta) IsEscape() bool {
	return tm.Config.IsEscape
}

func (tm *TypeMeta) WithName(name string) *TypeMeta {
	tm.name = name
	return tm
}

func (tm *TypeMeta) WithLength(len *IntType) *TypeMeta {
	tm.Length = len
	return tm
}

func (tm *TypeMeta) WithConfig(config TypeConfig) *TypeMeta {
	tm.Config = config
	return tm
}

func (tm *TypeMeta) Equal(o *TypeMeta) bool {
	if tm == nil || o == nil {
		return tm == o
	}

	if tm.name != o.name {
		return false
	}

	if (tm.Length == nil) != (o.Length == nil) {
		return false
	}
	if tm.Length != nil && !tm.Length.Equal(o.Length) {
		return false
	}

	if !cmp.Equal(tm.Config, o.Config) {
		return false
	}

	return true
}
