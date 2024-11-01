package protocol

import (
	"encoding/json"
	"testing"

	"github.com/CnTeng/rx-serial/internal/types"
)

var topConfig = `{
  "name": "top",
  "is_top": true,
  "global": {
    "is_escape": true,
    "byte_order": "big"
  },
  "escape": {
    "escape_char": 16,
    "escaped_char": [
      2,
      3,
      16
    ]
  },
  "fields": [
    {
      "name": "const_test",
      "length": 1,
      "type": "const",
      "value": 42,
      "config": {
        "is_escape": false
      }
    },
    {
      "name": "int_test",
      "length": 2,
      "type": "int"
    },
    {
      "name": "string_test",
      "length": 16,
      "type": "string"
    },
    {
      "name": "timestamp_test",
      "length": 7,
      "type": "timestamp"
    },
		{
			"name": "dynamic_test",
			"length": "int_test",
			"type": "int"
		},
		{
			"name": "sub",
			"type": "complex"
		}
  ]
}`

var subConfig = `{
  "name": "sub",
  "global": {
    "byte_order": "big"
  },
  "fields": [
    {
      "name": "int_test",
      "length": 2,
      "type": "int"
    },
    {
      "name": "string_test",
      "length": 16,
      "type": "string"
    }
  ]
}`

func getTopWant() *types.ComplexType {
	want := &types.ComplexType{
		IsTop: true,
		Escape: types.Escape{
			EscapeChar:  16,
			EscapedChar: []byte{2, 3, 16},
		},
		TypeMeta: types.NewTypeMeta().
			WithName("top").
			WithConfig(types.TypeConfig{
				IsEscape:  true,
				ByteOrder: types.BigEndian,
			}),
		TypesMap: make(map[string]types.Type),
	}
	want.WithTypes([]types.Type{
		&types.ConstType{
			Value: 42,
			TypeMeta: types.NewTypeMeta().
				WithName("const_test").
				WithLength(types.NewConstInt(1)).
				WithConfig(types.TypeConfig{
					IsEscape:  false,
					ByteOrder: types.BigEndian,
				}),
		},
		&types.IntType{
			TypeMeta: types.NewTypeMeta().
				WithName("int_test").
				WithLength(types.NewConstInt(2)).
				WithConfig(types.TypeConfig{
					IsEscape:  true,
					ByteOrder: types.BigEndian,
				}),
		},
		&types.StringType{
			TypeMeta: types.NewTypeMeta().
				WithName("string_test").
				WithLength(types.NewConstInt(16)).
				WithConfig(types.TypeConfig{
					IsEscape:  true,
					ByteOrder: types.BigEndian,
				}),
		},
		&types.TimeStamp{
			TypeMeta: types.NewTypeMeta().
				WithName("timestamp_test").
				WithLength(types.NewConstInt(7)).
				WithConfig(types.TypeConfig{
					IsEscape:  true,
					ByteOrder: types.BigEndian,
				}),
		},
		&types.IntType{
			TypeMeta: types.NewTypeMeta().
				WithName("dynamic_test").
				WithConfig(types.TypeConfig{
					IsEscape:  true,
					ByteOrder: types.BigEndian,
				}),
		},
		&types.ComplexType{
			TypeMeta: types.NewTypeMeta().
				WithName("sub").
				WithConfig(types.TypeConfig{
					IsEscape:  true,
					ByteOrder: types.BigEndian,
				}),
		},
	})

	want.Types[4].(*types.IntType).WithLength(want.Types[1].(*types.IntType))

	return want
}

func getSubWant() *types.ComplexType {
	want := &types.ComplexType{
		TypeMeta: types.NewTypeMeta().
			WithName("sub").
			WithConfig(types.TypeConfig{ByteOrder: types.BigEndian}),
		TypesMap: make(map[string]types.Type),
	}
	want.WithTypes([]types.Type{
		&types.IntType{
			TypeMeta: types.NewTypeMeta().
				WithName("int_test").
				WithLength(types.NewConstInt(2)).
				WithConfig(types.TypeConfig{
					ByteOrder: types.BigEndian,
				}),
		},
		&types.StringType{
			TypeMeta: types.NewTypeMeta().
				WithName("string_test").
				WithLength(types.NewConstInt(16)).
				WithConfig(types.TypeConfig{
					ByteOrder: types.BigEndian,
				}),
		},
	})

	return want
}

func TestComplexConfig_toComplexType(t *testing.T) {
	cc := &complexConfig{}

	if err := json.Unmarshal([]byte(topConfig), cc); err != nil {
		t.Fatal(err)
	}

	c, err := cc.toComplexType()
	if err != nil {
		t.Fatal(err)
	}

	want := getTopWant()
	if !c.Equal(want) {
		t.Fatalf("want: %v, got: %v", want, c)
	}
}

func TestProtocol_sort(t *testing.T) {
	topCc := &complexConfig{}
	if err := json.Unmarshal([]byte(topConfig), topCc); err != nil {
		t.Fatal(err)
	}

	top, err := topCc.toComplexType()
	if err != nil {
		t.Fatal(err)
	}

	topWant := getTopWant()
	if !top.Equal(topWant) {
		t.Fatalf("want: %v, got: %v", topWant, top)
	}

	subCc := &complexConfig{}
	if err := json.Unmarshal([]byte(subConfig), subCc); err != nil {
		t.Fatal(err)
	}

	sub, err := subCc.toComplexType()
	if err != nil {
		t.Fatal(err)
	}

	subWant := getSubWant()
	if !sub.Equal(subWant) {
		t.Fatalf("want: %v, got: %v", subWant, sub)
	}

	wantProto := &Protocol{
		top:        top,
		complexMap: make(map[string]*types.ComplexType, 2),
	}
	wantProto.complexMap[top.Name()] = topWant
	wantProto.complexMap[sub.Name()] = subWant
	wantProto.complexMap[top.Name()].Types[5] = subWant

	p := &Protocol{
		complexMap: make(map[string]*types.ComplexType, 2),
	}
	p.complexMap[top.Name()] = top
	p.complexMap[sub.Name()] = sub

	if err := p.sort(); err != nil {
		t.Fatal(err)
	}

	if !p.top.Equal(wantProto.top) {
		t.Fatalf("want: %v, got: %v", wantProto.top, p.top)
	}

	if !p.complexMap[top.Name()].Equal(wantProto.complexMap[top.Name()]) {
		t.Fatalf("want: %v, got: %v", wantProto.complexMap[top.Name()], p.complexMap[top.Name()])
	}
}
