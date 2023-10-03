package uuid

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"testing"
	"time"
)

var (
	_ driver.Valuer            = UUID{}
	_ encoding.BinaryMarshaler = UUID{}
	_ encoding.TextMarshaler   = UUID{}
	_ json.Marshaler           = UUID{}

	_ encoding.BinaryUnmarshaler = (*UUID)(nil)
	_ encoding.TextUnmarshaler   = (*UUID)(nil)
	_ json.Unmarshaler           = (*UUID)(nil)
	_ sql.Scanner                = (*UUID)(nil)
)

func TestNewV3(t *testing.T) {
	namespace := newUUID()
	name := []byte("testing")

	u1 := NewV3(namespace, name)
	verifyVariant(t, u1)
	verifyVersion(t, u1, 3)

	u2 := NewV3(namespace, name)
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("NewV3 returned different UUIDs with the same namespace & name: %s vs %s",
			u1.Format(), u2.Format())
	}
}

func TestNewV4(t *testing.T) {
	u1 := newUUID()
	verifyVariant(t, u1)
	verifyVersion(t, u1, 4)

	u2 := newUUID()
	if bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("NewV4 returned equal UUIDs: %s vs %s", u1.Format(), u2.Format())
	}
}

func TestNewV5(t *testing.T) {
	namespace := newUUID()
	name := []byte("testing")

	u1 := NewV5(namespace, name)
	verifyVariant(t, u1)
	verifyVersion(t, u1, 5)

	u2 := NewV5(namespace, name)
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("NewV5 returned different UUIDs with the same namespace & name: %s vs %s",
			u1.Format(), u2.Format())
	}
}

func TestVersion(t *testing.T) {
	var table = []struct {
		name       string
		u          func() UUID
		expVersion int
	}{
		{
			name: "v3",
			u: func() UUID {
				u := newUUID()
				return NewV3(u, []byte("test"))
			},
			expVersion: 3,
		},
		{
			name:       "v4",
			u:          func() UUID { return newUUID() },
			expVersion: 4,
		},
		{
			name: "v5",
			u: func() UUID {
				u := newUUID()
				return NewV5(u, []byte("test"))
			},
			expVersion: 5,
		},
		{
			name: "v7",
			u: func() UUID {
				u, err := NewV7(time.UnixMilli(1000))
				if err != nil {
					panic(err)
				}
				return u
			},
			expVersion: 7,
		},
	}

	for i := 0; i < len(table); i++ {
		ts := table[i]
		t.Run(ts.name, func(t *testing.T) {
			v := ts.u().Version()
			if v != ts.expVersion {
				t.Fatalf("Incorrect version: %d", v)
			}
		})
	}
}

func verifyVariant(t *testing.T, u UUID) {
	v := u[8] >> 6
	if v != 2 {
		t.Fatalf("Expected variant '10', got '%x'", v)
	}
}

func verifyVersion(t *testing.T, u UUID, version byte) {
	v := u[6] >> 4
	if v != version {
		t.Fatalf("Expected version '%x', got '%x'", version, v)
	}
}

func TestFormat(t *testing.T) {
	u := newUUID()
	f := u.Format()
	if f[8] != '-' || f[13] != '-' || f[18] != '-' || f[23] != '-' {
		t.Fatalf("Invalid UUID format: %s", f)
	}
}

func TestBytes(t *testing.T) {
	u := newUUID()
	b := u.Bytes()
	if len(b) != 36 {
		t.Fatalf("Invalid UUID length: %d (expected 36)", len(b))
	}
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		t.Fatalf("Invalid UUID format: %s", b)
	}
	bb := u.Format()
	if !bytes.Equal(bb[:], b) {
		t.Fatalf("Format and FormatString return different UUIDs: %s vs %s",
			bb, b)
	}
}

func TestString(t *testing.T) {
	u := newUUID()
	s := u.String()
	if len(s) != 36 {
		t.Fatalf("Invalid UUID length: %d (expected 36)", len(s))
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		t.Fatalf("Invalid UUID format: %s", s)
	}
	b := u.Format()
	if !bytes.Equal(b[:], []byte(s)) {
		t.Fatalf("Format and FormatString return different UUIDs: %s vs %s",
			b, s)
	}
}

