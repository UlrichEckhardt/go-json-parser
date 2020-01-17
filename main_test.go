package main

import (
	"testing"
)

type jsonTest struct {
	data     []byte
	elements []JSONElement
}

func TestParseJSON(t *testing.T) {
	cases := map[string]jsonTest{
		"empty": {
			data:     []byte(``),
			elements: []JSONElement{},
		},
		"string 1": {
			data: []byte(`""`),
			elements: []JSONElement{
				JSONElement{tpe: tString, offset: 0},
			},
		},
		"string 2": {
			data: []byte(`"string"`),
			elements: []JSONElement{
				JSONElement{tpe: tString, offset: 0},
			},
		},
		"null": {
			data: []byte(`null`),
			elements: []JSONElement{
				JSONElement{tpe: tNull, offset: 0},
			},
		},
		"bool true": {
			data: []byte(`true`),
			elements: []JSONElement{
				JSONElement{tpe: tBool, offset: 0},
			},
		},
		"bool false": {
			data: []byte(`false`),
			elements: []JSONElement{
				JSONElement{tpe: tBool, offset: 0},
			},
		},
		"array 1": {
			data: []byte(`[]`),
			elements: []JSONElement{
				JSONElement{tpe: tArrayStart, offset: 0},
				JSONElement{tpe: tArrayEnd, offset: 1},
			},
		},
		"object 1": {
			data: []byte(`{}`),
			elements: []JSONElement{
				JSONElement{tpe: tObjectStart, offset: 0},
				JSONElement{tpe: tObjectEnd, offset: 1},
			},
		},
		"whitespace 1": {
			data: []byte{'\n', '"', '"'},
			elements: []JSONElement{
				JSONElement{tpe: tString, offset: 1},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			elements, err := parseJSON(c.data)
			if err != nil {
				t.Error("unexpected failure", err)
				return
			}
			for i, e := range c.elements {
				if i >= len(elements) {
					t.Errorf("element %d is missing", i)
					break
				}
				if e != elements[i] {
					t.Errorf("element %d differs", i)
					break
				}
			}
			if len(elements) > len(c.elements) {
				t.Error("too many elements in output")
			}

		})
	}
}
