// SPDX-FileCopyrightText: 2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"encoding/json/v2"
	"encoding/xml"
	"testing"
)

type obj struct {
	S string `json:"s"`
	I int    `json:"i"`
}

func TestAssertion_EncodingEqual(t *testing.T) {
	a := New(t, false)
	jsonU := func(data []byte, v any) error { return json.Unmarshal(data, v) }

	a.EncodingEqual[obj]([]byte(`{"s":"s","i":5}`), []byte(`{"s":"s",	"i":5}`), jsonU)
	a.EncodingEqual[obj]([]byte(`<xml><S>s</S><I>5</I></xml>`), []byte(` <xml> <S>s</S><I>5</I></xml>`), xml.Unmarshal)
}

func TestAssertion_EncodingNotEqual(t *testing.T) {
	a := New(t, false)
	jsonU := func(data []byte, v any) error { return json.Unmarshal(data, v) }

	a.EncodingNotEqual[obj]([]byte(`{"s":"s","i":5}`), []byte(`{"s":"s","i":6}`), jsonU)
	a.EncodingNotEqual[obj]([]byte(`<xml><S>s</S><I>6</I></xml>`), []byte(` <xml> <S>s</S><I>5</I></xml>`), xml.Unmarshal)
}
