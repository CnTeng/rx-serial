package types

type StringType struct {
	Value string
	*TypeMeta
}

func (s *StringType) Parse(data []byte) error {
	s.Value = string(data[:s.Len()])
	return nil
}

func (s *StringType) Serialize() []byte {
	return []byte(s.Value)[:s.Len()]
}

func (s *StringType) Equal(other Type) bool {
	if os, ok := other.(*StringType); !ok {
		return false
	} else {
		return s.Value == os.Value && s.TypeMeta.Equal(os.TypeMeta)
	}
}
