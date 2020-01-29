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
	tNumber
	tBool
)

// ErrInvalidToken signals that something could not be converted to a token.
var ErrInvalidToken = errors.New("invalid token")

// ErrInvalidStructure signals that a valid token was encountered in the wrong place.
// In particular, that means closing tokens (")", "}") outside the scope of the
// according aggregate value type. Further, it means commas outside of aggregate types
// and colons anywhere but as a separator between key and value of an object value.
var ErrInvalidStructure = errors.New("invalid structure")

// JSONElement is an element of the JSON syntax tree.
type JSONElement struct {
	tpe    int // type according to the t* constants above
	offset int // offset of the element within the input data
	parent int // index of the parent element in the output data
}

func findMatchingQuotes(data []byte, cur, length int) (int, error) {
	// skip opening quotes
	cur++

	backslashed := false
	for res := 0; cur+res != length; res++ {
		switch c := data[cur+res]; {
		case c == '"':
			if !backslashed {
				return res, nil
			}
			backslashed = false
		case c == '\\':
			// invert backslash-state
			backslashed = !backslashed
		case c < 32:
			// control byte
			return res, ErrInvalidToken
		default:
			backslashed = false
		}
	}
	return 0, ErrInvalidToken
}

func findEndOfNumber(data []byte, cur, length int) (int, error) {
	const (
		optionalSign = iota
		nonfractionStart
		nonfractionContinued
		radixSeparator
		fractionStart
		fractionContinued
		exponentSeparator
		exponentSign
		exponentStart
		exponentContinued
	)

	res := 0
	state := optionalSign
loop:
	for {
		// get next glyph
		if cur+res == length {
			break loop
		}
		c := data[cur+res]

		switch state {
		case optionalSign:
			// if it's a minus sign, skip it
			if c == '-' {
				res++
			}
			state = nonfractionStart

		case nonfractionStart:
			switch c {
			case '0':
				// consume non-fractional digit
				res++
				state = radixSeparator
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// consume non-fractional digit
				res++
				state = nonfractionContinued
			default:
				break loop
			}

		case nonfractionContinued:
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// consume non-fractional digits
				res++
				state = nonfractionContinued
			default:
				// Anything else isn't consumed. Instead, we treat it as
				// (optional) radix separator and continue from that point.
				state = radixSeparator
			}

		case radixSeparator:
			switch c {
			case '.':
				// consume radix separator
				res++
				state = fractionStart
			default:
				state = exponentSeparator
			}

		case fractionStart:
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// consume fractional digits
				res++
				state = nonfractionContinued
			default:
				break loop
			}

		case fractionContinued:
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// consume fractional digits
				res++
			default:
				state = exponentSeparator
			}

		case exponentSeparator:
			switch c {
			case 'e', 'E':
				// consume exponent separator
				res++
				state = exponentSign
			default:
				break loop
			}

		case exponentSign:
			switch c {
			case '+', '-':
				// consume exponent sign
				res++
				state = exponentStart
			default:
				state = exponentStart
			}

		case exponentStart:
			// Note: It seems that "1.e01" is valid, although "01.2" isn't, hence the
			// numbers of the exponent are not parsed like the nonfractional digits.
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// consume exponent digit
				res++
				state = exponentContinued
			default:
				break loop
			}

		case exponentContinued:
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// consume exponent digit
				res++
				state = exponentContinued
			default:
				break loop
			}
		}
	}

	// check final state, there must not be incomplete parts
	switch state {
	case optionalSign, nonfractionStart, fractionStart, exponentSign, exponentStart:
		// incomplete number token
		return 0, ErrInvalidToken
	case nonfractionContinued, radixSeparator, fractionContinued, exponentSeparator, exponentContinued:
		return res, nil
	default:
		return 0, errors.New("invalid state parsing number")
	}
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
			case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				fmt.Println(cur, "number")
				size, err := findEndOfNumber(data, cur, length)
				if err != nil {
					exc <- err
					return
				}
				tokens <- JSONElement{tpe: tNumber, offset: cur}
				cur += size
			default:
				fmt.Println(cur, "unexpected")
				exc <- ErrInvalidToken
				return
			}
		}
	}()

	res := make([]JSONElement, 0, 10)
	res = append(res, JSONElement{tpe: tRoot})
	context := 0
	for {
		select {
		case err := <-exc:
			// Note that "err" can be nil, which happens when the channel
			// is closed and it just means that the goroutine finished.
			fmt.Println("received error", err)
			return res, err
		case elem := <-tokens:
			fmt.Println("received element", elem)
			// determine context changes
			switch elem.tpe {
			case tArrayStart, tObjectStart:
				// remember parent index for aggregate value
				elem.parent = context
				context = len(res)
			case tArrayEnd:
				if res[context].tpe != tArrayStart {
					// current context must be an array
					return nil, ErrInvalidStructure
				}
				// validate all intermediate tokens
				const (
					start = iota // initial state, next token must be a value if present
					comma        // next token must be a comma if present
					next         // next token must be present and not a comma
				)
				state := start
				for i := context + 1; i != len(res); i++ {
					t := res[i]
					// if this is not a direct child, ignore it
					if t.parent != context {
						continue
					}
					// if this is the end of a nested structure, ignore it
					if t.tpe == tObjectEnd || t.tpe == tArrayEnd {
						continue
					}

					switch t.tpe {
					case tObjectStart, tArrayStart, tBool, tNumber, tNull, tString:
						if state == comma {
							// expected a comma as separator, not a value
							return nil, ErrInvalidStructure
						}
						state = comma
					case tComma:
						if state != comma {
							// expected a value, not a comma as separator
							return nil, ErrInvalidStructure
						}
						state = next
					default:
						// unexpected token as array element
						return nil, ErrInvalidStructure
					}
				}
				if state == next {
					// brackets are not empty but don't end in a value
					return nil, ErrInvalidStructure
				}
				context = res[context].parent
				elem.parent = context
			case tObjectEnd:
				if res[context].tpe != tObjectStart {
					// current context must be an object
					return nil, ErrInvalidStructure
				}
				// validate all intermediate tokens
				const (
					start = iota // initial state, next token must be a string if present
					colon        // next token must be present and a colon
					value        // next token must be present and a value
					comma        // next token must be a comma if present
					next         // next token must be present and be a string
				)
				state := start
				for i := context + 1; i != len(res); i++ {
					t := res[i]
					// if this is not a direct child, ignore it
					if t.parent != context {
						continue
					}
					// if this is the end of a nested structure, ignore it
					if t.tpe == tObjectEnd || t.tpe == tArrayEnd {
						continue
					}

					switch state {
					case start, next:
						if t.tpe != tString {
							// expected a string as key
							return nil, ErrInvalidStructure
						}
						state = colon
					case colon:
						if t.tpe != tColon {
							// expected a colon as separator
							return nil, ErrInvalidStructure
						}
						state = value
					case value:
						switch t.tpe {
						case tObjectStart, tArrayStart, tBool, tNumber, tNull, tString:
							state = comma
						default:
							// expected a value
							return nil, ErrInvalidStructure
						}
					case comma:
						if t.tpe != tComma {
							// expected a comma as separator
							return nil, ErrInvalidStructure
						}
						state = next
					}
				}
				switch state {
				case colon, value, next:
					// braces are not empty but don't end in a value
					return nil, ErrInvalidStructure
				}
				context = res[context].parent
				elem.parent = context
			case tComma:
				if res[context].tpe != tArrayStart && res[context].tpe != tObjectStart {
					return nil, ErrInvalidStructure
				}
				elem.parent = context
			case tColon:
				if res[context].tpe != tObjectStart {
					return nil, ErrInvalidStructure
				}
				elem.parent = context
			default:
				elem.parent = context
			}
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
