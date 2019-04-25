package rpn

import "fmt"

func isSep(r byte, sq bool) bool {
	x := r == ' ' || r == '\t' || r == '\r' || r == '\n' || r == '[' || r == ']' || r == '"' || r == '!'
	if sq {
		x = x || r == '\''
	}
	return x
}

func isAlpha(r byte, num bool) bool {
	if num {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
	}
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isVar(v string) bool {
	for i := 0; i < len(v); i++ {
		if !isAlpha(v[i], i > 0) {
			return false
		}
	}
	return true
}

func isNum(r byte, e bool) bool {
	x := (r >= '0' && r <= '9') || r == '.' || r == '+' || r == '-'
	if e {
		x = x || r == 'e' || r == 'E'
	}
	return x
}

func unexpectedToken(tok *Token) error {
	return fmt.Errorf(ErrUnexpectedInput, tok)
}
