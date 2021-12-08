package lib

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"

	"github.com/coyove/nj"
	"github.com/coyove/nj/internal"
)

func getEncB64(enc *base64.Encoding, padding rune) *base64.Encoding {
	if padding != '=' {
		enc = enc.WithPadding(padding)
	}
	return enc
}

func getEncB32(enc *base32.Encoding, padding rune) *base32.Encoding {
	if padding != '=' {
		enc = enc.WithPadding(padding)
	}
	return enc
}

var encDecProto = nj.NamedObject("EncodeDecode", 0).
	SetMethod("encode", func(e *nj.Env) {
		i := e.This("_e")
		e.A = nj.Str(i.(interface{ EncodeToString([]byte) string }).EncodeToString(e.Get(0).Safe().Bytes()))
	}, "").
	SetMethod("decode", func(e *nj.Env) {
		i := e.This("_e")
		v, err := i.(interface{ DecodeString(string) ([]byte, error) }).DecodeString(e.Str(0))
		internal.PanicErr(err)
		e.A = nj.Bytes(v)
	}, "").
	SetPrototype(nj.NamedObject("EncoderDecoder", 0).
		SetMethod("encoder", func(e *nj.Env) {
			enc := nj.Nil
			buf := &bytes.Buffer{}
			switch encoding := e.This("_e").(type) {
			default:
				enc = nj.ValueOf(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = nj.ValueOf(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = nj.ValueOf(base64.NewEncoder(encoding, buf))
			}
			e.A = nj.NamedObject("Encoder", 0).
				SetProp("_f", nj.ValueOf(enc)).
				SetProp("_b", nj.ValueOf(buf)).
				SetMethod("value", func(e *nj.Env) {
					e.A = nj.Str(e.This("_b").(*bytes.Buffer).String())
				}, "").
				SetMethod("bytes", func(e *nj.Env) {
					e.A = nj.Bytes(e.This("_b").(*bytes.Buffer).Bytes())
				}, "").
				SetPrototype(nj.WriteCloserProto).
				ToValue()
		}, "").
		SetMethod("decoder", func(e *nj.Env) {
			src := nj.NewReader(e.Get(0))
			dec := nj.Nil
			switch encoding := e.This("_e").(type) {
			case *base64.Encoding:
				dec = nj.ValueOf(base64.NewDecoder(encoding, src))
			case *base32.Encoding:
				dec = nj.ValueOf(base32.NewDecoder(encoding, src))
			default:
				dec = nj.ValueOf(hex.NewDecoder(src))
			}
			e.A = nj.NamedObject("Decoder", 0).
				SetProp("_f", nj.ValueOf(dec)).
				SetPrototype(nj.ReaderProto).
				ToValue()
		}, ""))

func init() {
	nj.Globals.SetProp("hex", nj.NamedObject("hex", 0).SetPrototype(encDecProto.Prototype()).ToValue())
	nj.Globals.SetProp("base64", nj.NamedObject("base64", 0).
		SetProp("std", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB64(base64.StdEncoding, '='))).ToValue()).
		SetProp("url", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB64(base64.URLEncoding, '='))).ToValue()).
		SetProp("std2", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB64(base64.StdEncoding, -1))).ToValue()).
		SetProp("url2", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB64(base64.URLEncoding, -1))).ToValue()).
		SetPrototype(encDecProto).
		ToValue())
	nj.Globals.SetProp("base32", nj.NamedObject("base32", 0).
		SetProp("std", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB32(base32.StdEncoding, '='))).ToValue()).
		SetProp("hex", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB32(base32.HexEncoding, '='))).ToValue()).
		SetProp("std2", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB32(base32.StdEncoding, -1))).ToValue()).
		SetProp("hex2", nj.NewObject(1).SetPrototype(encDecProto).SetProp("_e", nj.ValueOf(getEncB32(base32.HexEncoding, -1))).ToValue()).
		SetPrototype(encDecProto).
		ToValue())
}
