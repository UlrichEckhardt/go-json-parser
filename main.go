package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	tNone = iota
	tRoot
	tObjectStart
	tObjectEnd
	tArrayStart
	tArrayEnd
	tString
	tNull
	tInt
	tFloat
	tBool
)

// JSONElement is an element of the JSON syntax tree.
type JSONElement struct {
	tpe    int // type according to the t* constants above
	offset int // offset of the element within the input data
}

func findMatchingQuotes(data []byte, cur, length int) (int, error) {
	// skip opening quotes
	cur++

	for res := 0; cur+res != length; res++ {
		if data[cur+res] == '"' {
			return res, nil
		}
	}
	return 0, errors.New("missing closing quotes for string")
}

func parseJSON(data []byte) ([]JSONElement, error) {
	res := make([]JSONElement, 0, 10)
	length := len(data)
	cur := 0
	for cur != length {
		switch data[cur] {
		case ' ', '\n', '\r', '\t':
			fmt.Println(cur, "whitespace")
			// skip whitespace
			cur++
		case '{':
			res = append(res, JSONElement{tpe: tObjectStart, offset: cur})
			cur++
		case '}':
			res = append(res, JSONElement{tpe: tObjectEnd, offset: cur})
			cur++
		case '[':
			res = append(res, JSONElement{tpe: tArrayStart, offset: cur})
			cur++
		case ']':
			res = append(res, JSONElement{tpe: tArrayEnd, offset: cur})
			cur++
		case '"':
			size, err := findMatchingQuotes(data, cur, length)
			if err != nil {
				return nil, err
			}
			res = append(res, JSONElement{tpe: tString, offset: cur})
			// two quotes plus payload
			cur += 2 + size
		default:
			return nil, errors.New("unexpected byte")
		}
	}

	return res, nil
}

func main() {
	defer fmt.Println("Done.")
	if len(os.Args) != 2 {
		fmt.Println("expected exactly one arguent")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("failed to open file", err)
		return
	}

	data := make([]byte, 0, 10000)
	size, err := file.Read(data)
	if err != nil {
		fmt.Println("failed to read file", err)
		return
	}
	if size == cap(data) {
		fmt.Println("file too large")
		return
	}

	res, err := parseJSON(data)
	if err != nil {
		fmt.Println("failed to parse", err)
		return
	}

	fmt.Println("parsed", res)
}
