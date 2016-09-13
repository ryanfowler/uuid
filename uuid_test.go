package uuid

import (
	"bytes"
	"testing"
)

func TestNewV3(t *testing.T) {
	namespace, err := NewV4()
	if err != nil {
		t.Fatal(err)
	}
	name := []byte("testing")

	u1 := NewV3(namespace, name)
	verifyVariant(t, u1)
	verifyVersion(t, u1, 3)

	u2 := NewV3(namespace, name)
	if !bytes.Equal(u1[:], u2[:]) {
		t.Errorf("NewV3 returned different UUIDs with the same namespace & name: %s vs %s",
			u1.Format(), u2.Format())
	}
}

func TestNewV4(t *testing.T) {
	u1, err := NewV4()
	if err != nil {
		t.Fatal(err)
	}
	verifyVariant(t, u1)
	verifyVersion(t, u1, 4)

	u2, err := NewV4()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(u1[:], u2[:]) {
		t.Errorf("NewV4 returned equal UUIDs: %s vs %s", u1.Format(), u2.Format())
	}
}

func TestNewV5(t *testing.T) {
	namespace, err := NewV4()
	if err != nil {
		t.Fatal(err)
	}
	name := []byte("testing")

	u1 := NewV5(namespace, name)
	verifyVariant(t, u1)
	verifyVersion(t, u1, 5)

	u2 := NewV5(namespace, name)
	if !bytes.Equal(u1[:], u2[:]) {
		t.Errorf("NewV5 returned different UUIDs with the same namespace & name: %s vs %s",
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
				u := newV4(t)
				return NewV3(u, []byte("test"))
			},
			expVersion: 3,
		},
		{
			name:       "v4",
			u:          func() UUID { return newV4(t) },
			expVersion: 4,
		},
		{
			name: "v5",
			u: func() UUID {
				u := newV4(t)
				return NewV5(u, []byte("test"))
			},
			expVersion: 5,
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

func newV4(t *testing.T) UUID {
	u, err := NewV4()
	if err != nil {
		t.Fatalf("Unable to create UUID: %s", err.Error())
	}
	return u
}

func verifyVariant(t *testing.T, u UUID) {
	v := u[8] >> 6
	if v != 2 {
		t.Errorf("Expected variant '10', got '%x'", v)
	}
}

func verifyVersion(t *testing.T, u UUID, version byte) {
	v := u[6] >> 4
	if v != version {
		t.Errorf("Expected version '%x', got '%x'", version, v)
	}
}

func TestFormat(t *testing.T) {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	f := u.Format()
	if f[8] != '-' || f[13] != '-' || f[18] != '-' || f[23] != '-' {
		t.Errorf("Invalid UUID format: %s", f)
	}
}

func TestFormatString(t *testing.T) {
	u, err := NewV4()
	if err != nil {
		panic(err)
	}
	fs := u.FormatString()
	if len(fs) != 36 {
		t.Errorf("Invalid UUID length: %d (expected 36)", len(fs))
	}
	if fs[8] != '-' || fs[13] != '-' || fs[18] != '-' || fs[23] != '-' {
		t.Errorf("Invalid UUID format: %s", fs)
	}
	b := u.Format()
	if !bytes.Equal(b[:], []byte(fs)) {
		t.Errorf("Format and FormatString return different UUIDs: %s vs %s",
			b, fs)
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
		_ = u.FormatString()
	}
}
