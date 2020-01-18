package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	tNone = iota
	tRoot
	tComma
	tColon
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

// ErrInvalidToken signals that something could not be converted to a token.
var ErrInvalidToken = errors.New("invalid token")

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
	return 0, ErrInvalidToken
}

func parseJSON(data []byte) ([]JSONElement, error) {
	// create a channel to receive errors from
	exc := make(chan error)

	// start parsing in a goroutine which emits the resulting tokens to this channel
	tokens := make(chan JSONElement)
	go func() {
		// close error channel on exit to terminate waiting loop
		defer close(exc)

		length := len(data)
		cur := 0
		for cur != length {
			switch data[cur] {
			case ' ', '\n', '\r', '\t':
				fmt.Println(cur, "whitespace")
				// skip whitespace
				cur++
			case '{':
				fmt.Println(cur, "opening braces")
				tokens <- JSONElement{tpe: tObjectStart, offset: cur}
				cur++
			case '}':
				fmt.Println(cur, "closing braces")
				tokens <- JSONElement{tpe: tObjectEnd, offset: cur}
				cur++
			case '[':
				fmt.Println(cur, "opening brackets")
				tokens <- JSONElement{tpe: tArrayStart, offset: cur}
				cur++
			case ']':
				fmt.Println(cur, "closing brackets")
				tokens <- JSONElement{tpe: tArrayEnd, offset: cur}
				cur++
			case ':':
				fmt.Println(cur, "colon")
				tokens <- JSONElement{tpe: tColon, offset: cur}
				cur++
			case ',':
				fmt.Println(cur, "comma")
				tokens <- JSONElement{tpe: tComma, offset: cur}
				cur++
			case '"':
				fmt.Println(cur, "string")
				size, err := findMatchingQuotes(data, cur, length)
				if err != nil {
					exc <- err
					return
				}
				tokens <- JSONElement{tpe: tString, offset: cur}
				// two quotes plus payload
				cur += 2 + size
			case 'n':
				fmt.Println(cur, "null")
				if cur+4 > length {
					exc <- ErrInvalidToken
					return
				}
				if (data[cur+1] != 'u') || (data[cur+2] != 'l') || (data[cur+3] != 'l') {
					exc <- ErrInvalidToken
				}
				tokens <- JSONElement{tpe: tNull, offset: cur}
				cur += 4
			case 't':
				fmt.Println(cur, "true")
				if cur+4 > length {
					exc <- ErrInvalidToken
					return
				}
				if (data[cur+1] != 'r') || (data[cur+2] != 'u') || (data[cur+3] != 'e') {
					exc <- ErrInvalidToken
				}
				tokens <- JSONElement{tpe: tBool, offset: cur}
				cur += 4
			case 'f':
				fmt.Println(cur, "false")
				if cur+5 > length {
					exc <- ErrInvalidToken
					return
				}
				if (data[cur+1] != 'a') || (data[cur+2] != 'l') || (data[cur+3] != 's') || (data[cur+4] != 'e') {
					exc <- ErrInvalidToken
				}
				tokens <- JSONElement{tpe: tBool, offset: cur}
				cur += 5
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				fmt.Println(cur, "number")
				cur++
			default:
				fmt.Println(cur, "unexpected")
				exc <- ErrInvalidToken
				return
			}
		}
	}()

	res := make([]JSONElement, 0, 10)
	for {
		select {
		case err := <-exc:
			// Note that "err" can be nil, which happens when the channel
			// is closed and it just means that the goroutine finished.
			fmt.Println("received error", err)
			return res, err
		case elem := <-tokens:
			fmt.Println("received element", elem)
			res = append(res, elem)
			break
		}
	}
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
