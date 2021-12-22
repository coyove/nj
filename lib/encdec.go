package lib

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"

	"github.com/coyove/nj/bas"
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

var encDecProto = bas.NamedObject("EncodeDecode", 0).
	SetMethod("encode", func(e *bas.Env) {
		i := e.This("_e")
		e.A = bas.Str(i.(interface{ EncodeToString([]byte) string }).EncodeToString(e.Get(0).Safe().Bytes()))
	}, "").
	SetMethod("decode", func(e *bas.Env) {
		i := e.This("_e")
		v, err := i.(interface{ DecodeString(string) ([]byte, error) }).DecodeString(e.Str(0))
		internal.PanicErr(err)
		e.A = bas.Bytes(v)
	}, "").
	SetPrototype(bas.NamedObject("EncoderDecoder", 0).
		SetMethod("encoder", func(e *bas.Env) {
			enc := bas.Nil
			buf := &bytes.Buffer{}
			switch encoding := e.This("_e").(type) {
			default:
				enc = bas.ValueOf(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = bas.ValueOf(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = bas.ValueOf(base64.NewEncoder(encoding, buf))
			}
			e.A = bas.NamedObject("Encoder", 0).
				SetProp("_f", bas.ValueOf(enc)).
				SetProp("_b", bas.ValueOf(buf)).
				SetMethod("value", func(e *bas.Env) {
					e.A = bas.Str(e.This("_b").(*bytes.Buffer).String())
				}, "").
				SetMethod("bytes", func(e *bas.Env) {
					e.A = bas.Bytes(e.This("_b").(*bytes.Buffer).Bytes())
				}, "").
				SetPrototype(bas.WriteCloserProto).
				ToValue()
		}, "").
		SetMethod("decoder", func(e *bas.Env) {
			src := bas.NewReader(e.Get(0))
			dec := bas.Nil
			switch encoding := e.This("_e").(type) {
			case *base64.Encoding:
				dec = bas.ValueOf(base64.NewDecoder(encoding, src))
			case *base32.Encoding:
				dec = bas.ValueOf(base32.NewDecoder(encoding, src))
			default:
				dec = bas.ValueOf(hex.NewDecoder(src))
			}
			e.A = bas.NamedObject("Decoder", 0).
				SetProp("_f", bas.ValueOf(dec)).
				SetPrototype(bas.ReaderProto).
				ToValue()
		}, ""))

func init() {
	bas.Globals.SetProp("hex", bas.NamedObject("hex", 0).SetPrototype(encDecProto.Prototype()).ToValue())
	bas.Globals.SetProp("base64", bas.NamedObject("base64", 0).
		SetProp("std", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB64(base64.StdEncoding, '='))).ToValue()).
		SetProp("url", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB64(base64.URLEncoding, '='))).ToValue()).
		SetProp("std2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB64(base64.StdEncoding, -1))).ToValue()).
		SetProp("url2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB64(base64.URLEncoding, -1))).ToValue()).
		SetPrototype(encDecProto).
		ToValue())
	bas.Globals.SetProp("base32", bas.NamedObject("base32", 0).
		SetProp("std", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB32(base32.StdEncoding, '='))).ToValue()).
		SetProp("hex", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB32(base32.HexEncoding, '='))).ToValue()).
		SetProp("std2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB32(base32.StdEncoding, -1))).ToValue()).
		SetProp("hex2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(getEncB32(base32.HexEncoding, -1))).ToValue()).
		SetPrototype(encDecProto).
		ToValue())
}