func TestMarshalBinary(t *testing.T) {
	u := newUUID()
	b, err := u.MarshalBinary()
	if err != nil {
		t.Fatalf("Unexpected binary marshaling error: %s", err.Error())
	}
	if !bytes.Equal(u[:], b[:]) {
		t.Fatalf("Unexpected binary marshaling result: %v", b)
	}
}

func TestUnmarshalBinary(t *testing.T) {
	u1 := newUUID()
	u2 := UUID{}
	err := u2.UnmarshalBinary(u1[:])
	if err != nil {
		t.Fatalf("Unexpected binary unmarshaling error: %s", err.Error())
	}
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("Unexpected binary unmarshaling result: %v", u2)
	}
	u2 = UUID{}
	err = u2.UnmarshalBinary([]byte{0})
	if err != ErrInvalidUUID {
		t.Fatalf("Unexpected binary unmarshaling error: %v", err)
	}
}

func TestMarshalJSON(t *testing.T) {
	u := newUUID()
	b, err := u.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected json marshaling error: %s", err.Error())
	}
	if !bytes.Equal(b[1:37], u.Bytes()) {
		t.Fatalf("Unexpected json marshaling result: %v", b)
	}
	if b[0] != '"' || b[37] != '"' {
		t.Fatalf("Unexpected json marshaling result: %v", b)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	u1 := newUUID()
	b, err := u1.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected json marshaling error: %s", err.Error())
	}
	u2 := UUID{}
	err = u2.UnmarshalJSON(b)
	if err != nil {
		t.Fatalf("Unexpected json unmarshaling error: %s", err.Error())
	}
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("Unexpected json unmarshaling result: %v", u2)
	}
	u2 = UUID{}
	err = u2.UnmarshalJSON([]byte{0})
	if err != ErrInvalidUUID {
		t.Fatalf("Unexpected json unmarshaling error: %v", err)
	}
}

func TestMarshalText(t *testing.T) {
	u := newUUID()
	b, err := u.MarshalText()
	if err != nil {
		t.Fatalf("Unexpected text marshaling error: %s", err.Error())
	}
	if !bytes.Equal(b[:], u.Bytes()) {
		t.Fatalf("Unexpected text marshaling result: %v", b)
	}
}

func TestUnmarshalText(t *testing.T) {
	u1 := newUUID()
	u2 := UUID{}
	err := u2.UnmarshalText(u1.Bytes())
	if err != nil {
		t.Fatalf("Unexpected text unmarshaling error: %s", err.Error())
	}
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("Unexpected text unmarshaling result: %v", u2)
	}
	u2 = UUID{}
	err = u2.UnmarshalText([]byte{0})
	if err != ErrInvalidUUID {
		t.Fatalf("Unexpected text unmarshaling error: %v", err)
	}
}

func TestValue(t *testing.T) {
	u := newUUID()
	v, err := u.Value()
	if err != nil {
		t.Fatalf("Unexpected value error: %s", err.Error())
	}
	if !bytes.Equal(v.([]byte), u.Bytes()) {
		t.Fatalf("Unexpected value result: %v", v)
	}
}

func TestScan(t *testing.T) {
	u1 := newUUID()
	u2 := UUID{}
	err := u2.Scan(u1.Bytes())
	if err != nil {
		t.Fatalf("Unexpected scan error: %s", err.Error())
	}
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("Unexpected scan result: %v", u2)
	}
	u2 = UUID{}
	err = u2.Scan(u1.String())
	if err != nil {
		t.Fatalf("Unexpected scan error: %s", err.Error())
	}
	if !bytes.Equal(u1[:], u2[:]) {
		t.Fatalf("Unexpected scan result: %v", u2)
	}
	u2 = UUID{}
	err = u2.Scan(1)
	if err != ErrInvalidUUID {
		t.Fatalf("Unexpected scan error: %v", err)
	}
}

