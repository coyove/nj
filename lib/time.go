package lib

import (
	"bytes"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/coyove/script"
)

func init() {
	script.AddGlobalValue("strtime", func(env *script.Env) {
		f := env.Get(0).StringDefault("")
		switch strings.ToLower(f) {
		case "ansic":
			f = time.ANSIC
		case "unixdate":
			f = time.UnixDate
		case "rubydate":
			f = time.RubyDate
		case "rfc822":
			f = time.RFC822
		case "rfc822z":
			f = time.RFC822Z
		case "rfc850":
			f = time.RFC850
		case "rfc1123":
			f = time.RFC1123
		case "rfc1123z":
			f = time.RFC1123Z
		case "rfc3339":
			f = time.RFC3339
		case "rfc3339nano":
			f = time.RFC3339Nano
		case "kitchen":
			f = time.Kitchen
		case "stamp":
			f = time.Stamp
		case "stampmilli":
			f = time.StampMilli
		case "stampmicro":
			f = time.StampMicro
		case "stampnano":
			f = time.StampNano
		default:
			buf := bytes.Buffer{}
			for len(f) > 0 {
				r, sz := utf8.DecodeRuneInString(f)
				if sz == 0 {
					break
				}
				switch r {
				case 'd':
					buf.WriteString("02")
				case 'D':
					buf.WriteString("Mon")
				case 'j':
					buf.WriteString("2")
				case 'l':
					buf.WriteString("Monday")
				case 'F':
					buf.WriteString("January")
				case 'z':
					buf.WriteString("002")
				case 'm':
					buf.WriteString("01")
				case 'M':
					buf.WriteString("Jan")
				case 'n':
					buf.WriteString("1")
				case 'Y':
					buf.WriteString("2006")
				case 'y':
					buf.WriteString("06")
				case 'a':
					buf.WriteString("pm")
				case 'A':
					buf.WriteString("PM")
				case 'g':
					buf.WriteString("3")
				case 'G':
					buf.WriteString("15")
				case 'h':
					buf.WriteString("03")
				case 'H':
					buf.WriteString("15")
				case 'i':
					buf.WriteString("04")
				case 's':
					buf.WriteString("05")
				case 'u':
					buf.WriteString("05.000000")
				case 'v':
					buf.WriteString("05.000")
				case 'O':
					buf.WriteString("+0700")
				case 'P':
					buf.WriteString("+07:00")
				case 'T':
					buf.WriteString("MST")
				case 'c': //	ISO 8601
					buf.WriteString("2006-01-02T15:04:05+07:00")
				case 'r': //	RFC 2822
					buf.WriteString("Mon, 02 Jan 2006 15:04:05 +0700")
				default:
					buf.WriteRune(r)
				}
				f = f[sz:]
			}
			f = buf.String()
		}

		tt, ok := env.Get(1).Interface().(time.Time)
		if !ok {
			ts := env.Get(1).IntDefault(0)
			if ts > 0 {
				if ts < 1<<33 {
					tt = time.Unix(ts, 0)
				} else if ts < 1<<33*1e3 {
					tt = time.Unix(0, ts*1e6)
				} else if ts < 1<<33*1e6 {
					tt = time.Unix(0, ts*1e3)
				} else {
					tt = time.Unix(0, ts)
				}
			} else {
				tt = time.Now()
			}
		}

		r := tt.Format(f)
		env.A = script.String(r)
	},
		"strtime(format_string) => string",
		"strtime(format_string, time.Time) => string",
		"strtime(format_string, unix_timestamp) => string",
		"\tformat doc: https://www.php.net/manual/datetime.format.php",
	)
}
