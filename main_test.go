package main

import (
	"testing"
)

type jsonTest struct {
	data     []byte
	elements []JSONElement
	err      error
}

func TestParseJSON(t *testing.T) {
	cases := map[string]jsonTest{
		"empty": {
			data: []byte(``),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
			},
		},
		"string 1": {
			data: []byte(`""`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 2": {
			data: []byte(`"string"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 3": {
			data: []byte(`"\""`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 4": {
			data: []byte(`"\\"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 5": {
			data: []byte(`"\/"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 6": {
			data: []byte(`"\b"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 7": {
			data: []byte(`"\n"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 8": {
			data: []byte(`"\t"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"string 9": {
			data: []byte(`"\u1234"`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 0, parent: 0},
			},
		},
		"null": {
			data: []byte(`null`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNull, offset: 0, parent: 0},
			},
		},
		"bool true": {
			data: []byte(`true`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tBool, offset: 0, parent: 0},
			},
		},
		"bool false": {
			data: []byte(`false`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tBool, offset: 0, parent: 0},
			},
		},
		"number 1": {
			data: []byte(`0`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 0, parent: 0},
			},
		},
		"number 2": {
			data: []byte(`-1234`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 0, parent: 0},
			},
		},
		"number 3": {
			data: []byte(`1.2`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 0, parent: 0},
			},
		},
		"number 4": {
			data: []byte(`1E1`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 0, parent: 0},
			},
		},
		"number 5": {
			data: []byte(`1e-1`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 0, parent: 0},
			},
		},
		"number 6": {
			data: []byte(`0.314e+1`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 0, parent: 0},
			},
		},
		"array 1": {
			data: []byte(`[]`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tArrayStart, offset: 0, parent: 0},
				JSONElement{tpe: tArrayEnd, offset: 1, parent: 0},
			},
		},
		"array 2": {
			data: []byte(`["", true]`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tArrayStart, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 1, parent: 1},
				JSONElement{tpe: tComma, offset: 3, parent: 1},
				JSONElement{tpe: tBool, offset: 5, parent: 1},
				JSONElement{tpe: tArrayEnd, offset: 9, parent: 0},
			},
		},
		"array 3": {
			data: []byte(`[[""]]`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tArrayStart, offset: 0, parent: 0},
				JSONElement{tpe: tArrayStart, offset: 1, parent: 1},
				JSONElement{tpe: tString, offset: 2, parent: 2},
				JSONElement{tpe: tArrayEnd, offset: 4, parent: 1},
				JSONElement{tpe: tArrayEnd, offset: 5, parent: 0},
			},
		},
		"array 5": {
			data: []byte(`[1, 2]`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tArrayStart, offset: 0, parent: 0},
				JSONElement{tpe: tNumber, offset: 1, parent: 1},
				JSONElement{tpe: tComma, offset: 2, parent: 1},
				JSONElement{tpe: tNumber, offset: 4, parent: 1},
				JSONElement{tpe: tArrayEnd, offset: 5, parent: 0},
			},
		},
		"object 1": {
			data: []byte(`{}`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tObjectStart, offset: 0, parent: 0},
				JSONElement{tpe: tObjectEnd, offset: 1, parent: 0},
			},
		},
		"object 2": {
			data: []byte(`{"k": true}`),
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tObjectStart, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 1, parent: 1},
				JSONElement{tpe: tColon, offset: 4, parent: 1},
				JSONElement{tpe: tBool, offset: 6, parent: 1},
				JSONElement{tpe: tObjectEnd, offset: 10, parent: 0},
			},
		},
		"whitespace 1": {
			data: []byte{'\n', '"', '"'},
			elements: []JSONElement{
				JSONElement{tpe: tRoot, offset: 0, parent: 0},
				JSONElement{tpe: tString, offset: 1, parent: 0},
			},
		},
		"invalid 1": {
			data: []byte(`a`),
			err:  ErrInvalidToken,
		},
		"invalid 3": {
			data: []byte{0},
			err:  ErrInvalidToken,
		},
		"invalid 4": {
			data: []byte(`nil`),
			err:  ErrInvalidToken,
		},
		"invalid 7": {
			data: []byte(".2"),
			err:  ErrInvalidToken,
		},
		"invalid 8": {
			data: []byte("1."),
			err:  ErrInvalidToken,
		},
		"invalid 9": {
			data: []byte("1.2E"),
			err:  ErrInvalidToken,
		},
		"invalid 10": {
			data: []byte("1.2E+"),
			err:  ErrInvalidToken,
		},
		"invalid 11": {
			data: []byte("1.2e-"),
			err:  ErrInvalidToken,
		},
		"invalid 12": {
			data: []byte("01.2"),
			err:  ErrInvalidToken,
		},
		"invalid 13": {
			data: []byte("01"),
			err:  ErrInvalidToken,
		},
		"invalid string 1": {
			data: []byte(`"`),
			err:  ErrInvalidToken,
		},
		"invalid string 2": {
			data: []byte{'"', 0, '"'},
			err:  ErrInvalidToken,
		},
		"invalid string 3": {
			data: []byte("\"\n\""),
			err:  ErrInvalidToken,
		},
		"invalid string 4": {
			data: []byte(`"\M"`),
			err:  ErrInvalidToken,
		},
		"invalid string 5": {
			data: []byte(`"\ux1234"`),
			err:  ErrInvalidToken,
		},
		"invalid structure 1": {
			data: []byte("}"),
			err:  ErrInvalidStructure,
		},
		"invalid structure 2": {
			data: []byte("]"),
			err:  ErrInvalidStructure,
		},
		"invalid structure 3": {
			data: []byte("[,]"),
			err:  ErrInvalidStructure,
		},
		"invalid structure 4": {
			data: []byte("[1,]"),
			err:  ErrInvalidStructure,
		},
		"invalid structure 5": {
			data: []byte("[1:2]"),
			err:  ErrInvalidStructure,
		},
		"invalid structure 6": {
			data: []byte("[true false]"),
			err:  ErrInvalidStructure,
		},
		"invalid structure 7": {
			data: []byte(`{"k":}`),
			err:  ErrInvalidStructure,
		},
		"invalid structure 8": {
			data: []byte(`{"k":"v",}`),
			err:  ErrInvalidStructure,
		},
		"invalid structure 9": {
			data: []byte(`{1: 2}`),
			err:  ErrInvalidStructure,
		},
		"invalid structure 10": {
			data: []byte(`{"1": 2 : 3}`),
			err:  ErrInvalidStructure,
		},
		"invalid structure 11": {
			data: []byte(`{"k" "v"}`),
			err:  ErrInvalidStructure,
		},
		"invalid structure 12": {
			data: []byte(`1, 2`),
			err:  ErrInvalidStructure,
		},
		"invalid structure 13": {
			data: []byte(`"k":"v"`),
			err:  ErrInvalidStructure,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			elements, err := parseJSON(c.data)
			if c.err != nil {
				if err == nil {
					t.Error("expected error missing")
				} else if c.err != err {
					t.Log("expected error", c.err)
					t.Log("received error", err)
					t.Error("wrong error")
				}
				return
			}
			if c.err == nil && err != nil {
				t.Error("unexpected failure", err)
				return
			}
			for i, e := range c.elements {
				if i >= len(elements) {
					t.Errorf("element %d is missing", i)
					break
				}
				if e != elements[i] {
					t.Log("expected element", e)
					t.Log("received element", elements[i])
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
