package types

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type ComplexType struct {
	IsTop    bool
	Escape   Escape
	Types    []Type
	TypesMap map[string]Type
	*TypeMeta
}

func (c *ComplexType) WithEscape(e Escape) *ComplexType {
	c.Escape = e
	return c
}

func (c *ComplexType) WithType(t Type) *ComplexType {
	c.Types = append(c.Types, t)
	c.TypesMap[t.Name()] = t
	return c
}

func (c *ComplexType) WithTypes(types []Type) *ComplexType {
	for _, t := range types {
		c.WithType(t)
	}
	return c
}

func (c *ComplexType) GetType(name string) (Type, error) {
	t, ok := c.TypesMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown type: %s", name)
	}
	return t, nil
}

func (c *ComplexType) Len() int {
	l := 0
	for _, t := range c.Types {
		l += t.Len()
	}
	return l
}

func (c *ComplexType) Parse(data []byte) error {
	data = c.Escape.unescape(data)
	for _, t := range c.Types {
		if err := t.Parse(data); err != nil {
			return err
		}
		data = data[t.Len():]
	}
	return nil
}

func (c *ComplexType) Serialize() []byte {
	var data []byte
	for _, t := range c.Types {
		b := t.Serialize()
		if t.IsEscape() {
			b = c.Escape.escape(b)
		}
		data = append(data, b...)
	}
	return data
}

func (c *ComplexType) Equal(o Type) bool {
	oc, ok := o.(*ComplexType)
	if !ok {
		return false
	}

	if c.IsTop != oc.IsTop {
		return false
	}

	if !cmp.Equal(c.Escape, oc.Escape) {
		return false
	}

	if len(c.Types) != len(oc.Types) {
		return false
	}
	for i, t := range c.Types {
		if !t.Equal(oc.Types[i]) {
			return false
		}
	}

	if len(c.TypesMap) != len(oc.TypesMap) {
		return false
	}
	for k, t := range c.TypesMap {
		ot, ok := oc.TypesMap[k]
		if !ok {
			return false
		}

		if !t.Equal(ot) {
			return false
		}
	}

	return c.TypeMeta.Equal(oc.TypeMeta)
}
