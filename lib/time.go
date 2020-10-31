package lib

import (
	"bytes"
	"time"
	"unicode/utf8"

	"github.com/coyove/script"
)

func init() {
	script.AddGlobalValue("Go_timefmt", func(env *script.Env) {
		f := env.InStr(0, "")
		switch f {
		case "ANSIC":
			f = time.ANSIC
		case "UnixDate":
			f = time.UnixDate
		case "RubyDate":
			f = time.RubyDate
		case "RFC822":
			f = time.RFC822
		case "RFC822Z":
			f = time.RFC822Z
		case "RFC850":
			f = time.RFC850
		case "RFC1123":
			f = time.RFC1123
		case "RFC1123Z":
			f = time.RFC1123Z
		case "RFC3339":
			f = time.RFC3339
		case "RFC3339Nano":
			f = time.RFC3339Nano
		case "Kitchen":
			f = time.Kitchen
		case "Stamp":
			f = time.Stamp
		case "StampMilli":
			f = time.StampMilli
		case "StampMicro":
			f = time.StampMicro
		case "StampNano":
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
			tt = time.Now()
		}

		r := tt.Format(f)
		env.A = env.NewString(r)
	},
		"Go_timefmt(format_string) => string",
		"Go_timefmt(format_string, time.Time) => string",
		"\tformat doc: https://www.php.net/manual/datetime.format.php",
	)
}
