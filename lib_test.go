package potatolang

import "testing"

func TestSprintf(t *testing.T) {
	assert := func(a, b string) {
		if a != b {
			t.Error(a, b)
		}
	}

	sprintf := func(a string, args ...interface{}) string {
		env := NewEnv(nil)
		env.LocalPush(NewStringValue(a))
		for _, arg := range args {
			switch arg.(type) {
			case string:
				env.LocalPush(NewStringValue(arg.(string)))
			case float64:
				env.LocalPush(NewNumberValue(arg.(float64)))
			}
		}
		return _sprintf(env)
	}

	assert(sprintf("a"), "a")
	assert(sprintf("~1", "a"), "a")
	assert(sprintf("~1", 1.0), "1")
	assert(sprintf("~1~1", "a"), "aa")
	assert(sprintf("~1~2~1", "a"), "anila")
	assert(sprintf("~1~2~1", "a", "b"), "aba")
	assert(sprintf("~1~~2~1", "a", "b"), "a~2a")
	assert(sprintf("~1~a~1", "a", "b"), "aaa")
	assert(sprintf("~1~a~1~", "a"), "aaa")
	assert(sprintf("~1%2s%", "a"), " a")
	assert(sprintf("~1%d%", 1.0), "1")
	assert(sprintf("~1%02x%", 10.1), "0a")
}
