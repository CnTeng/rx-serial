package protocol

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CnTeng/rx-serial/internal/types"
)

type Protocol struct {
	top        *types.ComplexType
	complexMap map[string]*types.ComplexType
}

func NewProtocol(dir string) (*Protocol, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	p := &Protocol{
		complexMap: make(map[string]*types.ComplexType, len(files)),
	}

	for _, file := range files {
		data, err := os.ReadFile(file.Name())
		if err != nil {
			return nil, err
		}

		cc := &complexConfig{}
		if err := json.Unmarshal(data, cc); err != nil {
			return nil, err
		}

		c, err := cc.toComplexType()
		if err != nil {
			return nil, err
		}

		p.complexMap[c.Name()] = c
	}

	if err := p.sort(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Protocol) sort() error {
	nTop := 0
	for _, c := range p.complexMap {
		if c.IsTop {
			p.top = c
			nTop++
		}
	}

	if nTop != 1 {
		return fmt.Errorf("invalid top count: %d", nTop)
	}

	for _, c := range p.complexMap {
		for i, t := range c.Types {
			if t, ok := t.(*types.ComplexType); ok {
				if _, ok := p.complexMap[t.Name()]; !ok {
					return fmt.Errorf("unknown type: %s", t.Name())
				}

				c.Types[i] = p.complexMap[t.Name()]
			}
		}
	}

	return nil
}

func (p *Protocol) Parse(data []byte) error {
	return p.top.Parse(data)
}

type complexConfig struct {
	Name   string       `json:"name"`
	IsTop  bool         `json:"is_top"`
	Global *FieldConfig `json:"global"`
	Escape types.Escape `json:"escape"`
	Fields []Field      `json:"fields"`
}

func (cc *complexConfig) toComplexType() (*types.ComplexType, error) {
	ct := &types.ComplexType{
		IsTop:    cc.IsTop,
		Escape:   cc.Escape,
		TypeMeta: types.NewTypeMeta().WithName(cc.Name),
		TypesMap: make(map[string]types.Type),
	}

	if cc.Global != nil {
		cc.Global.patch(&ct.Config)
	}

	for _, f := range cc.Fields {
		t, err := f.toType(ct)
		if err != nil {
			return nil, err
		}
		ct.WithType(t)
	}

	return ct, nil
}
