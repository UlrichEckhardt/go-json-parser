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