func TestParseString(t *testing.T) {
	s := "9e754ef6-8dd9-4903-af43-7aea99bfb1fe"
	u, err := ParseString(s)
	if err != nil {
		t.Fatalf("Unexpected parsing error: %s", err.Error())
	}
	if u.String() != s {
		t.Fatalf("Invalid parsed UUID: %s", u.String())
	}
}

func TestParse16(t *testing.T) {
	b := newUUID()
	u, err := Parse(b[:])
	if err != nil {
		t.Fatalf("Unexpected parsing error: %s", err.Error())
	}
	if b.String() != u.String() {
		t.Fatalf("Invalid parsed UUID: %s", u.String())
	}
}

func TestParse32(t *testing.T) {
	b := []byte("9e754ef68dd94903af437aea99bfb1fe")
	u, err := Parse(b)
	if err != nil {
		t.Fatalf("Unexpected parsing error: %s", err.Error())
	}
	f := u.Format()
	if !bytes.Equal(f[:], []byte("9e754ef6-8dd9-4903-af43-7aea99bfb1fe")) {
		t.Fatalf("Unexpected parsing result: %s", f)
	}

	b = []byte("9e754ef68dd94903af437aea99bfb1fg")
	_, err = Parse(b)
	if err == nil {
		t.Fatalf("Unexpected parsing success: %s", b)
	}
}

func TestParse36(t *testing.T) {
	b := []byte("9e754ef6-8dd9-4903-af43-7aea99bfb1fe")
	u, err := Parse(b)
	if err != nil {
		t.Fatalf("Unexpected parsing error: %s", err.Error())
	}
	f := u.Format()
	if !bytes.Equal(f[:], b) {
		t.Fatalf("Unexpected parsing result: %s", f)
	}
}

func TestParse36Error(t *testing.T) {
	bb := [][]byte{
		[]byte("9e754ef6-8dd9-4903-af437aea99bfb1fef"),
		[]byte("9e754gf6-8dd9-4903-af43-7aea99bfb1fe"),
	}
	for _, b := range bb {
		_, err := Parse(b)
		if err != ErrInvalidUUID {
			t.Fatalf("Unexpected parsing pass: %s", b)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	s := "bad"
	_, err := ParseString(s)
	if err == nil {
		t.Fatal("Unexpected parsing success")
	}
}

func TestTime(t *testing.T) {
	now := time.UnixMilli(time.Now().UnixMilli())
	u, err := NewV7(now)
	if err != nil {
		t.Fatalf("Unexpected error generating uuid v7: %s", err.Error())
	}
	ut, ok := u.Time()
	if !ok {
		t.Fatal("Unable to parse time from V7 UUID")
	}
	if !now.Equal(ut) {
		t.Fatalf("Time not equal to original: %v vs %v", now, ut)
	}

	u = newUUID()
	if _, ok := u.Time(); ok {
		t.Fatal("Should not be able to parse time from a V4 UUID")
	}
}

func TestMust(t *testing.T) {
	u := Must(NewV4())
	if u.Version() != 4 {
		t.Fatalf("Unexpected UUID version: %d", u.Version())
	}
}

func newUUID() UUID {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	return u
}

func BenchmarkParse(b *testing.B) {
	buf := []byte("9e754ef6-8dd9-4903-af43-7aea99bfb1fe")
	for i := 0; i < b.N; i++ {
		_, _ = Parse(buf)
	}
}

func BenchmarkParseString(b *testing.B) {
	str := "9e754ef6-8dd9-4903-af43-7aea99bfb1fe"
	for i := 0; i < b.N; i++ {
		_, _ = ParseString(str)
	}
}

func BenchmarkNewV3(b *testing.B) {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	name := []byte("test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewV3(u, name)
	}
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewV4()
	}
}

func BenchmarkNewV5(b *testing.B) {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	name := []byte("test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewV5(u, name)
	}
}

func BenchmarkFormat(b *testing.B) {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.Format()
	}
}

func BenchmarkFormatString(b *testing.B) {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}
