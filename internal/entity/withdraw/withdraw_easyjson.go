// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package withdraw

import (
	json "encoding/json"
	order "github.com/k1nky/gophermart/internal/entity/order"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson4f4a6fc6DecodeGithubComK1nkyGophermartInternalEntityWithdraw(in *jlexer.Lexer, out *Withdraw) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "sum":
			out.Sum = float32(in.Float32())
		case "order":
			out.Number = order.OrderNumber(in.String())
		case "processed_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ProcessedAt).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4f4a6fc6EncodeGithubComK1nkyGophermartInternalEntityWithdraw(out *jwriter.Writer, in Withdraw) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"sum\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float32(float32(in.Sum))
	}
	{
		const prefix string = ",\"order\":"
		out.RawString(prefix)
		out.String(string(in.Number))
	}
	{
		const prefix string = ",\"processed_at\":"
		out.RawString(prefix)
		out.Raw((in.ProcessedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Withdraw) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4f4a6fc6EncodeGithubComK1nkyGophermartInternalEntityWithdraw(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdraw) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4f4a6fc6EncodeGithubComK1nkyGophermartInternalEntityWithdraw(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdraw) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4f4a6fc6DecodeGithubComK1nkyGophermartInternalEntityWithdraw(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdraw) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4f4a6fc6DecodeGithubComK1nkyGophermartInternalEntityWithdraw(l, v)
}
